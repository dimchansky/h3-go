//go:build cgo && c2go

package h3

// Bridge to the public reverseDirectedEdge added to h3api.h in H3 4.5.0
// (absent from the 4.4.0 tree).

/*
#include "h3api.h"
*/
import "C"

// reverseDirectedEdgeC calls the original C implementation.
func reverseDirectedEdgeC(edge h3Index) (h3Index, h3Error) {
	var out C.H3Index
	err := h3Error(C.reverseDirectedEdge(C.H3Index(edge), &out))
	return h3Index(out), err
}
