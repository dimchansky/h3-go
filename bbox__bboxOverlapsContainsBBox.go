package h3

// bboxOverlapsBBox reports whether two bboxes overlap (port of bbox.c)
// Ported from H3 C: bbox.c::bboxOverlapsBBox.
func bboxOverlapsBBox(a, b *bbox) bool {
	// Latitude overlap
	if a.North < b.South || a.South > b.North {
		return false
	}
	// Longitude overlap with normalization
	var aNorm, bNorm longitudeNormalization
	bboxNormalization(a, b, &aNorm, &bNorm)
	if normalizeLng(a.East, aNorm) < normalizeLng(b.West, bNorm) ||
		normalizeLng(a.West, aNorm) > normalizeLng(b.East, bNorm) {
		return false
	}
	return true
}

// bboxContainsBBox reports whether a contains b (port of bbox.c)
// Ported from H3 C: bbox.c::bboxContainsBBox.
func bboxContainsBBox(a, b *bbox) bool {
	if a.North < b.North || a.South > b.South {
		return false
	}
	var aNorm, bNorm longitudeNormalization
	bboxNormalization(a, b, &aNorm, &bNorm)
	return normalizeLng(a.West, aNorm) <= normalizeLng(b.West, bNorm) &&
		normalizeLng(a.East, aNorm) >= normalizeLng(b.East, bNorm)
}
