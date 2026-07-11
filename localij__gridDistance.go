package h3

// gridDistance calculates the grid distance between two H3 indexes.
// This function may fail to find the distance between two indexes, for
// example if they are very far apart. It may also fail when finding
// distances for indexes on opposite sides of a pentagon.
// Ported from H3 C: localij.c::gridDistance.
func gridDistance(origin h3Index, index h3Index, out *int64) h3Error {
	var originIjk, h3Ijk coordIJK

	originError := cellToLocalIjk(origin, origin, &originIjk)
	if originError != 0 {
		return originError
	}

	destError := cellToLocalIjk(origin, index, &h3Ijk)
	if destError != 0 {
		return destError
	}

	*out = int64(ijkDistance(&originIjk, &h3Ijk))
	return 0 // eSuccess
}
