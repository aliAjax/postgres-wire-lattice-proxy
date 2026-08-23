package session_test

import (
	"testing"

	"github.com/acme/pg-lattice-proxy/internal/session"
)

func TestR007CommittedSessionIsIdleR007(t *testing.T) {
	s := session.New("id", "user", "db", "tenant", 1, 2)
	s.Begin()
	s.Commit()
	if s.State != session.Idle {
		t.Fatalf("commit left session in %q", s.State)
	}
}
