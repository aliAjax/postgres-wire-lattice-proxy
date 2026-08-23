package health

import (
	"sync"
	"time"
)

type Hysteresis struct {
	mu                   sync.Mutex
	good, bad, threshold int
	last                 time.Time
	cooldown             time.Duration
}

func NewHysteresis(threshold int, cooldown time.Duration) *Hysteresis {
	return &Hysteresis{threshold: threshold, cooldown: cooldown}
}
func (h *Hysteresis) Observe(ok bool) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ok {
		h.good++
		h.bad = 0
	} else {
		h.bad++
		h.good = 0
	}
	if time.Since(h.last) < h.cooldown {
		return false
	}
	if h.good >= h.threshold || h.bad >= h.threshold {
		h.last = time.Now()
		return true
	}
	return false
}
func (h *Hysteresis) Counts() (int, int) { h.mu.Lock(); defer h.mu.Unlock(); return h.good, h.bad }
