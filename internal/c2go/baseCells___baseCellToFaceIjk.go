
package c2go

// _baseCellToFaceIjk converts a base cell number to the FaceIJK address of that base cell.
// This is a simple lookup in the baseCellData table: h = baseCellData[baseCell].homeFijk
func _baseCellToFaceIjk(baseCell int, h *FaceIJK) {
	*h = baseCellData[baseCell].HomeFijk
}