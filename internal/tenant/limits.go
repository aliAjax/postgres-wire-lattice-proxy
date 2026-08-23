package tenant

import (
	"errors"
	"sync"
	"time"
)

type Limits struct {
	MaxConnections, MaxQueriesPerMinute int
	QueryTimeout, SessionTimeout        time.Duration
}
type Limiter struct {
	mu     sync.Mutex
	limits map[string]Limits
	counts map[string]int
	window time.Time
}

func NewLimiter() *Limiter {
	return &Limiter{limits: map[string]Limits{}, counts: map[string]int{}, window: time.Now()}
}
func (l *Limiter) Set(id string, v Limits) { l.mu.Lock(); l.limits[id] = v; l.mu.Unlock() }
func (l *Limiter) Allow(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.window) >= time.Minute {
		l.counts = map[string]int{}
		l.window = time.Now()
	}
	lim := l.limits[id]
	if lim.MaxQueriesPerMinute > 0 && l.counts[id] >= lim.MaxQueriesPerMinute {
		return errors.New("tenant query rate exceeded")
	}
	l.counts[id]++
	return nil
}
