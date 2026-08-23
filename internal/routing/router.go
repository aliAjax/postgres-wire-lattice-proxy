package routing

import (
	"github.com/acme/pg-lattice-proxy/internal/pgwire"
	"github.com/acme/pg-lattice-proxy/internal/session"
	"strings"
	"sync"
	"time"
)

type Node struct {
	ID, Addr                    string
	ReadOnly, Draining, Healthy bool
	Latency                     time.Duration
	LastChange                  time.Time
}
type Router struct {
	mu       sync.RWMutex
	primary  Node
	replicas map[string]Node
}

func New(primary Node) *Router      { return &Router{primary: primary, replicas: map[string]Node{}} }
func (r *Router) AddReplica(n Node) { r.mu.Lock(); defer r.mu.Unlock(); r.replicas[n.ID] = n }
func (r *Router) Choose(q pgwire.QueryInfo, s *session.Session) Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// A pinned session (SET/LISTEN/temp/advisory-lock/unknown) carries
	// session-local state that only the primary holds, as does any session
	// mid-transaction. Keep such sessions on the primary.
	if s != nil && (s.Pinned || s.State != session.Idle) {
		return r.primary
	}
	if q.Kind == pgwire.ReadOnly && !q.Pins {
		for _, n := range r.replicas {
			if n.Healthy && !n.Draining {
				return n
			}
		}
	}
	return r.primary
}
func (r *Router) Update(id string, healthy, draining bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if id == r.primary.ID {
		r.primary.Healthy = healthy
		r.primary.Draining = draining
		r.primary.LastChange = time.Now()
		return
	}
	n := r.replicas[id]
	n.Healthy = healthy
	n.Draining = draining
	n.LastChange = time.Now()
	r.replicas[id] = n
}
func (r *Router) Status() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := map[string]any{"primary": r.primary}
	reps := []Node{}
	for _, n := range r.replicas {
		reps = append(reps, n)
	}
	out["replicas"] = reps
	return out
}
func IsTransaction(sql string) bool {
	l := strings.ToLower(strings.TrimSpace(sql))
	return strings.HasPrefix(l, "begin") || strings.HasPrefix(l, "start transaction") || strings.HasPrefix(l, "commit") || strings.HasPrefix(l, "rollback")
}
