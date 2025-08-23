package c2go

import "fmt"

// h3ToString converts an H3 index to its lowercase hex string form.
// Ported behavior: always succeeds; returns 0 error code.
// Ported from H3 C: h3Index.c::h3ToString
func h3ToString(h H3Index) (string, uint32) { return fmt.Sprintf("%x", uint64(h)), 0 }
