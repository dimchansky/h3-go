# Versioning and release policy

The authoritative policy for how this project is versioned and released.
Other documents (README, [CHANGELOG](../CHANGELOG.md),
[FUTURE_WORK](FUTURE_WORK.md)) summarize or link here rather than
restating it.

## Two independent version axes

The project has two separate versions, and they must not be conflated:

1. **The h3-go module version** — the semantic version of this Go module
   and its public API. This is what a Git tag names.
2. **The H3 Core compatibility target** — the exact upstream H3 C release
   this library is behaviorally equivalent to (currently v4.4.0; every
   difference must be listed in [DEVIATIONS.md](DEVIATIONS.md)).

An illustrative pair (not a commitment to the next release numbers):

```text
h3-go release:             v0.3.0
H3 Core behavioral target: v4.5.0
```

The Git tag is the h3-go module version and nothing else. The H3 Core
target is release *metadata*, recorded prominently (below) but never
encoded into the tag.

## Semantic versioning for h3-go

Tags are ordinary, canonical Go-module SemVer: `v0.x.y`, `v1.x.y`,
`v2.x.y`, …

- **While at v0**: minor releases may contain breaking API changes, always
  called out explicitly in the [CHANGELOG](../CHANGELOG.md); patch
  releases contain compatible fixes and improvements.
- **From v1.0.0 on**: normal SemVer — major for backward-incompatible
  public API or documented-behavior changes, minor for backward-compatible
  functionality, patch for compatible fixes and optimizations. The
  pre-v1.0.0 checklist is in
  [FUTURE_WORK.md](FUTURE_WORK.md#release-considerations-before-v100).
- The module path stays `github.com/dimchansky/h3-go`. It is **not**
  migrated to `/v4` (or any `/vN`) merely to mirror H3 Core's major
  version; a `/v2` path would appear only if this library itself ever
  ships a v2.0.0.

## One version per tag

Do not encode both axes into one tag:

- **No four-component tags** (`v0.4.4.0`) — not valid SemVer, not
  canonical Go module versions.
- **No prerelease or build-metadata second axis** (`v0.3.0-h3.4.5.0`,
  `v0.3.0+h3.4.5.0`) — prerelease suffixes have ordering and stability
  semantics of their own (a `-h3.…` "release" would sort *before* the
  release and be treated as unstable), and build metadata is ignored for
  precedence and not preserved as a distinct canonical module version.
- Prerelease suffixes such as `-rc.1` remain valid for genuine
  prereleases of an h3-go version — never for naming the H3 target.

## Tag immutability

While the repository is private and no external consumer or module proxy
can have fetched a version, rewriting or deleting the experimental tags
is permitted as an intentional, coordinated operation. From the moment
the repository is public — or a module version may have been consumed or
cached (e.g. by the Go module proxy) — published tags are immutable:
never moved, reused, or silently replaced. A bad release is followed by a
new version and, when appropriate, a Go module
[`retract` directive](https://go.dev/ref/mod#go-mod-file-retract) for the
broken one.

## Where the H3 Core target is recorded

Every release states the exact H3 Core target in:

- the **GitHub Release title** — preferred form:
  `h3-go vX.Y.Z — H3 Core vA.B.C` (e.g. `h3-go v0.3.0 — H3 Core v4.5.0`);
- the **first lines of the release notes** (outline below);
- the **README** H3-version badge and its Compatibility and versioning
  section;
- the exported **`VersionMajor` / `VersionMinor` / `VersionPatch`**
  constants, which report the H3 Core target (not the module version) —
  the CLI separately reports its own build version, stamped from the Git
  tag by the [release-builds workflow](../.github/workflows/release-builds.yml);
- the **CHANGELOG** entry for the release;
- the comparison, migration, benchmark, and CLI documents when they pin
  versions ([comparison-uber-h3-go.md](comparison-uber-h3-go.md#versions-compared),
  [migration-from-uber-h3-go.md](migration-from-uber-h3-go.md),
  [benchmarks](benchmarks/README.md),
  [cli-compatibility.md](cli-compatibility.md)).

The annotated Git tag message may stay short; the GitHub Release and the
CHANGELOG carry the full details.

## Choosing the h3-go increment

An H3 Core upgrade does not mechanically determine the h3-go version
number. Pick the increment from the upgrade's effect on *this module's*
API and documented behavior. Illustrative examples only — not a
commitment to the next release number:

| Event | h3-go version | H3 Core target |
|---|---:|---:|
| Adopt H3 Core 4.5.0 during pre-v1 development | `v0.3.0` | `v4.5.0` |
| Compatible h3-go bug fix | `v0.3.1` | `v4.5.0` |
| Breaking pre-v1 API change | `v0.4.0` | `v4.5.0` |
| Stabilize the public API | `v1.0.0` | exact target stated separately |
| H3 upgrade after v1 adding backward-compatible functionality | next minor release | new exact target |
| H3 upgrade after v1 with only compatible fixes/optimizations | next patch release | new exact target |
| H3 upgrade after v1 changing documented behavior incompatibly | next major release | new exact target |
| Breaking Go API change after v1 | next major release | exact target stated separately |

When an H3 patch release is adopted, the h3-go patch number does not need
to equal the H3 Core patch number; the exact target is recorded
separately as described above.

## Why upstream-aligned tags were rejected

The official binding tags releases like `v4.5.0` because its module path
is `github.com/uber/h3-go/v4` and its release line deliberately follows
the H3 major/minor line. That scheme was considered and rejected for this
independent implementation because it couples two different compatibility
lifecycles:

- a breaking h3-go API change would require a new Go major version (and
  import path) even if H3 Core stayed on the same line;
- a new H3 major version would force a new Go import path even if this
  library's API remained compatible;
- even the binding's patch number is not always identical to H3 Core's
  patch number, so the coupling is already imperfect where it is used.

Do not reopen this decision without new evidence.

## Release-note outline

Recommended structure for GitHub Release notes:

```md
# h3-go vX.Y.Z — H3 Core vA.B.C

Library version: vX.Y.Z
Behavioral compatibility target: H3 Core vA.B.C

## Upstream H3 changes adopted
## Go API changes
## Fixed
## Compatibility and migration
## Validation
## CLI artifacts
```

When an upstream version is adopted, link the corresponding upstream H3
release (and the reviewed sync record under [sync/](sync/4.3.0-to-4.4.0.md));
the notes must also describe the h3-go-specific API, implementation,
documentation, and validation changes — the release is of this library,
not of H3. *CLI artifacts* lists the reproducible archives and
`SHA256SUMS` built by the
[release-builds workflow](../.github/workflows/release-builds.yml).
