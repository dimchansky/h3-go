package c2go

// vertexRotations gets the number of CCW rotations of the cell's vertex numbers
// compared to the directional layout of its neighbors.
// Ported from H3 C: vertex.c::vertexRotations
func vertexRotations(cell H3Index, out *int32) H3Error {
	// Get the face and other info for the origin
	var fijk FaceIJK
	err := _h3ToFaceIjk(cell, &fijk)
	if err != E_SUCCESS {
		return err
	}
	baseCell := getBaseCellNumber(cell)
	cellLeadingDigit := _h3LeadingNonZeroDigit(cell)

	// get the base cell face
	var baseFijk FaceIJK
	_baseCellToFaceIjk(baseCell, &baseFijk)

	ccwRot60 := _baseCellToCCWrot60(baseCell, fijk.Face)

	if _isBaseCellPentagon(baseCell) {
		// Find the appropriate direction-to-face mapping
		var dirFaces PentagonDirectionFaces
		// We never hit the end condition
		p := int32(0)
		for ; p < NUM_PENTAGONS; p++ {
			if pentagonDirectionFaces[p].baseCell == baseCell {
				dirFaces = pentagonDirectionFaces[p]
				break
			}
		}
		if p == NUM_PENTAGONS {
			return E_FAILED
		}

		// additional CCW rotation for polar neighbors or IK neighbors
		if fijk.Face != baseFijk.Face &&
			(_isBaseCellPolarPentagon(baseCell) ||
				fijk.Face == dirFaces.faces[IK_AXES_DIGIT-DIRECTION_INDEX_OFFSET]) {
			ccwRot60 = (ccwRot60 + 1) % 6
		}

		// Check whether the cell crosses a deleted pentagon subsequence
		if cellLeadingDigit == int32(JK_AXES_DIGIT) &&
			fijk.Face == dirFaces.faces[IK_AXES_DIGIT-DIRECTION_INDEX_OFFSET] {
			// Crosses from JK to IK: Rotate CW
			ccwRot60 = (ccwRot60 + 5) % 6
		} else if cellLeadingDigit == int32(IK_AXES_DIGIT) &&
			fijk.Face == dirFaces.faces[JK_AXES_DIGIT-DIRECTION_INDEX_OFFSET] {
			// Crosses from IK to JK: Rotate CCW
			ccwRot60 = (ccwRot60 + 1) % 6
		}
	}
	*out = ccwRot60
	return E_SUCCESS
}
