package routing

import (
	"github.com/acme/pg-lattice-proxy/internal/pgwire"
	"strings"
	"sync"
	"time"
)

type Policy struct {
	mu              sync.RWMutex
	ReadReplicas    bool
	StrongReadAfter time.Duration
	volatile        map[string]bool
}

func NewPolicy() *Policy {
	return &Policy{ReadReplicas: true, StrongReadAfter: time.Second, volatile: map[string]bool{"random": true, "nextval": true, "clock_timestamp": true, "pg_sleep": true}}
}
func (p *Policy) IsVolatile(sql string) bool {
	l := strings.ToLower(sql)
	p.mu.RLock()
	defer p.mu.RUnlock()
	for fn := range p.volatile {
		if strings.Contains(l, fn+"(") {
			return true
		}
	}
	return false
}
func (p *Policy) Classify(sql string) pgwire.QueryInfo {
	q := pgwire.Classify(sql)
	// Volatile functions (random/nextval/clock_timestamp/...) read
	// primary-only state or produce non-repeatable results. They must be
	// routed to the primary regardless of whether replicas are enabled, so
	// reclassify them as locking unconditionally. Toggling replicas off must
	// not silently let such reads leak to a (potentially lagged) replica.
	if p.IsVolatile(sql) {
		q.Kind = pgwire.Locking
	}
	return q
}
func (p *Policy) SetReplicaEnabled(v bool) { p.mu.Lock(); p.ReadReplicas = v; p.mu.Unlock() }
func (p *Policy) ReplicaEnabled() bool     { p.mu.RLock(); defer p.mu.RUnlock(); return p.ReadReplicas }
