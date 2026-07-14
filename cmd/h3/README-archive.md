# h3 — standalone CLI for Uber's H3 grid system (pure Go)

This archive contains a single self-contained `h3` executable (no runtime
dependencies, no C library needed), a copy of the license and notice files,
and this README. It is the upstream-compatible command-line interface of
[github.com/dimchansky/h3-go](https://github.com/dimchansky/h3-go), a pure-Go
port of [Uber's H3](https://github.com/uber/h3) — behaviorally equivalent to
H3 C v4.4.0.

> This file ships inside release archives, so every link below is absolute.
> The full CLI documentation lives at
> <https://github.com/dimchansky/h3-go/blob/master/cmd/h3/README.md>.

## Install

Extract and place the binary anywhere on your `PATH`:

```sh
tar -xzf h3-<version>-<os>-<arch>.tar.gz        # linux / macOS
# or unzip h3-<version>-windows-<arch>.zip      # Windows
install -m 0755 h3-<version>-<os>-<arch>/h3 ~/bin/h3
```

macOS Gatekeeper note: the binaries are not notarized; if macOS quarantines
the download, clear it with `xattr -d com.apple.quarantine h3`.

## Quick start

```sh
h3 --help                                        # list every command
h3 latLngToCell -r 9 --lat 37.7759 --lng -122.4180
# "8928308280fffff"
h3 cellToBoundary -c 8928308280fffff
```

All 63 commands of the H3 C v4.4.0 CLI are supported with matching flags,
output, and exit codes.

## Version

```sh
h3 --version
# h3 4.4.0 (v0.3.0)
```

The first number is the H3 C compatibility target; the value in parentheses
is the h3-go release this binary was built from (injected from the Git tag
at build time). The two are independent axes — see
<https://github.com/dimchansky/h3-go/blob/master/docs/versioning.md>.

## Verifying your download

Every release publishes a `SHA256SUMS` manifest next to the archives
(tokenless, works anywhere):

```sh
shasum -a 256 -c SHA256SUMS      # or: sha256sum -c SHA256SUMS
```

Releases are published as GitHub **immutable releases** (tag and assets are
locked at publication and carry a release attestation). With an
**authenticated** GitHub CLI session (`gh auth login` — these commands do
not work anonymously), you can verify cryptographically from anywhere:

```sh
gh release verify <version> --repo dimchansky/h3-go
gh release verify-asset <version> <downloaded-file> --repo dimchansky/h3-go
```

Archives are built reproducibly (pinned Go toolchain, `CGO_ENABLED=0`,
`-trimpath`, normalized archive metadata); the release notes state the exact
toolchain and a one-line rebuild command.

## Supported platforms

Archives are published for linux, macOS (darwin), and Windows on amd64 and
arm64. The release notes for each version state which platforms were
runtime-tested on real hardware for that release and which are cross-built
only — see
<https://github.com/dimchansky/h3-go/releases>.

## License

Apache License 2.0 — see the bundled `LICENSE` and `NOTICE` files, or
<https://github.com/dimchansky/h3-go/blob/master/LICENSE>.
