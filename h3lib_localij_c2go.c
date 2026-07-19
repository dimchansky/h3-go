//go:build cgo && c2go

#include "localij.c"

// Test-only wrapper
// exposing the file-static gridPathCellsInterpolate to the parity
// harness, in the same translation unit as the static it calls.
H3Error h3goTest_gridPathCellsInterpolate(H3Index start, H3Index end,
                                          int64_t distance, H3Index *out,
                                          int64_t outOffset,
                                          int64_t outStep) {
    return gridPathCellsInterpolate(start, end, distance, out, outOffset,
                                    outStep);
}
