package c2go

// _hasChildAtRes determines whether childRes is a valid child resolution for h.
// Each resolution is considered a valid child resolution of itself.
func _hasChildAtRes(h H3Index, childRes int) bool {
	parentRes := getResolution(h)
	if childRes < parentRes || childRes > MAX_H3_RES {
		return false
	}
	return true
}
