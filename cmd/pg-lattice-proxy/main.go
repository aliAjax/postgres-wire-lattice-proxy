package main

import (
	"context"
	"flag"
	"github.com/acme/pg-lattice-proxy/internal/adapter"
	"github.com/acme/pg-lattice-proxy/internal/backend"
	"github.com/acme/pg-lattice-proxy/internal/config"
	"github.com/acme/pg-lattice-proxy/internal/frontend"
	"github.com/acme/pg-lattice-proxy/internal/health"
	"github.com/acme/pg-lattice-proxy/internal/platform"
	"github.com/acme/pg-lattice-proxy/internal/pool"
	"github.com/acme/pg-lattice-proxy/internal/routing"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	path := flag.String("config", "configs/config.yaml", "config path")
	flag.Parse()
	cfg, e := config.Load(*path)
	if e != nil {
		panic(e)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log := platform.NewLogger()
	metrics := &adapter.Metrics{}
	router := routing.New(routing.Node{ID: "primary", Addr: cfg.BackendAddr, Healthy: true})
	router.AddReplica(routing.Node{ID: "replica-1", Addr: "simulator-replica", Healthy: true})
	p := pool.New(cfg.MaxConnections, pool.Transaction)
	sim := backend.NewSimulator()
	_ = sim
	auth := adapter.BuildAuth(cfg.AuthMode, cfg.APIKey)
	srv := frontend.New(cfg, auth, p, router, metrics, log)
	h := health.New()
	h.Set("control", true, "ok")
	control := &adapter.Control{Health: h, Pool: p, Router: router, Metrics: metrics, Sessions: srv.Sessions}
	httpSrv := &http.Server{Addr: cfg.ControlListen, Handler: adapter.Middleware(control.Handler(), cfg.APIKey), ReadHeaderTimeout: 3 * time.Second}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("control server", map[string]any{"error": err.Error()})
		}
	}()
	go func() {
		if err := srv.ListenAndServe(ctx); err != nil {
			log.Error("proxy server", map[string]any{"error": err.Error()})
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	time.Sleep(100 * time.Millisecond)
}
