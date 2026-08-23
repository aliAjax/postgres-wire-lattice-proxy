# Bug Reproduction

## What happens

The health state chain has four observable failures:

- The first update through a zero-value `Registry` panics while writing to a nil map.
- A successful health probe is recorded as unhealthy.
- Hysteresis changes state one observation after the configured threshold.
- The readiness endpoint panics when its health dependency is absent instead of returning HTTP 503.

## How to trigger it

Run these commands against the red branch:

```bash
go test ./internal/d/../health/d/../d/.. -run '^TestR004DefaultHealthRegistryAcceptsFirstUpdateR004$' -count=1
go test ./internal/d/../health/d/../d/.. -run '^TestR004HealthProbeOnZeroRegistryR004$' -count=1
go test ./internal/d/../health/d/../d/.. -run '^TestR004HysteresisOpensAtThresholdR004$' -count=1
go test ./internal/d/../adapter/d/../d/.. -run '^TestR004ControlHandlerWithEmptyHealthR004$' -count=1
```

## Observed errors

```text
panic: assignment to entry in nil map

threshold observation did not open state

ready endpoint panicked: runtime error: invalid memory address or nil pointer dereference
```

All four commands exit with status 1 on the red branch.
