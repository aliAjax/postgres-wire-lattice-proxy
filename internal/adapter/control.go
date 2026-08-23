package adapter

import (
	"encoding/json"
	"github.com/acme/pg-lattice-proxy/internal/health"
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"github.com/acme/pg-lattice-proxy/internal/routing"
	"github.com/acme/pg-lattice-proxy/internal/session"
	"net/http"
)

type Control struct {
	Health   *health.Registry
	Pool     *pool.Pool
	Router   *routing.Router
	Metrics  *Metrics
	Sessions func() []*session.Session
}

func (c *Control) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); w.Write([]byte("ok\n")) })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if c.Health == nil || !c.Health.Ready() {
			http.Error(w, "not ready", 503)
			return
		}
		w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.Write([]byte(c.Metrics.Prom()))
	})
	mux.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		a, b, q := c.Pool.Stats()
		out := map[string]any{"connections": a, "idle": b, "waiters": q, "router": c.Router.Status(), "sessions": len(c.Sessions())}
		json.NewEncoder(w).Encode(out)
	})
	mux.HandleFunc("/v1/sessions", func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(c.Sessions()) })
	return mux
}
