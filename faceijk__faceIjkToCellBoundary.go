package h3

// _faceIjkToCellBoundary finds the boundary vertices for the given face-ijk cell.
// This function handles edge-crossing vertices for Class III resolutions when
// hexagon edges cross icosahedron face boundaries requiring interpolation.
// Ported from H3 C: faceijk.c::_faceIjkToCellBoundary.
func _faceIjkToCellBoundary(h *faceIJK, res int32, start int32, length int32, g *CellBoundary) {
	adjRes := res
	centerIJK := *h
	var fijkVerts [numHexVerts]faceIJK
	_faceIjkToVerts(&centerIJK, &adjRes, fijkVerts[:])

	// If we're returning the entire loop, we need one more iteration in case
	// of a distortion vertex on the last edge
	additionalIteration := int32(0)
	if length == numHexVerts {
		additionalIteration = 1
	}

	// convert each vertex to lat/lng
	// adjust the face of each vertex as appropriate and introduce
	// edge-crossing vertices as needed
	g.numVerts = 0
	lastFace := int32(-1)
	lastOverage := noOverage

	for vert := start; vert < start+length+additionalIteration; vert++ {
		v := vert % numHexVerts
		fijk := fijkVerts[v]
		pentLeading4 := false
		overage := _adjustOverageClassII(&fijk, adjRes, pentLeading4, true)

		/*
			Check for edge-crossing. Each face of the underlying icosahedron is a
			different projection plane. So if an edge of the hexagon crosses an
			icosahedron edge, an additional vertex must be introduced at that
			intersection point. Then each half of the cell edge can be projected
			to geographic coordinates using the appropriate icosahedron face
			projection. Note that Class II cell edges have vertices on the face
			edge, with no edge line intersections.
		*/
		if isResolutionClassIII(res) && vert > start &&
			fijk.Face != lastFace && lastOverage != faceEdge {
			// find hex2d of the two vertexes on original face
			lastV := (v + 5) % numHexVerts
			var orig2d0 vec2d
			_ijkToHex2d(&fijkVerts[lastV].Coord, &orig2d0)
			var orig2d1 vec2d
			_ijkToHex2d(&fijkVerts[v].Coord, &orig2d1)

			// find the appropriate icosa face edge vertexes
			maxDim := float64(maxDimByCIIres[adjRes])
			v0 := vec2d{3.0 * maxDim, 0.0}
			v1 := vec2d{-1.5 * maxDim, 3.0 * mSqrt32 * maxDim}
			v2 := vec2d{-1.5 * maxDim, -3.0 * mSqrt32 * maxDim}

			face2 := lastFace
			if lastFace == centerIJK.Face {
				face2 = fijk.Face
			}

			var edge0, edge1 *vec2d
			switch adjacentFaceDir[centerIJK.Face][face2] {
			case quadIJ:
				edge0 = &v0
				edge1 = &v1
			case quadJK:
				edge0 = &v1
				edge1 = &v2
			default: // quadKI case
				edge0 = &v2
				edge1 = &v0
			}

			// find the intersection and add the lat/lng point to the result
			inter := _v2dIntersect(&orig2d0, &orig2d1, edge0, edge1)

			/*
				If a point of intersection occurs at a hexagon vertex, then each
				adjacent hexagon edge will lie completely on a single icosahedron
				face, and no additional vertex is required.
			*/
			isIntersectionAtVertex := _v2dAlmostEquals(&orig2d0, &inter) ||
				_v2dAlmostEquals(&orig2d1, &inter)
			if !isIntersectionAtVertex {
				var v3 vec3d
				_hex2dToVec3(&inter, centerIJK.Face, adjRes, 1, &v3)
				g.verts[g.numVerts] = vec3ToLatLng(v3)
				g.numVerts++
			}
		}

		// convert vertex to lat/lng and add to the result
		// vert == start + numHexVerts is only used to test for possible
		// intersection on last edge
		if vert < start+numHexVerts {
			var vec vec2d
			_ijkToHex2d(&fijk.Coord, &vec)
			var v3 vec3d
			_hex2dToVec3(&vec, fijk.Face, adjRes, 1, &v3)
			g.verts[g.numVerts] = vec3ToLatLng(v3)
			g.numVerts++
		}
		lastFace = fijk.Face
		lastOverage = overage
	}
}
