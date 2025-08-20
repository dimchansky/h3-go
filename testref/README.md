# H3 Reference Test Oracle

This directory will contain the external C oracle CLI for testing our pure-Go H3 implementation against the reference H3 C v4.3.0 implementation.

## Overview

The Go library is intentionally pure Go (no cgo, no external dependencies). To ensure correctness, we build a separate C binary that wraps H3 C v4.3.0 and provides a simple CLI interface for Go tests to invoke via `exec.Command`.

## Architecture

- **Go library**: Pure Go implementation in `../` (this package)
- **C oracle**: Separate C binary built from H3 v4.3.0 source (this directory)
- **Test harness**: Go tests that invoke the C oracle and compare results

## Building the Oracle

```bash
make ref  # From repository root
```

This will:
1. Download H3 v4.3.0 source 
2. Build the `testref/h3ref` CLI binary
3. The binary provides a JSON stdin/stdout interface for test queries

## Protocol

The `h3ref` binary expects JSON commands on stdin and returns JSON responses on stdout:

```json
# Input
{"function": "latLngToCell", "args": {"lat": 37.775, "lng": -122.418, "res": 9}}

# Output  
{"result": "0x8928308280fffff", "error": null}
```

## Error Code Mapping

The C oracle returns numeric error codes that must be mapped to Go errors:

```go
// C error code -> Go error sentinel mapping (parity with H3 C v4.3.0).
var cErrToGo = map[uint32]error{
    0:  nil,                    // Success
    1:  ErrFailed,
    2:  ErrDomain,
    3:  ErrLatLngDomain,
    4:  ErrResolutionDomain,
    5:  ErrCellInvalid,
    6:  ErrDirectedEdgeInvalid,
    7:  ErrUndirectedEdgeInvalid,
    8:  ErrVertexInvalid,
    9:  ErrPentagon,
    10: ErrDuplicateInput,
    11: ErrNotNeighbors,
    12: ErrResolutionMismatch, // corrected name
    13: ErrMemoryAlloc,
    14: ErrMemoryBounds,
    15: ErrOptionInvalid,
}
```

Unknown error codes map to `ErrFailed` with logging.

## Status

- [ ] Download and build H3 C v4.3.0
- [ ] Implement `h3ref` CLI with JSON protocol  
- [ ] Add `make ref` target to root Makefile
- [ ] Wire up Go test harness to use oracle
- [ ] Add golden test datasets for stable operations

## Design Goals

1. **No linking**: The Go library never links to C code - oracle is a separate process
2. **Hermetic builds**: Oracle build is reproducible and isolated
3. **Behavioral parity**: All operations match H3 C v4.3.0 exactly
4. **Performance isolation**: Oracle is only used in tests, not production code