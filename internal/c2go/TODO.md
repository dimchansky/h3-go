Scope
- Track C→Go function conversions performed in `internal/c2go`.
- Record dependencies and testing status for each function.

Checklist
- [x] Scaffold package and docs (README.md)
- [x] Add first dependency-free function: `mathExtensions.c::_ipow`
- [x] Port `latLng.c::_posAngleRads` (angle normalization) — parity test via `_posAngleRadsC`
- [x] Port `latLng.c::constrainLng` and `constrainLat`
- [x] Port `latLng.c::degsToRads` and `radsToDegs`
- [x] Port `latLng.c::geoAlmostEqualThreshold` and `geoAlmostEqual`
- [ ] Port additional small helpers from `mathExtensions.c` and `latLng.c`
- [ ] Identify next targets with minimal dependencies; add TODO chains in code
- [ ] Mirror select C tests (apps/testapps, fuzzers) where practical

Functions
- mathExtensions.c
  - [x] `_ipow(int64_t base, int64_t exp)` — implemented; parity-tested via `_ipowC` wrapper
- latLng.c
  - [x] `_posAngleRads(double rads)` — DONE; parity via `_posAngleRadsC`
  - [x] `constrainLng(double lng)` — DONE
  - [x] `constrainLat(double lat)` — DONE
  - [x] `H3_EXPORT(degsToRads)(double degrees)` — DONE
  - [x] `H3_EXPORT(radsToDegs)(double radians)` — DONE
  - [x] `geoAlmostEqualThreshold(const LatLng*, const LatLng*, double)` — DONE
  - [x] `geoAlmostEqual(const LatLng*, const LatLng*)` — DONE

Conventions
- One function per Go file: `<cfile>__<function>.go`
- cgo interop per C module: `<cfile>_cgo.go` (build tag `cgo && c2go`), includes `"<cfile>.c"` by name only.
- C wrappers named distinctly (e.g., `_ipow_c_wrapper`), Go helpers mirror with a `C` suffix (e.g., `_ipowC`).
- Parity tests named `<cfile>__<function>_parity_test.go` with `//go:build c2go`. Reserve `<cfile>__<function>_test.go` for plain-Go tests later.
- Do not hardcode H3 version in code; include directories provided via `CGO_CPPFLAGS` in `make test-c2go` with `H3VER`.

Issue encountered: linking large C modules (resolved)
- Symptom: Including `latLng.c` directly from a single cgo file pulled in references to many other H3 C symbols (e.g., `cellToBoundary`, `directedEdgeToBoundary`), causing undefined symbol errors or duplicate static symbol definitions when multiple C files were included together.
- Constraints: We must not copy or reimplement C logic in the cgo layer.
- Solution applied:
  - Compile original H3 C modules as separate translation units via build-tagged C shim files in this package (e.g., `h3lib_latLng_c2go.c` contains `#include "latLng.c"`, `h3lib_h3Index_c2go.c` contains `#include "h3Index.c"`, etc.). This avoids static symbol redefinition (e.g., duplicate `DIRECTIONS`) by keeping each module isolated.
  - Keep `mathExtensions.c` referenced only once (through `mathExtensions_cgo.go`) to avoid duplicate `_ipow`.
  - Pass include directories via `CGO_CPPFLAGS` in `make test-c2go` using `H3VER` (no version in code).
  - Add `CGO_CFLAGS=-ffunction-sections -fdata-sections` and `CGO_LDFLAGS=-Wl,-dead_strip` so the linker drops unused functions referenced by the compiled C objects.
  - Tests call tiny Go wrappers around the original C symbols; cgo usage remains only in non-test files with the `c2go` tag. No C code is copied into Go files.

Guideline going forward when a C function links in extra deps
- First, add or reuse a per-module C shim file `h3lib_<module>_c2go.c` that `#include`s the needed original C file(s) by name.
- Ensure each original C file is compiled in its own shim to avoid static symbol redefinitions across modules.
- If link still fails due to further deps, add shims for those modules too (one TU per module).
- Keep the Go interop file (`<cfile>_cgo.go`) limited to declaring wrappers that forward to C; never copy C logic.


Note on C bool interop (cgo)
- Preferred: include `<stdbool.h>` and compare C.bool return values directly to `0` in Go (`C.fn(...) != 0`). This is valid because `_Bool` is an integer type in C99.
- Toolchain caveat: some cgo toolchains reject direct comparison `C.bool != 0`. In those cases, use a tiny inline C helper to normalize to `int` (e.g., `static int h3_bool_to_int(_Bool b) { return b ? 1 : 0; }`) and compare its result to `0` from Go.
