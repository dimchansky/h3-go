# Releasing

The complete, durable procedure for cutting an h3-go release. Versioning
policy (tag syntax, the two independent version axes, release titles, the
release-note outline) is fixed in [versioning.md](versioning.md); CI tiers
and what runs when are in [ci-policy.md](ci-policy.md). This page is the
operational runbook: gates, commands, and failure handling.

Terms: "rc" is a genuine `vX.Y.Z-rc.N` prerelease of the module version
(rc counting starts at 1); the final tag must point at the exact commit the
last rc validated.

## Invariants (never violated)

- **Published tags are immutable.** A bad release is fixed forward with a
  new version and, when warranted, a Go module `retract` directive — never
  by moving, deleting, or reusing a tag ([versioning.md](versioning.md)).
- **Release assets are immutable.** Releases are published with GitHub
  release immutability enabled; a bad asset means a corrective release,
  never `gh release upload --clobber` or asset edits.
- **The final tag is never the first validation of anything.** Every gate
  runs at rc first; any post-rc commit — code or docs — requires a new rc.
- **No direct pushes to `master`**; every change lands through a PR with
  the `CI / required` check green. The repository-admin ruleset bypass
  exists for emergencies only; using it requires an explanatory note in the
  next PR description.

## Toolchain pin

Releases build with one exact Go toolchain, pinned in two places that must
stay in lockstep:

- `requiredGoVersion` in [tools/releasepack/main.go](../tools/releasepack/main.go)
  (enforced — the build refuses any other toolchain), and
- `RELEASE_GO_VERSION` in [.github/workflows/release-builds.yml](../.github/workflows/release-builds.yml).

To bump: update both in one PR, state the reason (new stable Go), and let
the next rc validate the change. Record the toolchain in the release notes.

## Release procedure

### 1. Release-content PR

- Promote `CHANGELOG.md` `[Unreleased]` into a dated `## [vX.Y.Z]` section
  (Keep-a-Changelog subsections, user-facing first); restore an empty
  `[Unreleased]`; update the compare links. The recorded date must equal
  the final tag's calendar date — if the tag slips, fix the date in a new
  commit and cut a new rc.
- Update anything version-specific (e.g. the `h3 --version` example in
  `cmd/h3/README.md`).
- Draft the release notes as reviewable Markdown files (outline in
  [versioning.md](versioning.md)), kept outside the repository until
  publication. Cite counts via tag-pinned links to the generated
  inventories instead of hard-coding numbers where exactness is not
  essential.

### 2. Local battery (on the exact release commit)

```sh
git status --porcelain && git rev-parse HEAD origin/master   # clean, equal
make fmt lint test check-unsafe check-layout check-docs check-ubercompare check-benchdocs
make -C testref h3-source && make check-api check-test-inventory check-cli-inventory
make test-c2go && make test-upstream-fixtures
go test -race ./...          # mandatory before tagging (ci-policy.md)
make vulncheck               # pinned govulncheck, root module
```

### 3. Pre-tag Nightly

`gh workflow run nightly.yml --ref master` → every job green (race, parity
+ fixtures, fuzz smoke, uberdiff/uberbench, CLI differential, CLI
cross-build, vulncheck, secret-scan).

### 4. Release candidate

```sh
git config user.name && git config user.email    # deliberate identity
git tag -a vX.Y.Z-rc.N -m "h3-go vX.Y.Z-rc.N — H3 Core vA.B.C (release candidate)"
git push origin vX.Y.Z-rc.N
```

The tag push triggers **two independent workflows — watch both**: Nightly
(full battery; the secret-scan run URL on this tag is the durable scan
evidence) and release-builds (`build` → `verify-reproducible` + `smoke`).

### 5. rc artifact verification

```sh
VERIFY=$(mktemp -d)                                    # always outside the repo
gh run download <release-builds-run-id> -n h3-vX.Y.Z-rc.N -D "$VERIFY/ci"
(cd "$VERIFY/ci" && shasum -a 256 -c SHA256SUMS)
```

- Inspect one archive's contents (binary + LICENSE + NOTICE + README.md;
  uid/gid 0; mtimes = tagged commit time).
