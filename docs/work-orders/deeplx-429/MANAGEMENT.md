# DeepLX 429 repair

**Current outcome: RESOLVED in this environment.** Default Go TLS restored real native and Docker JA->KO translation. Historical partial conclusions below are superseded by the final follow-up results.
Objective: determine actual 429 cause, compare Docker, repair verified defects and validate translation.
Scope: DLX source/runtime; preserve existing edits. No paid credentials or rate-limit evasion.
Architecture: Gin HTTP service calls shared oneshot transport. Native binary and Docker should use the same code.
Plan review: independent read-only source and upstream/Docker investigations precede coordinator-owned edits and runtime integration.

| ID | Task | Status | Owner | Dependencies | Group |
|---|---|---|---|---|---|
| 01 | [Source audit](tasks/01-source.md) | completed | source agent | none | A |
| 02 | [Upstream and Docker](tasks/02-upstream.md) | completed | research agent | none | A |
| 03 | Repair and integration | completed | coordinator + source agent | source audit findings | B |

Integration: 01 + 02 -> 03. Coordinator owns production edits, shared contracts and this document.
Verification: unit/race tests for defects; native live short JA->KO request; Docker same-network request when runtime available; provenance/version comparison.
Risks: upstream rate limits cannot be guaranteed removable; Docker may be unavailable; pre-existing edits must be retained.
Completion: [ ] source audit [ ] upstream evidence [ ] Docker assessment [ ] fix/tests [ ] live validation.
Progress: initialized; discovered pre-existing retry edits and no docker/go in PATH.

Decision update: source agent owns only translate/translate.go and translate/retry_test.go for conservative retry fixes. Coordinator owns docs, build/deployment and verification. Research agent owns Docker installation/testing (Colima VM disk external), no production edits. No overlapping writes.
Live evidence: installed v1.3.1 and pre-existing modified source binary both return 429 for one JA->KO greeting. Official App Store currently resolves 26.58; earlier diagnostic and current source support runtime resolution. A stale pinned app version alone does not explain this environment.

Final integration evidence: conservative retry repair reviewed and tested; native built/deployed after backing up original; actual patched request still 429. Official Docker latest/v1.3.1 also returned 429. No packaging-based recovery demonstrated.
Verification checklist: [x] source audit [x] upstream evidence [x] actual Docker assessment [x] fix/tests [x] live validation [ ] successful live translation (blocked by upstream rejection of unconfirmed cause).
User-visible outcome is PARTIAL: proven local request-amplification/error-masking defects repaired and deployed; translation restoration unresolved. See [full investigation](../../deeplx-429-investigation.md). No unrequested provider substitution or proxy rotation.

Docker cleanup verified by research agent: test container/VM stopped and port 1190 not listening; external image/VM retained for reproduction.

## Follow-up: feasibility of live recovery
Plan reviewed: compare official browser behavior on this network, research verified upstream/client changes, and audit protocol evidence independently. Only evidence-backed experiments, no blind profile rotation. Source changes require a live-success candidate plus regression checks; production remains unchanged until then.
04 (related-client research) and 05 (protocol audit) are independent read-only investigations; 06 (coordinator browser comparison and integration) combines them. Status: in_progress. Verification: official web JA->KO outcome, version-specific source evidence, bounded DLX smoke only for credible fixes.

Recovery candidate CONFIRMED on native: Foundation probe200; DLX overlay deleting only SetTLSFingerprintIOS returned200 and 안녕하세요. Source endpoint/auth remains current according to official Chrome extension1.107.0. Coordinator now owns production translate.go change and native deployment; Docker agent owns isolated Linux build/test, no production source writes. Final acceptance requires native single/batch and Docker patched live smoke plus unit/race tests.
07 Docker validation: ready, depends on coordinator source change (complete), owner research agent, group C.

## Final follow-up completion
04 related-client research: completed;05 protocol audit: completed;06 native integration: completed;07 actual patched Docker: completed; container/VM stopped and1190 listener removed.
Decision: use default Go TLS directly; remove forced legacy iOS fingerprint. Do not adopt fork proxy rotation or automatic429 fallback.
Evidence: controlled one-line overlay200; production1188 single/batch200 including blank alignment; Linux Docker200; go test-race all passed; build and diff checks passed.
[x] source-backed root cause [x] production patch [x] native live restoration [x] Docker live restoration [x] tests [x] original binaries backed up. Browser UI comparison blocked by capture API but no longer necessary to establish the successful controlled fix.
