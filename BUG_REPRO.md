# Bug Reproduction

## What happens

Four protocol and routing state transitions fail:

- Completing an extended-query flow leaves the protocol state in phase 4 instead of Ready.
- Committing a transaction leaves the session in transaction state instead of Idle.
- A pinned session can route a read to a replica.
- With replicas disabled, a volatile query loses its primary-only classification.

## How to trigger it

Run these commands against the red branch:

```bash
go test ./internal/g/../pgwire/g/../g/.. -run '^TestR007StateMachineCompletesExtendedFlowR007$' -count=1
go test ./internal/g/../session/g/../g/.. -run '^TestR007CommittedSessionIsIdleR007$' -count=1
go test ./internal/g/../routing/g/../g/.. -run '^TestR007PinnedSessionRoutesPrimaryR007$' -count=1
go test ./internal/g/../routing/g/../g/.. -run '^TestR007VolatilePolicyWhenReplicasDisabledR007$' -count=1
```

## Observed errors

```text
extended flow stayed in phase 4
commit left session in "T"
pinned session was sent to "replica"
volatile query lost locking classification: "read"
```

All four commands exit with status 1 on the red branch.
