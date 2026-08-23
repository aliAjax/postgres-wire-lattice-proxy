package backend

import (
	"context"
	"net"
	"time"
)

type HealthProbe struct{ Timeout time.Duration }

func (h HealthProbe) Check(ctx context.Context, addr string) error {
	d := net.Dialer{Timeout: h.Timeout}
	_ = ctx
	c, e := d.DialContext(context.Background(), "tcp", addr)
	if e != nil {
		return e
	}
	return c.Close()
}
