package tenant

import (
	"errors"
	"sync"
)

type Quota struct {
	mu                   sync.Mutex
	global, perTenant    map[string]int
	maxGlobal, maxTenant int
}

func NewQuota(g, t int) *Quota {
	return &Quota{global: map[string]int{}, perTenant: map[string]int{}, maxGlobal: g, maxTenant: t}
}
func (q *Quota) Acquire(tenant string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.global["all"] >= q.maxGlobal {
		return errors.New("global connection quota exceeded")
	}
	if q.perTenant[tenant] >= q.maxTenant {
		return errors.New("tenant connection quota exceeded")
	}
	q.global["all"]++
	q.perTenant[tenant]++
	return nil
}
func (q *Quota) Release(tenant string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.global["all"] > 0 {
		q.global["all"]--
	}
	if q.perTenant[tenant] > 0 {
		q.perTenant[tenant]--
		if q.perTenant[tenant] == 0 {
			delete(q.perTenant, tenant)
		}
	}
}
func (q *Quota) Snapshot() (int, map[string]int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	m := map[string]int{}
	for k, v := range q.perTenant {
		m[k] = v
	}
	return q.global["all"], m
}
