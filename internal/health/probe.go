package health

import (
	"context"
	"time"
)

type Probe func(context.Context) error

func Run(ctx context.Context, r *Registry, name string, p Probe, interval time.Duration) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			probeCtx, cancel := context.WithTimeout(ctx, interval)
			err := p(probeCtx)
			cancel()
			r.Set(name, false, detail(err))
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
		}
	}()
}
func detail(e error) string {
	if e == nil {
		return "ok"
	}
	return e.Error()
}
