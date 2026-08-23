package tenant

import (
	"errors"
	"sync"
	"time"
)

type Tenant struct {
	ID, Name       string
	MaxConnections int
	CreatedAt      time.Time
}
type Registry struct {
	mu    sync.RWMutex
	items map[string]Tenant
}

func NewRegistry() *Registry { return &Registry{items: map[string]Tenant{}} }
func (r *Registry) Put(t Tenant) error {
	if t.ID == "" {
		return errors.New("tenant id required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[t.ID] = t
	return nil
}
func (r *Registry) Get(id string) (Tenant, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	return t, ok
}
func (r *Registry) List() []Tenant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Tenant, 0, len(r.items))
	for _, t := range r.items {
		out = append(out, t)
	}
	return out
}
