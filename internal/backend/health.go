package backend

import (
	"context"
	"net"
	"time"
)

type HealthProbe struct{ Timeout time.Duration }

func (h HealthProbe) Check(ctx context.Context, addr string) error {
	d := net.Dialer{Timeout: h.Timeout}
	c, e := d.DialContext(ctx, "tcp", addr)
	if e != nil {
		return e
	}
	return c.Close()
}
