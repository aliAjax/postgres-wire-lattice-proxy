package adapter

import (
	"fmt"
	"sync/atomic"
)

type Metrics struct {
	connections atomic.Int64
	queries     atomic.Int64
	errors      atomic.Int64
	cancels     atomic.Int64
}

func (m *Metrics) Conn(delta int64) { m.connections.Add(delta) }
func (m *Metrics) Query()           { m.queries.Add(1) }
func (m *Metrics) Error()           { m.errors.Add(1) }
func (m *Metrics) Cancel()          { m.cancels.Add(1) }
func (m *Metrics) Prom() string {
	return fmt.Sprintf("proxy_connections %d\nproxy_queries_total %d\nproxy_errors_total %d\nproxy_cancels_total %d\n", m.connections.Load(), m.queries.Load(), m.errors.Load(), m.cancels.Load())
}
