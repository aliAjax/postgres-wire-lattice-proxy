# Bug Reproduction

## What happens

Four zero-value and state boundaries fail:

- A zero-capacity cancellation audit panics with a slice-bounds error.
- The Prometheus output omits the accumulated error counter.
- A zero-capacity adapter audit panics with the same slice-bounds error.
- An expired tenant rate window retains its old count and rejects the first request in the new window.

## How to trigger it

Run these commands against the red branch:

```bash
go test ./internal/f/../cancel/f/../f/.. -run '^TestR006CancelAuditZeroLimitIsSafeR006$' -count=1
go test ./internal/f/../adapter/f/../f/.. -run '^TestR006MetricsErrorCounterVisibleR006$' -count=1
go test ./internal/f/../adapter/f/../f/.. -run '^TestR006AdapterAuditZeroLimitIsSafeR006$' -count=1
go test ./internal/f/../tenant/f/../f/.. -run '^TestR006LimiterWindowResetR006$' -count=1
```

## Observed errors

```text
zero-limit audit panicked: runtime error: slice bounds out of range [1:0]
error counter is missing from metrics output
zero-limit adapter audit panicked: runtime error: slice bounds out of range [1:0]
expired rate window did not reset: tenant query rate exceeded
```

All four commands exit with status 1 on the red branch.
