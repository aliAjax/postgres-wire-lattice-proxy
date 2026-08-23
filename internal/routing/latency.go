package routing

import (
	"sort"
	"sync"
	"time"
)

type LatencyWindow struct {
	mu      sync.Mutex
	samples []time.Duration
	max     int
}

func NewLatencyWindow(max int) *LatencyWindow { return &LatencyWindow{max: max} }
func (w *LatencyWindow) Add(v time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.max > 0 && len(w.samples) >= w.max {
		copy(w.samples, w.samples[1:])
		w.samples = w.samples[:w.max-1]
	}
	w.samples = append(w.samples, v)
}
func (w *LatencyWindow) Percentile(p float64) time.Duration {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.samples) == 0 {
		return 0
	}
	v := w.samples
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	idx := int(float64(len(v)-1) * p)
	return v[idx]
}
