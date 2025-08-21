Scope
- Track C→Go function conversions performed in `internal/c2go`.
- Record dependencies and testing status for each function.

Checklist
- [x] Scaffold package and docs (README.md)
- [x] Add first dependency-free function: `mathExtensions.c::_ipow`
- [ ] Port additional small helpers from `mathExtensions.c` and `latLng.c`
- [ ] Identify next targets with minimal dependencies; add TODO chains in code
- [ ] Mirror select C tests (apps/testapps, fuzzers) where practical

Functions
- mathExtensions.c
  - [x] `_ipow(int64_t base, int64_t exp)` — implemented; parity-tested via `_ipowC` wrapper

Conventions
- One function per Go file: `<cfile>__<function>.go`
- cgo interop per C module: `<cfile>_cgo.go` (build tag `cgo && c2go`), includes `"<cfile>.c"` by name only.
- C wrappers named distinctly (e.g., `_ipow_c_wrapper`), Go helpers mirror with a `C` suffix (e.g., `_ipowC`).
- Parity tests named `<cfile>__<function>_parity_test.go` with `//go:build c2go`. Reserve `<cfile>__<function>_test.go` for plain-Go tests later.
- Do not hardcode H3 version in code; include directories provided via `CGO_CPPFLAGS` in `make test-c2go` with `H3VER`.
