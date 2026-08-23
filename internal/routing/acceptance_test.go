package routing_test

import (
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/pgwire"
	"github.com/acme/pg-lattice-proxy/internal/routing"
	"github.com/acme/pg-lattice-proxy/internal/session"
)

func TestR007PinnedSessionRoutesPrimaryR007(t *testing.T) {
	r := routing.New(routing.Node{ID: "primary", Healthy: true})
	r.AddReplica(routing.Node{ID: "replica", Healthy: true})
	s := session.New("id", "user", "db", "tenant", 1, 2)
	s.Pin("session setting")
	n := r.Choose(pgwire.QueryInfo{Kind: pgwire.ReadOnly}, s)
	if n.ID != "primary" {
		t.Fatalf("pinned session was sent to %q", n.ID)
	}
}

func TestR007VolatilePolicyWhenReplicasDisabledR007(t *testing.T) {
	p := routing.NewPolicy()
	p.SetReplicaEnabled(false)
	if got := p.Classify("select clock_timestamp()"); got.Kind != pgwire.Locking {
		t.Fatalf("volatile query lost locking classification: %q", got.Kind)
	}
}
