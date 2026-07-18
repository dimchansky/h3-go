//go:build cgo && c2go

// Tree-shape-adaptive shim (docs/sync/4.4.0-to-4.5.0.md §15.1): the H3
// 4.4.0 tree implements coordijk in lib/coordijk.c; 4.5.0 moved the same
// functions (verified body-identical) into include/coordijk.h as static
// inline definitions. When the .c file exists we compile it here; when it
// does not, this becomes an empty translation unit and every cgo bridge
// that includes coordijk.h gets the inline definitions in its own TU.
#if __has_include("coordijk.c")
#include "coordijk.c"
#endif
