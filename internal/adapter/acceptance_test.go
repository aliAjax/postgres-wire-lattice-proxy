package adapter_test

import (
	"strings"
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/adapter"
)

func TestR006MetricsErrorCounterVisibleR006(t *testing.T) {
	m := &adapter.Metrics{}
	m.Error()
	if !strings.Contains(m.Prom(), "proxy_errors_total 1") {
		t.Fatal("error counter is missing from metrics output")
	}
}

func TestR006AdapterAuditZeroLimitIsSafeR006(t *testing.T) {
	a := adapter.NewAuditLog(0)
	defer func() {
		if v := recover(); v != nil {
			t.Fatalf("zero-limit adapter audit panicked: %v", v)
		}
	}()
	a.Add(adapter.AuditEvent{Action: "cancel"})
}
