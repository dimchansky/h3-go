package h3

// _iterateBaseCellIndexesAtRes iterates through all H3 indexes at the given resolution
// within a specific base cell. The callback function is called for each valid H3 index.
// This is useful for iterating through a subset of the H3 grid constrained to one base cell.
// Ported from H3 C: utility.c::iterateBaseCellIndexesAtRes.
func _iterateBaseCellIndexesAtRes(res int32, callback func(H3Index), baseCell int32) {
	iter := iterInitBaseCellNum(baseCell, res)
	for iter.H != 0 {
		callback(iter.H)
		iterStepChild(&iter)
	}
}
