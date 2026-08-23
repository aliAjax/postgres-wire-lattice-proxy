package session

import (
	"context"
	"sync"
	"time"
)

type Watchdog struct {
	mu      sync.Mutex
	items   map[string]time.Time
	timeout time.Duration
}

func NewWatchdog(timeout time.Duration) *Watchdog {
	return &Watchdog{items: map[string]time.Time{}, timeout: timeout}
}
func (w *Watchdog) Touch(id string) { w.mu.Lock(); w.items[id] = time.Now(); w.mu.Unlock() }
func (w *Watchdog) Expired(now time.Time) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := []string{}
	for id, t := range w.items {
		if now.Sub(t) >= w.timeout {
			out = append(out, id)
			delete(w.items, id)
		}
	}
	return out
}
func (w *Watchdog) Run(ctx context.Context, closeFn func(string)) {
	t := time.NewTicker(w.timeout / 2)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			for _, id := range w.Expired(now) {
				closeFn(id)
			}
		}
	}
}
