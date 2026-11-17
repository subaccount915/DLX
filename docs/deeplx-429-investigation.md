# DeepLX / DLX 429 investigation — 2026-11-17

Published paths use `$DLX_CHECKOUT` (source directory), `$DLX_RUNTIME` (local runtime/data), `$DLX_BUILD` (temporary build tools/cache), and `$DLX_DIAGNOSTICS` (local diagnostic artifacts). Machine-specific paths and binary artifacts are not included.

## Current status: live translation restored

The follow-up investigation found a reproducible transport-specific failure: forcing the dependency's legacy iOS TLS profile returned429, while removing only `SetTLSFingerprintIOS()` returned200 with the same application payload/headers on this Mac. Normal Go TLS is now the default; certificate verification remains enabled. Native single/batch translations and the patched Linux Docker image returned correct Korean translations. Earlier failure results below are retained as historical evidence.

## Findings

This investigation distinguishes an application defect from an upstream rejection. A successful mock translation does not demonstrate that live DeepL translation works.

- Installed native binary: v1.3.1, source revision `83f97d336d145c56b06fb0316f7187f41954c664`, Go 1.27.1 arm64. It already uses oneshot, not the obsolete JSON-RPC endpoint.
- Local checkout: `0f2b50d`; pre-existing user changes already added upstream error detail, Retry-After forwarding, context propagation and retries. These were preserved and reviewed.
- Before this repair, one `JA -> KO` request with `こんにちは。` returned 429 from both the installed native service on 1188 and the locally built service on 1189. The latter reported upstream detail `Too many requests. Please slow down and try again later.` and no Retry-After.
- Apple's public lookup returned DeepL iOS version `26.58` (release date November 13). The existing runtime resolver and prior local diagnostic resolve this same version. Merely updating the fallback version is not an established fix.

## Source and issue evidence

- [Issue #235](https://github.com/OwO-Network/DLX/issues/235): the original IP-block diagnosis is DLX's generic error text, not proof supplied by upstream.
- [Issue #236](https://github.com/OwO-Network/DLX/issues/236): a stale pinned app version caused 429 on both Docker v1.2.4 and native binary. This explains historical incidents, not necessarily today's rejection.
- [Issue #248](https://github.com/OwO-Network/DLX/issues/248): another 429 report remains open; no verified universal fix was established during this review.
- [Client source](https://github.com/OwO-Network/DLX/blob/0f2b50d/translate/translate.go) and [runtime version resolver](https://github.com/OwO-Network/DLX/blob/0f2b50d/translate/appversion.go): Docker packaging does not choose a separate translation backend.

## Repair

The local pre-existing retry loop sent up to three upstream translation requests even without Retry-After. Stop such speculative retries; honor explicit valid retry instructions only within the bounded budget. Preserve subsequent transport failures and avoid sending a request with an already cancelled context. Retain the existing bounded upstream error detail and Retry-After forwarding.

No proxy rotation, credential substitution, alternate translation provider, or forged successful response is used. Live upstream acceptance remains an external dependency.

## Verification and Docker comparison

| Runtime | Version/source | Live JA → KO result |
|---|---|---|
| Installed native, before repair | v1.3.1 / 83f97d3 | 429, generic IP-block wording |
| Local pre-existing modified binary | checkout 0f2b50d + existing edits | 429, upstream detail, no Retry-After |
| Official Docker linux/arm64 | latest = v1.3.1 / 83f97d3 | 429, generic IP-block wording, 0.87 s |
| Rebuilt and deployed native | checkout 0f2b50d + preserved edits + this repair | 429, upstream detail, no Retry-After, about 0.9 s |

Official Docker index digest: `sha256:f1af0bce495dedd71708317d098a37d4d23ee8f40fe3c634e72cb2e86e6ad93f`.
Arm64 child digest: `sha256:679a6d8948eea232b2125741ac370d94ab9da13c50ff3028775fd21c034a5e3f`.
Docker test was isolated on localhost port 1190. Colima and Docker CLI were installed for this test; VM data is on the external SSD, with `~/.colima` linked to `$DLX_RUNTIME/docker-runtime`. The test container and VM are stopped after verification; CLI tools remain installed for reproducibility.

`go test ./...`, `go test -race ./...`, `go build -o dlx .`, and `git diff --check` passed. Coordinator independently reran all package tests and built the binary. Tests cover a successful retry only with an explicit Retry-After, a single attempt without one, invalid header handling, retry budget, cancelled context, and preservation of subsequent transport errors. Mock success is not live translation success.

The patched binary was deployed to `$DLX_RUNTIME/bin/deeplx` and is serving `127.0.0.1:1188`. The source checkout binary is `$DLX_CHECKOUT/dlx`. No startup/login service was added.
Backup: `$DLX_RUNTIME/backups/deeplx-429-20261117-141309/deeplx`.
Runtime log: `$DLX_RUNTIME/logs/deeplx.log`.
For rollback, stop the matching running deeplx process, restore the backup binary to the runtime path, then start it with `-ip 127.0.0.1`.

## Follow-up root cause and successful verification

A recent [fork fix](https://github.com/NotKeKe/DeepLX/commit/8e9df8dd178420308bff93f8fb67f1b38b7e086e) reported the same transport-dependent rejection. Inspecting the installed dependency confirmed `utls HelloIOS_Auto = HelloIOS_14`, despite modern app metadata. Official Chrome extension1.107.0 still uses the same oneshot endpoint/auth/body essentials, so no endpoint migration was justified.

Controlled probes:

1. Existing iOS-fingerprint service: 429 at14:42:59KST; previous service request was14:13:15.
2. Foundation URLSession using current Easydict-style request:200, `안녕하세요.`. This was a clue, not a pure single-variable test because cookie behavior also differed.
3. Go overlay deleting **only** `SetTLSFingerprintIOS()` from current translate.go:200, `안녕하세요.` on port1191. Application payload/headers, dynamic version, cookies, timeout and retry policy unchanged. This isolates the transport configuration as the practical fix; server-side classification internals remain unknown.
4. Production patched native on1188: `今日は天気がいいですね。` -> `오늘 날씨가 참 좋네요.` (200,1.229s); batch `おはようございます。`, empty string, `明日は図書館へ行きます。` -> `좋은 아침입니다.`, empty string, `내일은 도서관에 갈 거예요.` (200,0.911s). Shape and blank alignment preserved.
5. Patched linux/arm64 binary in pinned official Alpine runtime on1190:200, `안녕하세요.`. See [Docker record](work-orders/deeplx-429/tasks/07-docker-patched.md) for image/source hashes and cleanup.
6. Alternative storefront endpoint from pot-desktop PR1290 returned404 and was rejected as a candidate. Official browser UI comparison was unavailable due to CUA ScreenCaptureKit error-3811; it is not claimed as completed.

Production change is default Go TLS directly, with no forced iOS fingerprint and no429-triggered transport/proxy rotation. Prior conservative retry/error handling remains. Unit/race suites and build passed; real translation success was checked separately from mocks.

Final deployed binary: `$DLX_RUNTIME/bin/deeplx`, running on `127.0.0.1:1188`.
Immediately preceding binary backup: `$DLX_RUNTIME/backups/deeplx-default-tls-20261117-144611/deeplx`.
Reproduction artifacts: `$DLX_RUNTIME/diagnostics/dlx-feasibility/` (isolated Go overlay, Swift probe, native live result JSON).

This fixes the persistent429 reproduced in this environment. It does not remove genuine upstream usage limits or guarantee that the unofficial endpoint will never change. No paid API/provider substitution was used.
