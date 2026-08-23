package health_test

import (
	"context"
	"testing"
	"time"

	"github.com/acme/pg-lattice-proxy/internal/health"
)

func TestR004DefaultHealthRegistryAcceptsFirstUpdateR004(t *testing.T) {
	var value health.Registry
	r := &value
	r.Set("control", true, "ok")
	if !r.Ready() {
		t.Fatal("healthy first update was not ready")
	}
}

func TestR004HealthProbeOnZeroRegistryR004(t *testing.T) {
	var value health.Registry
	r := &value
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	health.Run(ctx, r, "backend", func(context.Context) error { return nil }, time.Millisecond)
	deadline := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(deadline) {
		if len(r.Snapshot()) == 1 {
			if !r.Ready() {
				t.Fatal("successful probe left registry unhealthy")
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("probe did not publish health state")
}

func TestR004HysteresisOpensAtThresholdR004(t *testing.T) {
	h := health.NewHysteresis(1, 0)
	if !h.Observe(true) {
		t.Fatal("threshold observation did not open state")
	}
}
