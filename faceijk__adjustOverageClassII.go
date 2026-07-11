package h3

// _adjustOverageClassII adjusts a faceIJK address for overage to an adjacent face.
// This handles coordinate transformations when IJK coordinates exceed the face boundary.
// Returns the overage status: noOverage, faceEdge, or newFace.
// Ported from H3 C: faceijk.c::_adjustOverageClassII.
func _adjustOverageClassII(fijk *faceIJK, res int32, pentLeading4 bool, substrate bool) overage {
	overage := noOverage
	ijk := &fijk.Coord

	// Get the maximum dimension value; scale if a substrate grid
	maxDim := int32(maxDimByCIIres[res])
	if substrate {
		maxDim *= 3
	}

	// Check for overage
	if substrate && ijk.I+ijk.J+ijk.K == maxDim { // on edge
		overage = faceEdge
	} else if ijk.I+ijk.J+ijk.K > maxDim { // overage
		overage = newFace

		var fijkOrient *faceOrientIJK
		if ijk.K > 0 {
			if ijk.J > 0 { // jk "quadrant"
				fijkOrient = &faceNeighbors[fijk.Face][quadJK]
			} else { // ik "quadrant"
				fijkOrient = &faceNeighbors[fijk.Face][quadKI]

				// Adjust for the pentagonal missing sequence
				if pentLeading4 {
					// Translate origin to center of pentagon
					origin := coordIJK{maxDim, 0, 0}
					var tmp coordIJK
					_ijkSub(ijk, &origin, &tmp)
					// Rotate to adjust for the missing sequence
					_ijkRotate60cw(&tmp)
					// Translate the origin back to the center of the triangle
					_ijkAdd(&tmp, &origin, ijk)
				}
			}
		} else { // ij "quadrant"
			fijkOrient = &faceNeighbors[fijk.Face][quadIJ]
		}

		fijk.Face = fijkOrient.Face

		// Rotate and translate for adjacent face
		for i := int32(0); i < fijkOrient.CcwRot60; i++ {
			_ijkRotate60ccw(ijk)
		}

		transVec := fijkOrient.Translate
		unitScale := int32(unitScaleByCIIres[res])
		if substrate {
			unitScale *= 3
		}
		_ijkScale(&transVec, unitScale)
		_ijkAdd(ijk, &transVec, ijk)
		_ijkNormalize(ijk)

		// overage points on pentagon boundaries can end up on edges
		if substrate && ijk.I+ijk.J+ijk.K == maxDim { // on edge
			overage = faceEdge
		}
	}

	return overage
}
