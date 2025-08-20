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
1. Download H3 v4.3.0 source to `testref/h3-4.3.0/`
2. Build the `testref/h3ref` CLI binary
3. The binary provides a simple command-line interface for test queries

## Direct Build

You can also build directly in the testref directory:

```bash
cd testref
make        # Download and build oracle
make test   # Run validation tests
make version    # Show current H3 version
make clean-all  # Remove all downloaded source and binaries
```

## Upgrading H3 Version

To upgrade to a newer H3 version:

1. Edit `H3_VERSION` in `testref/Makefile`
2. Run `make clean-all && make` to rebuild with new version
3. Run `make test` to verify compatibility

The oracle source code is version-agnostic and requires no changes.

## Protocol

The `h3ref` binary provides a simple command-line interface:

```bash
# Test pentagon detection
./h3ref pentagon 4        # Returns: 1 (is pentagon)
./h3ref pentagon 0        # Returns: 0 (is hexagon)

# Convert FaceIJK to H3 index  
./h3ref faceijk 0 0 0 1 0  # Returns: 0x8025fffffffffff

# Convert LatLng to H3 index
./h3ref latlng 37.775 -122.418 9  # Returns: 0x8928308280fffff 0

# Rotate H3 indices
./h3ref rotate60cw 0x8021fffffffffff   # Returns: 0x8021fffffffffff
./h3ref rotate60ccw 0x8021fffffffffff  # Returns: 0x8021fffffffffff
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

- [x] Download and build H3 C v4.3.0
- [x] Implement `h3ref` CLI with command-line protocol
- [x] Add `make ref` target to root Makefile
- [ ] Wire up Go test harness to use oracle
- [ ] Add golden test datasets for stable operations
- [x] Validate pentagon handling and rotation functions

## Design Goals

1. **No linking**: The Go library never links to C code - oracle is a separate process
2. **Hermetic builds**: Oracle build is reproducible and isolated
3. **Behavioral parity**: All operations match H3 C v4.3.0 exactly
4. **Performance isolation**: Oracle is only used in tests, not production code