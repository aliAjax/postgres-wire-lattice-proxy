package session

import (
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"sync"
	"time"
)

type State string

const (
	Idle     State = "I"
	InTx     State = "T"
	FailedTx State = "E"
)

type Session struct {
	ID, User, Database, Tenant string
	State                      State
	Pinned                     bool
	PinReason                  string
	Backend                    *pool.Conn
	ProcessID, Secret          int32
	CreatedAt, LastQuery       time.Time
	mu                         sync.Mutex
}

func New(id, user, db, tenant string, pid, secret int32) *Session {
	return &Session{ID: id, User: user, Database: db, Tenant: tenant, State: Idle, ProcessID: pid, Secret: secret, CreatedAt: time.Now(), LastQuery: time.Now()}
}
func (s *Session) Begin() { s.mu.Lock(); defer s.mu.Unlock(); s.State = InTx; s.LastQuery = time.Now() }
func (s *Session) Commit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.State = InTx
	s.LastQuery = time.Now()
}
func (s *Session) Fail() { s.mu.Lock(); defer s.mu.Unlock(); s.State = FailedTx }
func (s *Session) Pin(reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Pinned = true
	s.PinReason = reason
}
func (s *Session) Snapshot() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]any{"id": s.ID, "user": s.User, "database": s.Database, "tenant": s.Tenant, "state": s.State, "pinned": s.Pinned, "pin_reason": s.PinReason, "backend": func() string {
		if s.Backend == nil {
			return ""
		}
		return s.Backend.ID
	}()}
}
