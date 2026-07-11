package h3

// getIcosahedronFaces finds all icosahedron faces intersected by a given H3 index.
//
// For a given h3Index, returns the list of icosahedron faces intersected by it.
// This can be used to accelerate other algorithms. For class II pentagons, the
// function recursively checks a direct child instead of the vertices, as pentagon
// vertices lie on icosahedron edges.
//
// It can happen that faces are returned that don't actually intersect the given
// h3Index. This is only false positives, though, so it's up to the caller to
// filter out invalid values.
//
// Ported from H3 C: h3Index.c::getIcosahedronFaces.
func getIcosahedronFaces(h h3Index, out []int32) h3Error {
	res := getResolution(h)
	isPent := isPentagon(h)

	// We can't use the vertex-based approach here for class II pentagons,
	// because all their vertices are on the icosahedron edges. Their
	// direct child pentagons cross the same faces, so use those instead.
	if isPent && !isResolutionClassIII(res) {
		// Note that this would not work for res 15, but this is only run on
		// Class II pentagons, it should never be invoked for a res 15 index.
		childPentagon := makeDirectChild(h, 0)
		return getIcosahedronFaces(childPentagon, out)
	}

	// convert to faceIJK
	var fijk faceIJK
	err := _h3ToFaceIjk(h, &fijk)
	if err != eSuccess {
		return err
	}

	// Get all vertices as faceIJK addresses. For simplicity, always
	// initialize the array with 6 verts, ignoring the last one for pentagons
	var fijkVerts [numHexVerts]faceIJK
	var vertexCount int32
	resCopy := res // Make a copy since the functions may modify it
	if isPent {
		vertexCount = numPentVerts
		_faceIjkPentToVerts(&fijk, &resCopy, fijkVerts[:])
	} else {
		resCopy = res // Reset for hex case
		vertexCount = numHexVerts
		_faceIjkToVerts(&fijk, &resCopy, fijkVerts[:])
	}

	// We may not use all of the slots in the output array,
	// so fill with invalid values to indicate unused slots
	var faceCount int32
	maxFaceCountError := maxFaceCount(h, &faceCount)
	if maxFaceCountError != eSuccess {
		return maxFaceCountError
	}
	for i := int32(0); i < faceCount; i++ {
		out[i] = invalidFace
	}

	// add each vertex face, using the output array as a hash set
	for i := int32(0); i < vertexCount; i++ {
		vert := &fijkVerts[i]
		// Adjust overage, determining whether this vertex is
		// on another face
		if isPent {
			_adjustPentVertOverage(vert, resCopy)
		} else {
			_adjustOverageClassII(vert, resCopy, false, true)
		}
		// Save the face to the output array
		face := vert.Face
		pos := int32(0)
		// Find the first empty output position, or the first position
		// matching the current face
		for out[pos] != invalidFace && out[pos] != face {
			pos++
			if pos >= faceCount {
				// Mismatch between the heuristic used in maxFaceCount and
				// calculation here - indicates an invalid index.
				return eFailed
			}
		}
		out[pos] = face
	}
	return eSuccess
}
