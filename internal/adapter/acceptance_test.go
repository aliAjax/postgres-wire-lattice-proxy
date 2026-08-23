package adapter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/adapter"
)

func TestR004ControlHandlerWithEmptyHealthR004(t *testing.T) {
	c := &adapter.Control{Health: nil}
	r := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	defer func() {
		if v := recover(); v != nil {
			t.Fatalf("ready endpoint panicked: %v", v)
		}
	}()
	c.Handler().ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 for missing health registry, got %d", w.Code)
	}
}
