// Package badpkg is a doclinkcheck test fixture.
//
// [NoSuchSymbol] must be reported as unresolved, while [Good] and
// [Good.Method] resolve. Non-candidates like [0, 1) and [lng, lat] are
// ignored.
package badpkg

// Good is referenced from the package comment.
type Good struct{}

// Method is referenced as [Good.Method].
func (Good) Method() {}
