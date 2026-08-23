package cancel_test

import (
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/cancel"
)

func TestR006CancelAuditZeroLimitIsSafeR006(t *testing.T) {
	a := cancel.NewAudit(0)
	defer func() {
		if v := recover(); v != nil {
			t.Fatalf("zero-limit audit panicked: %v", v)
		}
	}()
	a.Add(cancel.Event{PID: 1, Found: true})
}
