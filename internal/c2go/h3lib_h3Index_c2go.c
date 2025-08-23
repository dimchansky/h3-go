//go:build cgo && c2go

#include "h3Index.c"

// Expose static helpers for cgo parity
int has_child_at_res_c(H3Index h, int childRes) { return _hasChildAtRes(h, childRes) ? 1 : 0; }
int first_one_index_c(H3Index h) { return _firstOneIndex(h); }
int has_good_top_bits_c(H3Index h) { return _hasGoodTopBits(h) ? 1 : 0; }