- **Local archive-level reproduction (macOS or Linux — same command):**

  ```sh
  git checkout --detach vX.Y.Z-rc.N
  GOTOOLCHAIN=go<pinned> make release-dist VERSION=vX.Y.Z-rc.N OUT=$(mktemp -d)
  # compare that OUT/SHA256SUMS with the CI SHA256SUMS: must be identical
  git checkout master && test -z "$(git status --porcelain)"
  ```

- Native smoke: `h3 --version` → `h3 A.B.C (vX.Y.Z-rc.N)`, plus one golden
  command; linux/arm64 via `docker run --platform linux/arm64 alpine:3`.
- Module-zip preview (rc hash validates the rc only — rc and final hashes
  legitimately differ because the version string is part of every path in
  the module zip):

  ```sh
  GOPRIVATE=github.com/dimchansky/h3-go GOPROXY=direct \
    go mod download -json github.com/dimchansky/h3-go@vX.Y.Z-rc.N
  ```

- Audit the draft release-notes files against the current tree (claim
  audit #1).

Any fix — including to this document — goes through a PR and produces
**rc.N+1**. rc tags are kept forever and never get Release objects.

### 6. Final tag + private draft/publish (immutability enabled)

Preconditions: `git rev-parse master vX.Y.Z-rc.N^{}` identical;
`gh api repos/dimchansky/h3-go/immutable-releases` reports `enabled: true`.

```sh
git tag -a vX.Y.Z -m "h3-go vX.Y.Z — H3 Core vA.B.C"
git push origin vX.Y.Z
```

Watch both workflows; fast re-verify (checksums, `--version`, one smoke);
record the **final** module hashes pre-publication:

```sh
GOPRIVATE=github.com/dimchansky/h3-go GOPROXY=direct \
  go mod download -json github.com/dimchansky/h3-go@vX.Y.Z
# record BOTH "Sum" and "GoModSum"
```

Create the Release as a **draft** from the reviewed notes file, attach the
six archives + SHA256SUMS downloaded from the verified CI run, perform
claim audit #2 (last edit before immutability), then publish:

```sh
gh release create vX.Y.Z --draft --title "h3-go vX.Y.Z — H3 Core vA.B.C" \
  --notes-file notes-vX.Y.Z.md "$VERIFY/ci"/h3-vX.Y.Z-*.tar.gz \
  "$VERIFY/ci"/h3-vX.Y.Z-*.zip "$VERIFY/ci"/SHA256SUMS
gh release edit vX.Y.Z --draft=false --latest
```

Immediately verify while authenticated: the Release reports **Immutable**;
`gh release verify vX.Y.Z --repo dimchansky/h3-go` passes;
`gh release verify-asset` passes for **every** uploaded asset from a fresh
`gh release download` into a new temp dir; SHA256SUMS re-verifies.

### 7. Post-publication verification (public repository)

Proxy-only — never `GOPROXY=…,direct` in an acceptance check (a direct
fallback would silently prove nothing about proxy publication):

```sh
tmp=$(mktemp -d); cd "$tmp"; unset GOPRIVATE GONOPROXY GONOSUMDB
export GOMODCACHE="$tmp/m" GOCACHE="$tmp/c" GOBIN="$tmp/bin"
export GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org
# bounded poll (20 × 30 s), then:
go mod download -json github.com/dimchansky/h3-go@vX.Y.Z
#   "Sum" and "GoModSum" MUST equal the recorded pre-publication values
go list -m github.com/dimchansky/h3-go@latest        # MUST print vX.Y.Z
curl -fsS "https://sum.golang.org/lookup/github.com/dimchansky/h3-go@vX.Y.Z"
#   the "… vX.Y.Z h1:…" line == recorded Sum; "… vX.Y.Z/go.mod h1:…" == GoModSum
mkdir consumer && cd consumer && go mod init example.com/probe
go get github.com/dimchansky/h3-go@vX.Y.Z && go mod verify
go install github.com/dimchansky/h3-go/cmd/h3@vX.Y.Z && "$GOBIN/h3" --version
cd "$REPO" && test -z "$(git status --porcelain)"    # probes never touch the repo
```

Then: pkg.go.dev renders the version; unauthenticated curl checks (Release
page, every asset, SHA256SUMS re-check, badges); a hash mismatch anywhere
is a hard incident — stop and fix forward with a new version.

### 8. Cleanup gate

`git status --porcelain` empty and `git clean -ndX` inspected: no leftover
`dist/`, staging, notes drafts, or verification downloads anywhere in the
tree — temporary artifacts are removed, not gitignored.

## Smoke matrix (proven, not assumed)

Runner labels were architecture-proven (PR #4 probe, run 29350016630) and
every smoke job permanently asserts `uname -m` before executing anything:

| Runner | Asserted arch | Executes |
|---|---|---|
| ubuntu-latest | x86_64 | linux-amd64 (tar.gz) |
| windows-latest | x86_64 | windows-amd64 (zip) |
| macos-latest | arm64 | darwin-arm64 (tar.gz) |

Runtime-tested per release: linux/amd64, windows/amd64, darwin/arm64 (CI)
plus darwin/arm64 natively and linux/arm64 emulated locally. Build-only:
darwin/amd64, windows/arm64. State this matrix in the release notes.

## Branch and tag rulesets

Created immediately after the repository became public (GitHub Free offers
rulesets on public repositories). The canonical payloads are committed —
apply and verify with:

```sh
gh api -X POST repos/dimchansky/h3-go/rulesets --input .github/rulesets/protect-master.json
gh api -X POST repos/dimchansky/h3-go/rulesets --input .github/rulesets/protect-v-tags.json
gh api repos/dimchansky/h3-go/rulesets   # record ids, rules, bypass_actors, source binding
```

- [`protect-master`](../.github/rulesets/protect-master.json) (branch):
  require PR (0 approvals — the owner cannot approve their own PR;
  conversation resolution required; squash/rebase only), required status
  check **`CI / required`** bound to the GitHub Actions integration
  (`integration_id: 15368`), strict up-to-date policy, linear history, no
  force pushes, no deletion. Bypass: Repository admin (emergency-only, see
  Invariants).
- [`protect-v-tags`](../.github/rulesets/protect-v-tags.json) (tag,
  `refs/tags/v*`): creation, update, and deletion restricted to bypass
  actors (deliberate release tagging only). Published immutable releases
  additionally lock their tags regardless of rulesets.

After any ruleset change: inspect `gh api repos/dimchansky/h3-go/rulesets`
(+ per-ruleset detail) and re-run the negative tests from an isolated
worktree with the bypass actor temporarily removed (rejected push to
master, rejected `v*` test tag; verified cleanup; bypass restored and
API-verified).

## Maintainer-account security preflight (recurring)

Before every release (and any visibility/permission change): 2FA enabled
with at least two independent recovery methods (passkey or hardware key +
TOTP), recovery codes stored and present; review SSH keys, PATs, OAuth
apps, GitHub Apps, active sessions, deploy keys — revoke anything unused;
review `gh auth status` on the release machine. Record only that the audit
was completed, never credentials or recovery details.

## Failure handling

| Situation | Response |
|---|---|
| Any rc-stage failure | Fix via PR → new rc. Never reuse an rc tag. |
| Defect found after the final tag but before Release publication | Prefer fix-forward (new patch version). Tag deletion is not an option once anything may have consumed the tag. |
| Bad asset discovered after publication | Corrective release (assets are immutable); add `retract` for the bad version if consumers must be steered away. |
| Code defect in a published version | Fix-forward vX.Y.Z+1 with `retract vX.Y.Z // reason` when warranted. |
| Proxy/sumdb hash ≠ recorded pre-publication hash | Hard incident: stop, investigate (implies content divergence), resolve with a new version. |
| Vulnerability report | Private security advisory (see SECURITY.md), coordinated fix-forward release. |

## History

- **v0.3.0** (first public release): procedure executed as written;
  benchmark A/B evidence, runner probe, and scan records live in the
  maintainer's release-evidence archive. Known deferred item: the
  `FuzzUpstreamPolygonOperations` fuzz-rotation exception (issue #3).
