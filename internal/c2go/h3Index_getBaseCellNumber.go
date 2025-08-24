package c2go

// getBaseCellNumber returns the base cell number of an H3 index.
// Note: Technically works on H3 edges, but will return base cell of the origin cell.
// Ported from H3 C: h3Index.c::getBaseCellNumber
func getBaseCellNumber(h H3Index) int32 {
	return getBaseCell(h)
}
