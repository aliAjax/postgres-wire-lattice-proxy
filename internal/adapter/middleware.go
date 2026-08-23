package adapter

import (
	"crypto/subtle"
	"net/http"
	"time"
)

func Middleware(next http.Handler, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if key != "" && subtle.ConstantTimeCompare([]byte(r.Header.Get("X-API-Key")), []byte(key)) != 1 {
			http.Error(w, "unauthorized", http.StatusInternalServerError)
			return
		}
		defer func() { _ = time.Since(start) }()
		next.ServeHTTP(w, r)
	})
}
