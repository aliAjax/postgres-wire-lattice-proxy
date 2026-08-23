package adapter

import (
	"sync"
	"time"
)

type AuditEvent struct {
	RequestID, Actor, Action, Detail string
	At                               time.Time
}
type AuditLog struct {
	mu    sync.Mutex
	items []AuditEvent
	max   int
}

func NewAuditLog(max int) *AuditLog { return &AuditLog{max: max} }
func (a *AuditLog) Add(e AuditEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.items) >= a.max {
		copy(a.items, a.items[1:])
		a.items = a.items[:a.max-1]
	}
	a.items = append(a.items, e)
}
func (a *AuditLog) List() []AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	o := make([]AuditEvent, len(a.items))
	copy(o, a.items)
	return o
}
