package c2go

// bboxNormalization determines longitude normalization for two bboxes.
// Ported from H3 C: bbox.c::bboxNormalization
func bboxNormalization(a, b BBox) (LongitudeNormalization, LongitudeNormalization) {
	aIsTrans := bboxIsTransmeridian(a)
	bIsTrans := bboxIsTransmeridian(b)
	aToBTrendsEast := a.West-b.East < b.West-a.East
	var aNorm, bNorm LongitudeNormalization
	if !aIsTrans {
		aNorm = NORMALIZE_NONE
	} else if bIsTrans {
		aNorm = NORMALIZE_EAST
	} else if aToBTrendsEast {
		aNorm = NORMALIZE_EAST
	} else {
		aNorm = NORMALIZE_WEST
	}
	if !bIsTrans {
		bNorm = NORMALIZE_NONE
	} else if aIsTrans {
		bNorm = NORMALIZE_EAST
	} else if aToBTrendsEast {
		bNorm = NORMALIZE_WEST
	} else {
		bNorm = NORMALIZE_EAST
	}
	return aNorm, bNorm
}
