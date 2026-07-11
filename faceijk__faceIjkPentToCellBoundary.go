package h3

// _faceIjkPentToCellBoundary finds the boundary vertices for the given pentagon face-ijk cell.
// Pentagon cells have 5 vertices and all Class III pentagon edges cross icosahedron edges.
// Class II pentagons have vertices directly on the edges with no edge intersections.
// Ported from H3 C: faceijk.c::_faceIjkPentToCellBoundary.
func _faceIjkPentToCellBoundary(h *faceIJK, res int32, start int32, length int32, g *CellBoundary) {
	adjRes := res
	centerIJK := *h
	var fijkVerts [numPentVerts]faceIJK
	_faceIjkPentToVerts(&centerIJK, &adjRes, fijkVerts[:])

	// If we're returning the entire loop, we need one more iteration in case
	// of a distortion vertex on the last edge
	additionalIteration := int32(0)
	if length == numPentVerts {
		additionalIteration = 1
	}

	// convert each vertex to lat/lng
	// adjust the face of each vertex as appropriate and introduce
	// edge-crossing vertices as needed
	g.NumVerts = 0
	var lastFijk faceIJK

	for vert := start; vert < start+length+additionalIteration; vert++ {
		v := vert % numPentVerts
		fijk := fijkVerts[v]
		_adjustPentVertOverage(&fijk, adjRes)

		// all Class III pentagon edges cross icosa edges
		// note that Class II pentagons have vertices on the edge,
		// not edge intersections
		if isResolutionClassIII(res) && vert > start {
			// find hex2d of the two vertexes on the last face
			tmpFijk := fijk
			var orig2d0 vec2d
			_ijkToHex2d(&lastFijk.Coord, &orig2d0)

			currentToLastDir := adjacentFaceDir[tmpFijk.Face][lastFijk.Face]
			fijkOrient := &faceNeighbors[tmpFijk.Face][currentToLastDir]

			tmpFijk.Face = fijkOrient.Face
			ijk := &tmpFijk.Coord

			// rotate and translate for adjacent face
			for i := int32(0); i < fijkOrient.CcwRot60; i++ {
				_ijkRotate60ccw(ijk)
			}
			transVec := fijkOrient.Translate
			_ijkScale(&transVec, unitScaleByCIIres[adjRes]*3)
			_ijkAdd(ijk, &transVec, ijk)
			_ijkNormalize(ijk)

			var orig2d1 vec2d
			_ijkToHex2d(ijk, &orig2d1)

			// find the appropriate icosa face edge vertexes
			maxDim := float64(maxDimByCIIres[adjRes])
			v0 := vec2d{3.0 * maxDim, 0.0}
			v1 := vec2d{-1.5 * maxDim, 3.0 * mSqrt32 * maxDim}
			v2 := vec2d{-1.5 * maxDim, -3.0 * mSqrt32 * maxDim}

			var edge0, edge1 *vec2d
			switch adjacentFaceDir[tmpFijk.Face][fijk.Face] {
			case quadIJ:
				edge0 = &v0
				edge1 = &v1
			case quadJK:
				edge0 = &v1
				edge1 = &v2
			case quadKI:
				fallthrough
			default:
				edge0 = &v2
				edge1 = &v0
			}

			// find the intersection and add the lat/lng point to the result
			inter := _v2dIntersect(&orig2d0, &orig2d1, edge0, edge1)

			// Ensure we have space in the boundary
			if len(g.Verts) <= int(g.NumVerts) {
				// Extend the slice if needed
				newVerts := make([]LatLng, g.NumVerts+1, (g.NumVerts+1)*2)
				copy(newVerts, g.Verts)
				g.Verts = newVerts
			}
			_hex2dToGeo(&inter, tmpFijk.Face, adjRes, 1, &g.Verts[g.NumVerts])
			g.NumVerts++
		}

		// convert vertex to lat/lng and add to the result
		// vert == start + numPentVerts is only used to test for possible
		// intersection on last edge
		if vert < start+numPentVerts {
			// Ensure we have space in the boundary
			if len(g.Verts) <= int(g.NumVerts) {
				// Extend the slice if needed
				newVerts := make([]LatLng, g.NumVerts+1, (g.NumVerts+1)*2)
				copy(newVerts, g.Verts)
				g.Verts = newVerts
			}
			var vec vec2d
			_ijkToHex2d(&fijk.Coord, &vec)
			_hex2dToGeo(&vec, fijk.Face, adjRes, 1, &g.Verts[g.NumVerts])
			g.NumVerts++
		}
		lastFijk = fijk
	}
}
