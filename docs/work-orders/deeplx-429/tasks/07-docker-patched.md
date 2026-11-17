# Patched Docker live verification
Status: completed; owner: Docker research agent; dependency: default TLS source patch; complexity: medium; group C.
Goal: verify current patched source actually translates inside Linux Docker.
Scope: external diagnostics artifacts, isolated Colima VM, test image/container on localhost1190; no source edits or production listener changes.
Steps: build current source linux/arm64 with existing Go cache; use small Alpine runtime/certificate trust; start external Colima; run single short JA->KO request; record exact output/source hash and cleanup container/VM.
Acceptance: real HTTP200 and correct Korean result, normal TLS certificate verification, no retries/proxies/provider substitution; if failed report exact evidence.
Handoff: commands/results/artifacts/cleanup status; write only this task document.

## Verified outcome 2026-11-17 14:46:42 KST
Patched source compiled with Go1.27.2, GOOS=linux GOARCH=arm64 CGO_ENABLED=0, -trimpath. Used existing GOCACHE=$DLX_BUILD/cache and GOMODCACHE=$DLX_BUILD/modules. No source changes by this agent. Root subsequently changed comments only; recorded hashes describe the compile-time snapshot.

Runtime test reuses the exact official v1.3.1 Alpine runtime layers (pinned index sha256:f1af0bce495dedd71708317d098a37d4d23ee8f40fe3c634e72cb2e86e6ad93f), replacing /app/deeplx with the patched Linux binary. This avoids changing runtime environment as a confounder. Verified Alpine3.24.2 and nonempty CA trust /etc/ssl/certs/ca-certificates.crt (181724 bytes). Normal certificate validation retained; no proxies or credentials.

Artifacts: $DLX_RUNTIME/diagnostics/dlx-feasibility/docker-patched/ contains Dockerfile, deeplx, source-sha256.json, live-response.txt.
- Binary SHA256: a1970613979b37aa0c15ad817f191eba0d67a0d3a02d06f71f4918ac9d2b1aa5
- Source manifest SHA256: e785a1a1cefc69e04e8587c690af5b939517f6e20fb955a3336793d5cdc6e729
- Image tag: dlx-patched-429:20261117
- Image ID: sha256:28b2d35ea66924c204584853dd6e873aa7d121fcb87f9f83a2e127410486aa35
- Build metadata: vcs.revision=0f2b50dd69fea3012d95a8d0e7a4ae8dca9ff5b2, vcs.modified=true. Inherited base labels remain upstream v1.3.1, so use local image ID and binary hash as patched identity.

Test container dlx-patched-429-verification bound only 127.0.0.1:1190. Single POST /translate with {"text":"こんにちは。","source_lang":"JA","target_lang":"KO"}:
```
HTTP/1.1 200 OK
{"alternatives":null,"code":200,"data":"안녕하세요.","id":1794894401814,"method":"Free","source_lang":"JA","target_lang":"KO"}
```
Container log confirms HTTP200 in1.2096s. Official original image returned429 on the same Mac/network earlier; patched binary succeeds in its runtime. No additional translation request was made by this agent.

Cleanup: test container stopped; Colima VM stopped; colima list confirms Stopped and lsof finds no1190 listener. Port1188 untouched. Stopped images/container and external diagnostics artifacts retained for reproducibility.
