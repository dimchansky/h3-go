Scope
- Track C→Go function conversions performed in `internal/c2go`.
- Record dependencies and testing status for each function.

Design rule (performance, Go‑friendly APIs)
- Use the dst‑buffer slice pattern for any function producing collections.
  - Prefer signatures like `func fn(dst []T, args ...) ([]T, H3Error)`.
  - Reuse capacity in `dst` when sufficient; allocate only when `cap(dst)` is too small.
  - Keep element ordering identical to the C reference.

Signature mirroring (C pointers → Go pointers)
- When the C function takes struct pointers (e.g., `const LatLng*`, `const BBox*`), mirror that in Go with pointer parameters (`*LatLng`, `*BBox`, `*Vec2d`, `*Vec3d`, etc.).
  - Rationale: preserves in/out semantics, avoids accidental copies, and keeps the port traceable to the C API.
  - Const semantics: if the C pointer is `const`, keep the Go pointer parameter but do not mutate the pointed value.

Iteration Workflow (standard operating procedure)
- Select target: pick a small, self-contained C function (or a tight cluster) from the planned list.
- Prepare interop: if the C function is not exported or uses structs, add minimal cgo helpers to call it directly (use C structs like C.LatLng/C.BBox, etc.); avoid splitting scalars.
- Transpile: implement a faithful Go port in `internal/c2go/<cfile>__<function>.go`, mirroring names/signatures (unexported) and behavior (NO build tag for pure Go).
- Document: follow established comment style with function description, technical details, and `// Ported from H3 C: <file>::<function>` attribution line.
- Parity test: add `internal/c2go/<cfile>__<function>_parity_test.go` with `//go:build cgo` that compares Go vs C output; use tight, justified tolerances for floats and safe bool handling.
- Sanity run: `make test-c2go` and ensure all parity tests pass.
- Format code: run `make fix-fmt` prior to committing.
- Update tracker: update this TODO to mark the function DONE and list the next planned items (do this BEFORE each commit).
- Commit: commit the minimal, focused changes with a message stating the ported function(s), parity, and TODO update.

Source Reference
- Local H3 C source code available at: `testref/h3-4.3.0/src/h3lib/lib/*.c`
- When porting C code, reference the local sources directly for accurate implementation
- Port as closely as possible to the C implementation, preserving algorithm details and edge cases

Build Tag Rules
- **Pure Go implementations** (`<cfile>__<function>.go`): NO build tags
- **CGO interop files** (`*_cgo.go`): Must have `//go:build cgo` as first line
- **Parity test files** (`*_parity_test.go`): Must have `//go:build cgo` as first line
- This ensures pure Go code is always available while CGO dependencies are isolated

Comment Style Requirements
- **Function description**: Brief summary of what the function does
- **Technical details**: Include important implementation notes, algorithm explanations, or edge case handling
- **Source attribution**: Must include `// Ported from H3 C: <file>::<function>` line
- **Example format**:
  ```go
  // _functionName calculates something important.
  // Additional technical details about the implementation.
  // Ported from H3 C: fileName.c::_functionName
  func _functionName(params) returnType {
  ```
- **Constants**: Should be declared in `<cfile>_constants.go` files, not inline

Interop helper note
- When performing Go→C struct conversions in `*_cgo.go`, prefer creating small helper functions to convert common types (e.g., `toCGeoLoop`, `toCGeoPolygon`, `toCBBox`). This avoids duplication when multiple functions need the same transformation and makes memory management (malloc/free) clearer and safer.

## Current Status

**Run the following command to get current implementation status:**
```bash
./scripts/update-h3-status.sh
```

This will generate/update:
- `internal/c2go/h3_functions.md` - Public API implementation status (checklist)
- `internal/c2go/ported_functions_report.md` - Detailed report of all ported internal functions

**Quick status check:**
```bash
./scripts/list-ported-functions.sh  # Live terminal output of ported functions
```

## Completed Milestones
- [x] Scaffold package and basic infrastructure
- [x] Core utility functions (math, angles, coordinates)  
- [x] Geographic and coordinate system functions
- [x] H3 index manipulation and validation
- [x] Automated status tracking and reporting system

Non-negotiable constraints (repo policy)
- Do not copy or vendor H3 `.c` or `.h` files into this repo. Only include original files via build-tagged shims (e.g., `h3lib_*.c`) and let `CGO_CPPFLAGS` provide include paths.
- Keep all c2go work scoped to `internal/c2go` only; do not introduce dependencies on other Go packages.

## Planning Next Targets

Use the automated reports to identify the next functions to implement:

1. **Check public API gaps**: Review `internal/c2go/h3_functions.md` for unchecked functions
2. **Review available building blocks**: Check `internal/c2go/ported_functions_report.md` for internal functions that can support new public APIs
3. **Select minimal dependency targets**: Choose functions that depend only on already-ported internal functions

**Strategy**: Focus on completing public API functions that have the most internal dependencies already satisfied.

Execution plan per function
- Extend `<cfile>_cgo.go` with direct calls using C structs (C.BBox/C.LatLng/GeoLoop); avoid scalar params.
- Implement faithful Go ports in `<cfile>__<function>.go` with same unexported names.
- Add `<cfile>__<function>_parity_test.go` with realistic tolerances; compare bools directly.
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

Execution plan per function
- Extend `<cfile>_cgo.go` with direct calls to the original C functions using C structs where applicable (no scalar explosion; use C.LatLng / C.Vec2d).
- Implement the Go port in `<cfile>__<function>.go` (faithful translation; same unexported name where possible).
- Add parity test `<cfile>__<function>_parity_test.go` under `//go:build c2go` with tight but realistic tolerances.
- If C returns `bool`, prefer `<stdbool.h>` and compare `!= 0`. If toolchain warns, route through a tiny C helper returning `int`.

## Update Status After Implementation

After completing function implementations, always run:
```bash
./scripts/update-h3-status.sh
```

This ensures the automated tracking stays current and provides accurate progress reporting.
