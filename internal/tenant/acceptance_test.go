package tenant

import (
	"testing"
	"time"
)

func TestR006LimiterWindowResetR006(t *testing.T) {
	l := NewLimiter()
	l.Set("tenant-a", Limits{MaxQueriesPerMinute: 1})
	if err := l.Allow("tenant-a"); err != nil {
		t.Fatal(err)
	}
	l.mu.Lock()
	l.window = time.Now().Add(-time.Minute)
	l.mu.Unlock()
	if err := l.Allow("tenant-a"); err != nil {
		t.Fatalf("expired rate window did not reset: %v", err)
	}
}
