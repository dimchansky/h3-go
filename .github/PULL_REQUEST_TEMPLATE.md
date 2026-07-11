<!-- Thanks! Please see CONTRIBUTING.md for the ground rules. -->

## What & why


## Checklist

- [ ] `make fmt lint test check-unsafe` pass locally
- [ ] Touched ported code or the harness → `make test-c2go` passes
- [ ] Changed the public API → golden surface regenerated
      (`UPDATE_API_SURFACE=1 go test -run TestAPISurface .`) and doc comments
      carry `H3 C API:` lines where applicable
- [ ] Intentional divergence from H3 C → documented in `docs/DEVIATIONS.md`
- [ ] Collection-returning API → allocation assertion added/updated
