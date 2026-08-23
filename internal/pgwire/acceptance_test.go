package pgwire_test

import (
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/pgwire"
)

func TestR007StateMachineCompletesExtendedFlowR007(t *testing.T) {
	m := pgwire.NewStateMachine()
	if err := m.Accept('P'); err != nil {
		t.Fatal(err)
	}
	m.Complete()
	if got := m.Phase(); got != pgwire.PhaseReady {
		t.Fatalf("extended flow stayed in phase %d", got)
	}
}
