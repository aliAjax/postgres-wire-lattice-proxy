package tenant_test

import (
	"sync"
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/tenant"
)

func TestR008TenantQuotaConcurrentAcquireReleaseR008(t *testing.T) {
	q := tenant.NewQuota(100, 100)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if err := q.Acquire("tenant-a"); err == nil {
				q.Release("tenant-a")
			}
		}()
	}
	close(start)
	wg.Wait()
	global, perTenant := q.Snapshot()
	if global != 0 || perTenant["tenant-a"] != 0 {
		t.Fatalf("quota leaked after concurrent lifecycle: global=%d tenant=%d", global, perTenant["tenant-a"])
	}
}
