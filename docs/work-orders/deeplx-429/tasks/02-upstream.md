# Upstream and Docker
Status: completed; dependencies: none; parallel group A; owner: assigned agent; complexity: medium.
Goal: evidence-backed diagnosis for 429 repair.
Context: $DLX_CHECKOUT; upstream OwO-Network/DLX; existing modified files are user-owned.
Scope: Read-only GitHub current source/issues/releases and Docker availability/image provenance. Do not install/start containers or send live translation requests; report exact reproducible comparison commands.
Write ownership: this task document only. Production fixes belong to coordinator.
Steps: inspect relevant evidence; test hypotheses read-only; report findings with source links/line locations.
Acceptance: distinguish established facts from hypotheses; identify tests and remaining blockers.
Verification: source inspection and appropriate non-mutating commands.
Handoff: files inspected, commands/results, decisions, risks and follow-ups. No secrets.

## Research evidence (2026-11-17)
- https://github.com/OwO-Network/DLX/issues/236 reports v1.2.4 429 caused by stale 26.42 iOS profile and fixed by 26.52 on BOTH missuo/deeplx:v1.2.4 and release Linux binary. This is an issue author's controlled reproduction, not proof of the current local cause.
- https://github.com/OwO-Network/DLX/issues/235 confirms the IP-block diagnosis is DLX's hardcoded wording.
- https://github.com/OwO-Network/DLX/issues/248 is a newer unresolved 429 report (Oct 15); browser translation still works. Comments do not establish a fix.
- https://github.com/OwO-Network/DLX/releases/tag/v1.3.1 is current latest release; local/runtime version resolver 5e7e212 already included.
- GitHub Docker workflow builds the same source into Linux images, with no Docker-specific translation implementation.
- Registry GET: latest and v1.3.1 index sha256:f1af0bce495dedd71708317d098a37d4d23ee8f40fe3c634e72cb2e86e6ad93f; arm64 manifest sha256:679a6d8948eea232b2125741ac370d94ab9da13c50ff3028775fd21c034a5e3f. Image revision 83f97d336d145c56b06fb0316f7187f41954c664, version v1.3.1.

## Authorized Docker verification expansion
Coordinator authorized installing colima/docker and a one-request test. Installed Homebrew bottles colima 0.10.3, lima 2.3.0, docker 29.10.0; Homebrew auto-updated 7.0.12 to 7.0.14. Unrelated pkgconf refresh reported old Xcode, but requested installs succeeded.
Created 2 CPU / 2 GiB VZ arm64 VM, 10 GiB root/data disks, no directory mounts, no SSH config, no Docker context activation. COLIMA_HOME was ignored during startup; stopped new VM and moved newly-created ~/.colima to $DLX_RUNTIME/docker-runtime with a ~/.colima symlink. Restart confirms disks external. Cache at $DLX_RUNTIME/docker-cache.
CLI uses explicit --host unix://$HOME/.colima/default/docker.sock to avoid default context mutation.

## Actual Docker result and cleanup
At 2026-11-17 14:13:18 KST official arm64 image was running as `dlx-429-verification`, bound only to `127.0.0.1:1190`. One POST /translate with `{"text":"こんにちは。","source_lang":"JA","target_lang":"KO"}` returned:

```
HTTP/1.1 429 Too Many Requests
Content-Type: application/json; charset=utf-8
Date: Tue, 17 Nov 2026 05:13:18 GMT
Content-Length: 141

{"code":429,"message":"too many requests, your IP has been blocked by DeepL temporarily, please don't request it frequently in a short time"}
```

Container log also records HTTP 429; no App Store lookup error was logged. Image inspect confirms `arm64`, v1.3.1, source revision 83f97d336d145c56b06fb0316f7187f41954c664 and index digest recorded above. This directly disproves a Docker-only workaround on this Mac/network at test time; it does not prove which upstream policy produced 429. No automatic repeated translation calls were made.

Cleanup verified 14:13:31 KST: test container stopped, `colima list` reports default Stopped, and `lsof -nP -iTCP:1190 -sTCP:LISTEN` finds no listener. Retained stopped container/image/VM for reproducibility; no login service enabled. ~/.colima is a symlink to external docker-runtime (1.9 GiB actual usage); external docker-cache is 320 MiB. CLI tools remain installed.

Reproduce when required:
```sh
XDG_CACHE_HOME=$DLX_RUNTIME/docker-cache colima start --activate=false
docker --host unix://$HOME/.colima/default/docker.sock start dlx-429-verification
curl --max-time 40 -sS -i http://127.0.0.1:1190/translate -H 'Content-Type: application/json' --data '{"text":"こんにちは。","source_lang":"JA","target_lang":"KO"}'
docker --host unix://$HOME/.colima/default/docker.sock stop dlx-429-verification
colima stop
```

Changed files: only this task document. Additional authorized host changes: installed named tools, created external VM/cache and ~/.colima symlink; no DLX source edits. Remaining limit: upstream 429 still exists and Docker does not resolve it.
