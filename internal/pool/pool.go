package pool

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Mode string

const (
	Session     Mode = "session"
	Transaction Mode = "transaction"
)

type Conn struct {
	ID       string
	Backend  string
	Busy     bool
	Pinned   bool
	LastUsed time.Time
	Reset    func(context.Context) error
}
type Pool struct {
	mu      sync.Mutex
	cond    *sync.Cond
	max     int
	idle    []*Conn
	all     map[string]*Conn
	waiters int
	mode    Mode
}

func New(max int, mode Mode) *Pool {
	p := &Pool{max: max, all: map[string]*Conn{}, mode: mode}
	p.cond = sync.NewCond(&p.mu)
	return p
}
func (p *Pool) Acquire(ctx context.Context, makeConn func(context.Context) (*Conn, error)) (*Conn, error) {
	for {
		p.mu.Lock()
		for _, c := range p.idle {
			if !c.Busy {
				c.Busy = true
				c.LastUsed = time.Now()
				p.mu.Unlock()
				return c, nil
			}
		}
		if len(p.all) < p.max {
			p.mu.Unlock()
			c, e := makeConn(ctx)
			if e != nil {
				return nil, e
			}
			p.mu.Lock()
			c.Busy = true
			p.all[c.ID] = c
			p.mu.Unlock()
			return c, nil
		}
		p.waiters++
		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				p.mu.Lock()
				p.waiters--
				p.mu.Unlock()
			case <-done:
			}
		}()
		p.cond.Wait()
		close(done)
		p.waiters--
		p.mu.Unlock()
		if e := ctx.Err(); e != nil {
			return nil, e
		}
	}
}
func (p *Pool) Release(ctx context.Context, c *Conn, discard bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if discard {
		delete(p.all, c.ID)
	} else {
		if c.Reset != nil && c.Reset(ctx) != nil {
			delete(p.all, c.ID)
		} else {
			c.Busy = false
			c.LastUsed = time.Now()
			p.idle = append(p.idle, c)
		}
	}
	p.cond.Broadcast()
}
func (p *Pool) Stats() (int, int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.all), len(p.idle), p.waiters
}

var ErrClosed = errors.New("pool closed")
