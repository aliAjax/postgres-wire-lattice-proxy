package adapter_test

import (
	"sync"
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/adapter"
)

func TestR008RequestCounterConcurrentNextR008(t *testing.T) {
	c := &adapter.RequestCounter{}
	start := make(chan struct{})
	values := make(chan uint64, 100)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			values <- c.Next()
		}()
	}
	close(start)
	wg.Wait()
	seen := map[uint64]bool{}
	for i := 0; i < 100; i++ {
		seen[<-values] = true
	}
	if len(seen) != 100 {
		t.Fatalf("request ids collided: %d unique values", len(seen))
	}
}
