package backend

import (
	"context"
	"fmt"
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"sync"
	"time"
)

type Connector struct {
	dialer Dialer
	mu     sync.Mutex
	opened int
	closed int
}

func NewConnector(d Dialer) *Connector { return &Connector{dialer: d} }
func (c *Connector) Open(ctx context.Context, addr string) (*pool.Conn, error) {
	v, e := c.dialer.Dial(ctx, addr)
	if e != nil {
		return nil, fmt.Errorf("dial backend: %w", e)
	}
	c.mu.Lock()
	c.opened++
	c.mu.Unlock()
	return v, nil
}
func (c *Connector) Stats() (int, int) { c.mu.Lock(); defer c.mu.Unlock(); return c.opened, c.closed }
func (c *Connector) Close(v *pool.Conn) {
	if v == nil {
		return
	}
	if v.Reset != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_ = v.Reset(ctx)
		cancel()
	}
	c.mu.Lock()
	c.closed++
	c.mu.Unlock()
}
