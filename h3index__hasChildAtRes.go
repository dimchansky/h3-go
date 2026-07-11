package h3

// _hasChildAtRes determines whether childRes is a valid child resolution for h.
// Each resolution is considered a valid child resolution of itself.
// Ported from H3 C: h3Index.c::_hasChildAtRes.
func _hasChildAtRes(h H3Index, childRes int32) bool {
	parentRes := getResolution(h)
	if childRes < parentRes || childRes > MAX_H3_RES {
		return false
	}
	return true
}
