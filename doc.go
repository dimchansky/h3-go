// Package h3 is a pure-Go implementation of Uber's H3 hexagonal hierarchical
// geospatial indexing system (reference version 4.3.0).
//
// The implementation layer is a function-by-function port of the H3 C library,
// validated against the original C objects by an opt-in cgo parity test suite
// (build tags cgo && c2go). The production library itself is safe Go only: no
// cgo and no unsafe.
package h3
