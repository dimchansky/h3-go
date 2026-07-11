package h3

// _iterateAllIndexesAtRes iterates through all H3 indexes at the given resolution.
// The callback function is called for each valid H3 index encountered.
// This covers all base cells (0 through numBaseCells-1) and their children.
// Ported from H3 C: utility.c::iterateAllIndexesAtRes.
func _iterateAllIndexesAtRes(res int32, callback func(h3Index)) {
	_iterateAllIndexesAtResPartial(res, callback, numBaseCells)
}
