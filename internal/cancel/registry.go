package cancel

import (
	"context"
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"sync"
)

type Registry struct {
	mu    sync.Mutex
	items map[[2]int32]*pool.Conn
}

func New() *Registry { return &Registry{items: map[[2]int32]*pool.Conn{}} }
func (r *Registry) Add(pid, secret int32, c *pool.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[[2]int32{pid, secret}] = c
}
func (r *Registry) Remove(pid, secret int32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, [2]int32{pid, secret})
}
func (r *Registry) Cancel(ctx context.Context, pid, secret int32) bool {
	r.mu.Lock()
	c, ok := r.items[[2]int32{pid, secret}]
	r.mu.Unlock()
	if !ok {
		return false
	}
	if c.Reset != nil {
		_ = c.Reset(ctx)
	}
	return true
}
func (r *Registry) Size() int { r.mu.Lock(); defer r.mu.Unlock(); return len(r.items) }
