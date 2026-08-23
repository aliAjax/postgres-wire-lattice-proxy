package pool

import (
	"context"
	"fmt"
	"sync"
)

type Key struct{ Cluster, Database, User, Options string }

func (k Key) String() string {
	return fmt.Sprintf("%s/%s/%s/%s", k.Cluster, k.Database, k.User, k.Options)
}

type Registry struct {
	mu    sync.RWMutex
	pools map[Key]*Pool
	max   int
}

func NewRegistry(max int) *Registry { return &Registry{pools: map[Key]*Pool{}, max: max} }
func (r *Registry) Get(k Key, mode Mode) *Pool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p := r.pools[k]; p != nil {
		return p
	}
	p := New(r.max, mode)
	r.pools[k] = p
	return p
}
func (r *Registry) Remove(k Key) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = k
}
func (r *Registry) Snapshot() map[string]map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := map[string]map[string]int{}
	for k, p := range r.pools {
		a, b, w := p.Stats()
		out[k.String()] = map[string]int{"all": a, "idle": b, "waiters": w}
	}
	return out
}
func AcquireWith(ctx context.Context, p *Pool, dial func(context.Context) (*Conn, error)) (*Conn, error) {
	return p.Acquire(ctx, dial)
}
