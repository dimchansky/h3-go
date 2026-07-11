package h3

// vertexRotations gets the number of CCW rotations of the cell's vertex numbers
// compared to the directional layout of its neighbors.
// Ported from H3 C: vertex.c::vertexRotations.
func vertexRotations(cell h3Index, out *int32) h3Error {
	// Get the face and other info for the origin
	var fijk faceIJK
	err := _h3ToFaceIjk(cell, &fijk)
	if err != eSuccess {
		return err
	}
	baseCell := getBaseCell(cell)
	cellLeadingDigit := _h3LeadingNonZeroDigit(cell)

	// get the base cell face
	var baseFijk faceIJK
	_baseCellToFaceIjk(baseCell, &baseFijk)

	ccwRot60 := _baseCellToCCWrot60(baseCell, fijk.Face)

	if _isBaseCellPentagon(baseCell) {
		// Find the appropriate direction-to-face mapping
		var dirFaces pentagonDirectionFacesEntry
		// We never hit the end condition
		p := int32(0)
		for ; p < numPentagons; p++ {
			if pentagonDirectionFaces[p].baseCell == baseCell {
				dirFaces = pentagonDirectionFaces[p]
				break
			}
		}
		if p == numPentagons {
			return eFailed
		}

		// additional CCW rotation for polar neighbors or IK neighbors
		if fijk.Face != baseFijk.Face &&
			(_isBaseCellPolarPentagon(baseCell) ||
				fijk.Face == dirFaces.faces[ikAxesDigit-directionIndexOffset]) {
			ccwRot60 = (ccwRot60 + 1) % 6
		}

		// Check whether the cell crosses a deleted pentagon subsequence
		if cellLeadingDigit == int32(jkAxesDigit) &&
			fijk.Face == dirFaces.faces[ikAxesDigit-directionIndexOffset] {
			// Crosses from quadJK to IK: Rotate CW
			ccwRot60 = (ccwRot60 + 5) % 6
		} else if cellLeadingDigit == int32(ikAxesDigit) &&
			fijk.Face == dirFaces.faces[jkAxesDigit-directionIndexOffset] {
			// Crosses from IK to quadJK: Rotate CCW
			ccwRot60 = (ccwRot60 + 1) % 6
		}
	}
	*out = ccwRot60
	return eSuccess
}
