package frontend

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/acme/pg-lattice-proxy/internal/adapter"
	"github.com/acme/pg-lattice-proxy/internal/cancel"
	"github.com/acme/pg-lattice-proxy/internal/config"
	"github.com/acme/pg-lattice-proxy/internal/pgwire"
	"github.com/acme/pg-lattice-proxy/internal/platform"
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"github.com/acme/pg-lattice-proxy/internal/routing"
	"github.com/acme/pg-lattice-proxy/internal/session"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	cfg      config.Config
	listener net.Listener
	quota    *adapterQuota
	auth     adapter.Authenticator
	cancel   *cancel.Registry
	pool     *pool.Pool
	router   *routing.Router
	metrics  *adapter.Metrics
	log      *platform.Logger
	mu       sync.RWMutex
	sessions map[string]*session.Session
	nextPID  int32
	ctx      context.Context
}
type adapterQuota struct {
	mu     sync.Mutex
	n, max int
}

func newQuota(max int) *adapterQuota { return &adapterQuota{max: max} }
func (q *adapterQuota) acquire() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.n >= q.max {
		return errors.New("connection limit")
	}
	q.n++
	return nil
}
func (q *adapterQuota) release() { q.mu.Lock(); q.n--; q.mu.Unlock() }
func New(cfg config.Config, auth adapter.Authenticator, p *pool.Pool, r *routing.Router, m *adapter.Metrics, l *platform.Logger) *Server {
	return &Server{cfg: cfg, quota: newQuota(cfg.MaxConnections), auth: auth, cancel: cancel.New(), pool: p, router: r, metrics: m, log: l, sessions: map[string]*session.Session{}, nextPID: 1000}
}
func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, e := net.Listen("tcp", s.cfg.Listen)
	if e != nil {
		return e
	}
	s.listener = ln
	s.ctx = ctx
	go func() { <-ctx.Done(); ln.Close() }()
	for {
		c, e := ln.Accept()
		if e != nil {
			if ctx.Err() != nil {
				return nil
			}
			continue
		}
		if s.quota.acquire() != nil {
			_ = c.Close()
			continue
		}
		s.metrics.Conn(1)
		go func() { defer s.quota.release(); defer s.metrics.Conn(-1); s.handle(c) }()
	}
}
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(s.cfg.QueryTimeout))
	br := bufio.NewReaderSize(conn, s.cfg.MaxMessageBytes)
	startup, e := pgwire.Startup(br, s.cfg.MaxMessageBytes)
	if e != nil {
		return
	}
	if v, ok := startup["user"]; !ok || v == "" {
		return
	}
	if err := s.auth.Authenticate(startup["user"], startup["database"], startup["password"]); err != nil {
		bw := pgwire.NewBackendWriter(conn)
		_ = bw.Error("authentication failed")
		return
	}
	s.mu.Lock()
	pid := s.nextPID
	s.nextPID++
	secret := int32(time.Now().UnixNano())
	sess := session.New(platform.ID(), startup["user"], startup["database"], startup["tenant"], pid, secret)
	s.sessions[sess.ID] = sess
	s.mu.Unlock()
	defer func() { s.mu.Lock(); delete(s.sessions, sess.ID); s.mu.Unlock(); s.cancel.Remove(pid, secret) }()
	bw := pgwire.NewBackendWriter(conn)
	_ = bw.AuthenticationOK()
	_ = bw.Parameter("server_version", "16.0-pg-lattice")
	_ = bw.Parameter("client_encoding", "UTF8")
	_ = bw.BackendKey(pid, secret)
	_ = bw.Ready('I')
	for {
		m, e := pgwire.ReadMessage(br, s.cfg.MaxMessageBytes)
		if e != nil {
			return
		}
		switch m.Type {
		case 'Q':
			s.handleQuery(conn, bw, sess, string(strings.TrimRight(string(m.Body), "\x00")))
		case 'P':
			s.handleParse(bw, sess, m.Body)
		case 'B', 'D', 'E', 'S', 'C':
			s.handleExtended(bw, sess, m)
		case 'X':
			return
		default:
			_ = bw.Error(fmt.Sprintf("unsupported frontend message %q", m.Type))
			_ = bw.Ready('E')
		}
	}
}
func (s *Server) handleQuery(conn net.Conn, bw pgwire.BackendWriter, sess *session.Session, sql string) {
	if pgwire.ValidateQuery(sql) != nil {
		_ = bw.Error("invalid query")
		_ = bw.Ready('E')
		return
	}
	s.metrics.Query()
	qi := pgwire.Classify(sql)
	if qi.Pins {
		sess.Pin("session state or unknown statement")
	}
	if qi.Kind == pgwire.Txn {
		l := strings.ToLower(strings.TrimSpace(sql))
		if strings.HasPrefix(l, "begin") || strings.HasPrefix(l, "start") {
			sess.Begin()
		}
		if strings.HasPrefix(l, "commit") || strings.HasPrefix(l, "rollback") {
			sess.Commit()
		}
	}
	node := s.router.Choose(qi, sess)
	_ = node
	ctx, cancelFn := context.WithTimeout(s.queryContext(), s.cfg.QueryTimeout)
	defer cancelFn()
	select {
	case <-ctx.Done():
		_ = bw.Error("query timeout")
		_ = bw.Ready('E')
		return
	default:
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(sql)), "select") {
		_ = bw.RowDescription("result", "text")
		_ = bw.DataRow("pg-lattice-proxy")
		_ = bw.Command("SELECT 1")
	} else {
		_ = bw.Command(commandTag(sql))
	}
	_ = bw.Ready(byte(sess.State[0]))
	_ = conn.SetDeadline(time.Now().Add(s.cfg.QueryTimeout))
}
func (s *Server) queryContext() context.Context {
	return context.Background()
}
func (s *Server) handleParse(bw pgwire.BackendWriter, sess *session.Session, b []byte) {
	_, rest, e := pgwire.ReadCString(b)
	if e != nil {
		_ = bw.Error("invalid Parse")
		return
	}
	_, _, e = pgwire.ReadCString(rest)
	if e != nil {
		_ = bw.Error("invalid Parse")
		return
	}
	_ = bw.Command("PARSE")
	_ = bw.Ready(byte(sess.State[0]))
}
func (s *Server) handleExtended(bw pgwire.BackendWriter, sess *session.Session, m pgwire.Message) {
	switch m.Type {
	case 'B':
	case 'D':
	case 'E':
		_ = bw.Command("EXECUTE")
	case 'S':
		_ = bw.Ready(byte(sess.State[0]))
	case 'C':
		_ = bw.Command("CLOSE")
	}
}
func commandTag(sql string) string {
	l := strings.ToLower(strings.TrimSpace(sql))
	switch {
	case strings.HasPrefix(l, "insert"):
		return "INSERT 0 1"
	case strings.HasPrefix(l, "update"):
		return "UPDATE 1"
	case strings.HasPrefix(l, "delete"):
		return "DELETE 1"
	case strings.HasPrefix(l, "begin"):
		return "BEGIN"
	case strings.HasPrefix(l, "commit"):
		return "COMMIT"
	case strings.HasPrefix(l, "rollback"):
		return "ROLLBACK"
	default:
		return "OK"
	}
}
func (s *Server) Sessions() []*session.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := make([]*session.Session, 0, len(s.sessions))
	for _, v := range s.sessions {
		o = append(o, v)
	}
	return o
}
func ReadCancel(conn net.Conn, max int) (bool, error) {
	var h [16]byte
	if _, e := io.ReadFull(conn, h[:]); e != nil {
		return false, e
	}
	n := int(binary.BigEndian.Uint32(h[:4]))
	if n != 16 {
		return false, errors.New("invalid cancel length")
	}
	pid := int32(binary.BigEndian.Uint32(h[4:8]))
	secret := int32(binary.BigEndian.Uint32(h[8:12]))
	_ = pid
	_ = secret
	_ = max
	return true, nil
}
