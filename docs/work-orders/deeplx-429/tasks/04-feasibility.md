# related-client research
Status: completed; owner: assigned agent; dependencies: none; group B; complexity: medium.
Goal: determine evidence-backed route to live DLX recovery.
Scope: Investigate recent DLX issues/PRs/forks and Easydict DeepL implementation for verified fixes after October 2026. Use actual sources, not repeated guesses. No live translation calls or production edits.
Write scope: this task document only. Preserve all existing edits.
Acceptance: report exact sources/lines, confirmed facts vs hypotheses, commands/results, proposed bounded tests and blockers. No claim of a working fix without real translation.

## Findings 2026-11-17
- Latest DLX PRs end October 10; #238 is the existing dynamic version resolver, later #241/#243/#245/#246 concern batch/UTF16/blank/null behavior, not a new 429 fix. Queried GitHub REST via gh because plugin search returned issues despite is:pr qualifier and fetch collection sort was not reliable.
- Easydict PR https://github.com/tisfeng/Easydict/pull/1341 (merged Sept29) shares dynamic marketing version in body and User-Agent, same fallback26.52/build5443737 as DLX. Latest DeepL-path commit Oct24 24d0b85 fixes concurrency, not transport/profile.
- https://github.com/tisfeng/Easydict/issues/1388 (Nov3) reports an outdated-version message; maintainer comment https://github.com/tisfeng/Easydict/issues/1388#issuecomment-6601251272 says v2.26.0 works. This is historical third-party evidence, not proof of current same-network success.
- Current Easydict source https://github.com/tisfeng/Easydict/blob/dev/Easydict/Swift/Service/DeepL/DeepLService%2BTranslate.swift blob b3c4d18a7941954cd6400268701ac0ef0d75b532: lines13-17 endpoint/OS26.0/build5443737; 85-125 payload and headers; 356-357 User-Agent. Uses Alamofire/Foundation transport, not emulated iOS uTLS. Request field schema is DeepLTranslateResponse.swift.
- Reviewed newest100 DLX forks for pushes after Oct11. tankle0/DLX only compose image changes; ysm827 and fyqh8040 match upstream core. NotKeKe/DeepLX has a relevant real patch below. This is a bounded scan, not exhaustive proof that no other fixes exist.

## Concrete candidate
https://github.com/NotKeKe/DeepLX/commit/8e9df8dd178420308bff93f8fb67f1b38b7e086e (Oct23) reports identical request failing with uTLS iOS but succeeding with Go default TLS on some networks. Diff makes SetTLSFingerprintIOS conditional, caches clients by proxy+TLS profile, and falls back after429; also fallback version26.55. Fork also contains proxy rotation, which is unrelated to the proposed diagnostic and should not be copied. The commit supplies source-level evidence and an author report, not a reproduced result here.

Proposed finite test: one default Go TLS request using existing payload/current version and same network; compare with already measured iOS profile. No proxies, no automated 429 fallback loop, no production modification until demonstrated. Alternative single Foundation URLSession control uses Easydict's known wire shape.

## Coordinator-authorized probe preparation
Created $DLX_RUNTIME/diagnostics/dlx-feasibility/native-probe.swift only; it uses version26.58 supplied by coordinator, same Easydict fields, 20s timeout, no cookies/cache, redirects disabled, one POST, only status/selected errors/translations printed. No live run by this agent.
Default CLT Swift failed to load standard library for macosx27.2; explicit target succeeds:
```sh
/Library/Developer/CommandLineTools/usr/bin/swiftc -target arm64-apple-macosx14.0 -sdk /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk -typecheck $DLX_RUNTIME/diagnostics/dlx-feasibility/native-probe.swift
```
Root owns execution and acceptance. No production edits. No working fix claimed by this investigation.
