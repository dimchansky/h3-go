//go:build cgo && c2go

#include "vertex.c"

// Wrapper to access the static vertexRotations function
H3Error vertexRotations_wrapper(H3Index cell, int *out) {
    return vertexRotations(cell, out);
}

