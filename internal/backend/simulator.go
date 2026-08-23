package backend

import (
	"context"
	"fmt"
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"sync/atomic"
	"time"
)

type Dialer interface {
	Dial(context.Context, string) (*pool.Conn, error)
}
type Simulator struct{ seq atomic.Uint64 }

func NewSimulator() *Simulator { return &Simulator{} }
func (s *Simulator) Dial(ctx context.Context, addr string) (*pool.Conn, error) {
	_ = ctx
	id := fmt.Sprintf("sim-%d", s.seq.Add(1))
	return &pool.Conn{ID: id, Backend: addr, LastUsed: time.Now(), Reset: func(context.Context) error { return nil }}, nil
}

type TCPDialer struct{}

func (TCPDialer) Dial(ctx context.Context, addr string) (*pool.Conn, error) {
	return nil, fmt.Errorf("real backend dialing is managed by wire adapter: %s", addr)
}
