# testref — upstream H3 C sources for parity testing

This directory fetches (never vendors) the pristine upstream
[H3 C](https://github.com/uber/h3) sources that the cgo parity suite and the
API-completeness gate compare against.

```sh
make -C testref h3-source              # download + extract the default version
make -C testref H3_VERSION=4.5.0 h3-source   # fetch another version (upstream syncs)
```

The downloaded trees (`h3-<version>/`) and tarballs are gitignored; only this
scaffolding is committed. The version consumed by the root Makefile is set by
`H3VER` there (`make test-c2go H3VER=...`); this Makefile's `H3_VERSION`
controls what gets downloaded.

How the sources are used:

- **Parity suite** (`make test-c2go`, build tags `cgo && c2go`): the root
  package's `h3lib_*_c2go.c` shims `#include` the original C files by name —
  include paths are injected via `CGO_CPPFLAGS`, so no H3 version is ever
  hardcoded in code — and `*_cgo.go` wrappers let 227 parity tests compare Go
  vs C behavior in-process.
- **API gates** (`make check-api`, `make api-inventory`):
  `tools/apiinventory` parses `h3api.h.in` to enforce that every C public
  function is ported and publicly represented.

`h3ref.c` builds a small standalone CLI (`make -C testref`) for manually
querying the reference implementation while debugging, e.g.
`./testref/h3ref latLngToCell 37.775 -122.418 9`. No Go test depends on it.
