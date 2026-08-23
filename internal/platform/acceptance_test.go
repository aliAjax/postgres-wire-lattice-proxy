package platform_test

import (
	"testing"
	"time"

	"github.com/acme/pg-lattice-proxy/internal/platform"
)

func TestR008RealClockAfterFiresR008(t *testing.T) {
	c := platform.RealClock{}
	select {
	case <-c.After(5 * time.Millisecond):
	case <-time.After(100 * time.Millisecond):
		t.Fatal("clock timer did not fire")
	}
}
