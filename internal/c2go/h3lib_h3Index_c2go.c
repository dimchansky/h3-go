//go:build cgo && c2go

#include "h3Index.c"

// Expose static helpers for cgo parity
int has_child_at_res_c(H3Index h, int childRes) { return _hasChildAtRes(h, childRes) ? 1 : 0; }
