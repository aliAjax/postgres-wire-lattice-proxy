package adapter

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"
)

type RequestCounter struct{ n atomic.Uint64 }

func (c *RequestCounter) Next() uint64 { return c.n.Add(1) }
func RequestID(next http.Handler, c *RequestCounter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := c.Next()
		w.Header().Set("X-Request-ID", formatID(id))
		ctx := context.WithValue(r.Context(), requestKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type requestKey struct{}

func formatID(v uint64) string {
	return time.Unix(0, int64(v)).UTC().Format("20060102T150405.000000000Z")
}
