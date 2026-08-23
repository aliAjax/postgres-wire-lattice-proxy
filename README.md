# pg-lattice-proxy

PostgreSQL wire-aware connection proxy implemented in pure Go. The data plane accepts the PostgreSQL startup, simple query, extended query, cancellation and termination messages. It provides bounded connection admission, session/transaction state tracking, conservative SQL routing, and an authenticated control API.

## Run

```bash
go run ./cmd/pg-lattice-proxy -config configs/config.yaml
```

The proxy listens on `127.0.0.1:6432`; the control API listens on `127.0.0.1:8080`. The default backend is an in-memory protocol-compatible simulator, so the project is runnable without PostgreSQL. Set `BACKEND_ADDR` to route to a real PostgreSQL endpoint when available.

## Verification

```bash
go test -race ./...
go vet ./...
go build ./...
printf 'health: '; curl -s http://127.0.0.1:8080/healthz
```

The simulator returns PostgreSQL `ReadyForQuery`, command completion and a single text result for `SELECT` statements. Authentication defaults to trust for local development; production deployments should configure TLS and an API key.
