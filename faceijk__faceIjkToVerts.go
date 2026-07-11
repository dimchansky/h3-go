package h3

// _faceIjkToVerts generates the cell vertices for an origin-centered cell.
// The vertices are generated in counterclockwise order from the i-axis.
// This function supports both Class II and Class III resolutions.
// Note: This function may modify the resolution value for Class III cells.
// Ported from H3 C: faceijk.c::_faceIjkToVerts.
func _faceIjkToVerts(fijk *FaceIJK, res *int32, fijkVerts []FaceIJK) {
	// the vertexes of an origin-centered cell in a Class II resolution on a
	// substrate grid with aperture sequence 33r. The aperture 3 gets us the
	// vertices, and the 3r gets us back to Class II.
	// vertices listed ccw from the i-axes
	vertsCII := [NUM_HEX_VERTS]CoordIJK{
		{2, 1, 0}, // 0
		{1, 2, 0}, // 1
		{0, 2, 1}, // 2
		{0, 1, 2}, // 3
		{1, 0, 2}, // 4
		{2, 0, 1}, // 5
	}

	// the vertexes of an origin-centered cell in a Class III resolution on a
	// substrate grid with aperture sequence 33r7r. The aperture 3 gets us the
	// vertices, and the 3r7r gets us to Class II.
	// vertices listed ccw from the i-axes
	vertsCIII := [NUM_HEX_VERTS]CoordIJK{
		{5, 4, 0}, // 0
		{1, 5, 0}, // 1
		{0, 5, 4}, // 2
		{0, 1, 5}, // 3
		{4, 0, 5}, // 4
		{5, 0, 1}, // 5
	}

	// get the correct set of substrate vertices for this resolution
	var verts *[NUM_HEX_VERTS]CoordIJK
	if isResolutionClassIII(*res) {
		verts = &vertsCIII
	} else {
		verts = &vertsCII
	}

	// Make a copy of the input fijk since we'll be modifying it
	centerFijk := *fijk

	// adjust the center point to be in an aperture 33r substrate grid
	// these should be composed for speed
	_downAp3(&centerFijk.Coord)
	_downAp3r(&centerFijk.Coord)

	// if res is Class III we need to add a cw aperture 7 to get to
	// icosahedral Class II
	if isResolutionClassIII(*res) {
		_downAp7r(&centerFijk.Coord)
		*res += 1
	}

	// The center point is now in the same substrate grid as the origin
	// cell vertices. Add the center point substrate coordinates
	// to each vertex to translate the vertices to that cell.
	for v := 0; v < NUM_HEX_VERTS; v++ {
		fijkVerts[v].Face = centerFijk.Face
		_ijkAdd(&centerFijk.Coord, &verts[v], &fijkVerts[v].Coord)
		_ijkNormalize(&fijkVerts[v].Coord)
	}
}
