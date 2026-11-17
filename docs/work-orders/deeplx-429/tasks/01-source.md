# Source audit
Status: completed; dependencies: none; parallel group A; owner: assigned agent; complexity: medium.
Goal: evidence-backed diagnosis for 429 repair.
Context: $DLX_CHECKOUT; upstream OwO-Network/DLX; existing modified files are user-owned.
Scope: Read-only DLX source and existing diff; determine retry/rate-limit correctness and concrete minimal repair. Inspect tests. Do not send live translation requests.
Write ownership: this task document only. Production fixes belong to coordinator.
Steps: inspect relevant evidence; test hypotheses read-only; report findings with source links/line locations.
Acceptance: distinguish established facts from hypotheses; identify tests and remaining blockers.
Verification: source inspection and appropriate non-mutating commands.
Handoff: files inspected, commands/results, decisions, risks and follow-ups. No secrets.

## Handoff
Coordinator authorized narrow implementation after audit; write scope expanded to `translate/translate.go` and `translate/retry_test.go` only. Existing edits retained.

Established findings:
- HEAD already uses oneshot and runtime App Store version lookup; old hardcoded IP-ban error was an unsupported attribution.
- Pre-existing retry code multiplied each bare 429 into three upstream calls and masked subsequent transport errors as 429.
- Docker release workflow publishes latest on tags; checkout and image version must be compared explicitly.
- External installed binary contains oneshot endpoint; cannot assume it is legacy JSONRPC from version memory.

Implemented: retry only for valid Retry-After; keep underlying transport error after earlier 429; reject already-cancelled contexts before post. Kept bounded explicit retries, existing diagnostic and header forwarding changes.
Tests: `go test ./...` and `go test -race ./...` passed with Go 1.27.2 at `$DLX_BUILD/go/bin/go`, GOCACHE `$DLX_BUILD/cache`, GOMODCACHE `$DLX_BUILD/modules`.
Coverage: absent/malformed headers one attempt; explicit zero recovery and retry cap; longer-than-budget no retry; HTTP-date and overflow parsing; original network error retained; cancelled-before-post sends nothing.
Risks/follow-ups: no shared cooldown across separate requests; version lookup ignores caller cancellation and explicit PROXY and can precede 20-second request budget by 3 seconds. No live translation requests made by this agent; coordinator owns integration/live evidence. Changes mitigate retry amplification, cannot promise upstream 429 removal.
