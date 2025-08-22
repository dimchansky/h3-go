Purpose
- Provide a dedicated, low-friction workspace for near line-by-line C→Go function ports from H3 C.
- Make verification easy: every function here has a cgo-backed parity test comparing results to the original C implementation.

Conventions
- Package: `internal/c2go` (internal-only; not part of public API).
- One function per Go file: `<cfile>__<function>.go` (unexported, keep original name when possible).
- One cgo interop file per C module: `<cfile>_cgo.go` — includes the original C file by name and exposes small C→Go wrappers for parity tests.
- Do not hardcode H3 versions in code. Include directories are provided via `CGO_CPPFLAGS` at test time (see Makefile target `test-c2go`).
- If a function depends on others not yet ported, add a `TODO:` with the exact C symbol names, then port dependencies recursively.

Testing
- Parity tests compare the Go port vs the C original and are guarded by build tag `c2go`.
- Parity test naming: `<cfile>__<function>_parity_test.go` (reserved `<cfile>__<function>_test.go` for plain-Go tests later).
- cgo is only used in non-test files (`*_cgo.go` with `//go:build cgo && c2go`). Test files import Go wrappers only.
- Keep test inputs within ranges that avoid C UB/overflows unless explicitly testing boundaries.

How to run
- Build C reference once: `make ref`
- Run parity tests: `make test-c2go` (defaults to `H3VER=4.3.0`)
- Override version: `make H3VER=4.4.0 test-c2go`
 - Format code before committing: `make fix-fmt`

Porting algorithm
- Pick a function in `testref/h3-<ver>/src/h3lib/lib/<cfile>.c` with minimal deps.
- Create `internal/c2go/<cfile>__<function>.go` with a faithful Go translation; keep the same unexported name.
- Add/extend `internal/c2go/<cfile>_cgo.go`:
  - `//go:build cgo && c2go`
  - `#include "<cfile>.c"` by name only (no version in code).
  - Add a small C wrapper (e.g., `_ipow_c_wrapper`) and a Go helper mirroring the name (e.g., `_ipowC`).
- Add `internal/c2go/<cfile>__<function>_parity_test.go` (build tag `c2go`) that compares Go vs C wrappers.
- If the function depends on other C helpers, add TODOs with exact symbol names and port those recursively.

Notes
- This workspace is for correctness parity and traceability. Optimized implementations used by the public API can live elsewhere.
- Library package remains pure Go. cgo exists only in `internal/c2go/*_cgo.go` and only builds with the `c2go` tag.
