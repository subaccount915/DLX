# Independent repository workflow

- Repository: https://github.com/subaccount915/DLX (public, independent).
- `upstream`: https://github.com/OwO-Network/DLX.git, source branch `main`.
- `origin/master`: exact upstream commit, fast-forward only; no user changes.
- `origin/custom`: user development, normal history-preserving merges from master.
- Local checkout and user pull requests target `custom`.

## Preserved migration history

The previous fork is retained at https://github.com/subaccount915/DLX-fork-backup-20261118, including its metadata. Its existing branches and tags were copied by explicit refs to the new repository. The old `main` remains at `a48b8ae7a74bf025f08bb981fd4d3b39bb4989e8`; this commit (including the working TLS fix) is an ancestor of `custom`. No existing commits were squashed, rebased, or rewritten.

The migration's pinned upstream and pristine master are `0f2b50dd69fea3012d95a8d0e7a4ae8dca9ff5b2`.
A complete local Git bundle and ref inventory were verified before migration. No uncommitted changes or unpublished user commits were present.

## Future updates

Fetch upstream and origin, pin the new source SHA, and verify that the existing master is its ancestor. Fast-forward master to that exact source SHA. Integrate it into custom with a normal merge in a separate clean worktree, preserving both prior custom and source ancestry. Test before pushing. Do not force-push, squash, rebase, or merge custom back into master. No scheduled synchronization has been configured.

## CI and publication

Use the native macOS binary for this customized installation. CI continues to build and test source on custom pushes and pull requests targeting custom. The custom Docker workflow has been removed; no Docker image build or publication is configured. Upstream Dockerfile and compose examples are retained for reference. User release tags use `custom-v*` and verify that the tagged commit belongs to custom history. Original upstream `v*` tags are preserved as history and are not republished as user releases. Upstream workflows remain untouched on master.

The original fork remains the backup for repository metadata; Git history transfer does not migrate issues, PRs, Actions logs/secrets, or release assets. No GitHub releases or issues existed in that fork at inspection. Source licenses and copyright notices are preserved.
