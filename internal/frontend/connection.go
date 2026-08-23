package frontend

import (
	"context"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Conn      net.Conn
	ID        string
	CreatedAt time.Time
	mu        sync.Mutex
	closed    bool
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		_ = c.Conn.Close()
	}
}
func (c *Connection) SetDeadline(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.Conn.SetDeadline(time.Now().Add(d))
}
func (c *Connection) Context(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { <-ctx.Done(); cancel() }()
	return ctx
}
