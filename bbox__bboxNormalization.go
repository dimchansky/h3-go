package h3

// bboxNormalization determines longitude normalization for two bboxes.
// Mirrors H3's bbox.c::bboxNormalization behavior.
// Ported from H3 C: bbox.c::bboxNormalization.
func bboxNormalization(a, b *bbox, aNorm, bNorm *longitudeNormalization) {
	aIsTrans := bboxIsTransmeridian(a)
	bIsTrans := bboxIsTransmeridian(b)
	aToBTrendsEast := a.West-b.East < b.West-a.East

	if !aIsTrans {
		*aNorm = normalizeNone
	} else if bIsTrans {
		*aNorm = normalizeEast
	} else if aToBTrendsEast {
		*aNorm = normalizeEast
	} else {
		*aNorm = normalizeWest
	}

	if !bIsTrans {
		*bNorm = normalizeNone
	} else if aIsTrans {
		*bNorm = normalizeEast
	} else if aToBTrendsEast {
		*bNorm = normalizeWest
	} else {
		*bNorm = normalizeEast
	}
}
