# DLX repository workflow

This is an independent copy managed by the github-sync skill.

- Work on `custom` or a feature branch based on `custom`; target pull requests at `custom`.
- `origin` is the user repository; `upstream` is OwO-Network/DLX, source branch `main`.
- `master` is a pristine upstream snapshot. Keep it exactly at the pinned source SHA using fast-forward updates only. Never put user changes or sync merge commits there.
- Integrate upstream into custom with normal merges. Preserve existing commits; do not squash, rebase, force-push, reset, or automatically stash user work.
- See `docs/github-sync.md` for migration provenance and CI/release conventions.
- Preserve the default Go TLS transport and conservative Retry-After handling. Run `go test ./...` for translation changes and distinguish live verification from mocks.
- No automatic release, image publication, visibility change, or scheduled sync is implied by a normal source update.

- This installation uses the native macOS binary. Do not reintroduce Docker CI, image builds, or a container runtime requirement unless requested.
