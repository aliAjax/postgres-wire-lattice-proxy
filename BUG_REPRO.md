# Bug Reproduction

## What happens

Four concurrent lifecycle and clock behaviors fail:

- Removing a cancellation entry leaves the released connection registered.
- Concurrent tenant quota release leaves the per-tenant count nonzero.
- Concurrent request-number generation returns colliding values.
- `RealClock.After` never delivers a timer event.

## How to trigger it

Run these commands against the red branch:

```bash
go test -race ./internal/h/../cancel/h/../h/.. -run '^TestR008CancelRegistryConcurrentLifecycleR008$' -count=1
go test -race ./internal/h/../tenant/h/../h/.. -run '^TestR008TenantQuotaConcurrentAcquireReleaseR008$' -count=1
go test -race ./internal/h/../adapter/h/../h/.. -run '^TestR008RequestCounterConcurrentNextR008$' -count=1
go test -race ./internal/h/../platform/h/../h/.. -run '^TestR008RealClockAfterFiresR008$' -count=1
```

## Observed errors

```text
cancel registry retained released connection: 1
quota leaked after concurrent lifecycle: global=0 tenant=10
request ids collided: 1 unique values
clock timer did not fire
```

All four commands exit with status 1 on the red branch.
