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

**Verify first, then extract** (see "Verifying your download" below for the
one-command single-asset check), and place the binary anywhere on your
`PATH`:

```sh
tar -xzf h3-<version>-<os>-<arch>.tar.gz        # linux / macOS
# or unzip h3-<version>-windows-<arch>.zip      # Windows
install -m 0755 h3-<version>-<os>-<arch>/h3 ~/bin/h3
```

### macOS Gatekeeper

The prebuilt macOS binaries are **ad-hoc linker-signed, but are not
Developer ID-signed or notarized**, so Gatekeeper can reject them when the
archive was downloaded through a browser. The safe order is: verify the
downloaded archive, extract it, remove the quarantine attribute **only from
the exact verified binary**, then run it:

```sh
archive=h3-v0.3.0-darwin-arm64.tar.gz            # darwin-amd64 on Intel Macs
actual=$(shasum -a 256 "$archive") &&
  grep -Fxq -- "$actual" SHA256SUMS &&
  echo "OK: checksum verified" &&
  tar -xzf "$archive" &&
  xattr -d com.apple.quarantine h3-v0.3.0-darwin-arm64/h3 &&
  h3-v0.3.0-darwin-arm64/h3 --version
```

The whole sequence is one `&&` chain: extraction and quarantine removal run
**only after** `OK: checksum verified` — a missing or corrupted archive
stops the chain with a non-zero status before anything is extracted.

Alternatives and cautions:

- If you already have Go, `go install github.com/dimchansky/h3-go/cmd/h3@v0.3.0`
  is the simplest path — a locally built binary never inherits
  browser-download quarantine.
- Apple's UI path after a blocked first launch: System Settings →
  Privacy & Security → **Open Anyway**.
- Remove quarantine only **after** checksum or release-attestation
  verification, never before.
- Never disable Gatekeeper globally, and never run broad commands such as
  `xattr -dr` over `~/Downloads` or any directory — target only the one
  verified binary.

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
(tokenless). It lists **all six** archives; verify just the one you
downloaded, from the directory containing the archive and `SHA256SUMS`.

macOS:

```sh
archive=h3-<version>-darwin-<arch>.tar.gz
actual=$(shasum -a 256 "$archive") &&
  grep -Fxq -- "$actual" SHA256SUMS &&
  echo "OK: checksum verified"
```

Linux:

```sh
archive=h3-<version>-linux-<arch>.tar.gz
actual=$(sha256sum "$archive") &&
  grep -Fxq -- "$actual" SHA256SUMS &&
  echo "OK: checksum verified"
```

Windows (PowerShell):

```powershell
$archive = "h3-<version>-windows-<arch>.zip"
$entries = @(Select-String -Path SHA256SUMS -SimpleMatch $archive)
if ($entries.Count -ne 1) { throw "expected exactly 1 manifest entry for $archive, found $($entries.Count)" }
$hash = (Get-FileHash $archive -Algorithm SHA256).Hash.ToLower()
if ($entries[0].Line -ne "$hash  $archive") { throw "checksum mismatch for $archive" }
"OK: checksum verified"
```

The Unix commands are fail-closed: `grep -Fxq` requires the complete
computed `<hash>  <name>` line to match a manifest line **exactly** (fixed
string, whole line — no regular expressions), so a match verifies the hash
and the exact file name together, and every failure — a missing archive
(the checksum tool's own error stops the `&&` chain), a corrupted archive,
or a name not in the manifest — ends with a non-zero exit status and
**without** the `OK: checksum verified` line, which is printed only after
successful verification. The PowerShell version fails explicitly on zero
or multiple manifest entries and on any hash mismatch. Every release run
executes these exact commands, plus corrupted-archive and missing-archive
negative cases, on real macOS, Linux, and Windows runners.
(`shasum -a 256 -c SHA256SUMS` on macOS or `sha256sum -c SHA256SUMS` on
Linux checks the whole manifest at once, but requires all six archives to
be present in the directory.)

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
