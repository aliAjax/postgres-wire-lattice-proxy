package health

import (
	"sync"
	"time"
)

type Check struct {
	Name      string
	Healthy   bool
	Detail    string
	CheckedAt time.Time
}
type Registry struct {
	mu     sync.RWMutex
	checks map[string]Check
}

func New() *Registry { return &Registry{checks: map[string]Check{}} }
func (r *Registry) Set(name string, ok bool, detail string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[name] = Check{name, ok, detail, time.Now()}
}
func (r *Registry) Ready() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.checks {
		if !c.Healthy {
			return false
		}
	}
	return true
}
func (r *Registry) Snapshot() []Check {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := make([]Check, 0, len(r.checks))
	for _, c := range r.checks {
		o = append(o, c)
	}
	return o
}
