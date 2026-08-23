package cancel_test

import (
	"context"
	"sync"
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/cancel"
	"github.com/acme/pg-lattice-proxy/internal/pool"
)

func TestR008CancelRegistryConcurrentLifecycleR008(t *testing.T) {
	r := cancel.New()
	conn := &pool.Conn{ID: "c", Reset: func(context.Context) error { return nil }}
	r.Add(7, 9, conn)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); <-start; r.Cancel(context.Background(), 7, 9) }()
	go func() { defer wg.Done(); <-start; r.Remove(7, 9) }()
	close(start)
	wg.Wait()
	if got := r.Size(); got != 0 {
		t.Fatalf("cancel registry retained released connection: %d", got)
	}
}
