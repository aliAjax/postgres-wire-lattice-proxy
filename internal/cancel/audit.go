package cancel

import (
	"sync"
	"time"
)

type Event struct {
	PID    int32
	Secret int32
	Found  bool
	At     time.Time
}
type Audit struct {
	mu     sync.Mutex
	events []Event
	max    int
}

func NewAudit(max int) *Audit { return &Audit{max: max} }
func (a *Audit) Add(e Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.events) >= a.max {
		copy(a.events, a.events[1:])
		a.events = a.events[:a.max-1]
	}
	a.events = append(a.events, e)
}
func (a *Audit) List() []Event {
	a.mu.Lock()
	defer a.mu.Unlock()
	o := make([]Event, len(a.events))
	copy(o, a.events)
	return o
}
