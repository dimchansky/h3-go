package h3

// cellToVertex gets a single vertex for a given cell, as an H3 index, or
// h3Null if the vertex is invalid.
// Ported from H3 C: vertex.c::cellToVertex.
func cellToVertex(cell h3Index, vertexNum int32, out *h3Index) h3Error {
	cellIsPentagon := isPentagon(cell)
	var cellNumVerts int32
	if cellIsPentagon {
		cellNumVerts = numPentVerts
	} else {
		cellNumVerts = numHexVerts
	}
	res := int32(getResolution(cell))

	// Check for invalid vertexes
	if vertexNum < 0 || vertexNum > cellNumVerts-1 {
		return eDomain
	}

	// Default the owner and vertex number to the input cell
	owner := cell
	ownerVertexNum := vertexNum

	// Determine the owner, looking at the three cells that share the vertex.
	// By convention, the owner is the cell with the lowest numerical index.

	// If the cell is the center child of its parent, it will always have
	// the lowest index of any neighbor, so we can skip determining the owner
	if res == 0 || h3GetIndexDigit(cell, res) != int32(centerDigit) {
		// Get the left neighbor of the vertex, with its rotations
		left := directionForVertexNum(cell, vertexNum)
		if left == invalidDigit {
			return eFailed
		}
		lRotations := int32(0)
		var leftNeighbor h3Index
		leftNeighborError := h3NeighborRotations(cell, left, &lRotations, &leftNeighbor)
		if leftNeighborError != eSuccess {
			return leftNeighborError
		}
		// Set to owner if lowest index
		if leftNeighbor < owner {
			owner = leftNeighbor
		}

		// As above, skip the right neighbor if the left is known lowest
		if res == 0 || h3GetIndexDigit(leftNeighbor, res) != int32(centerDigit) {
			// Get the right neighbor of the vertex, with its rotations
			// Note that vertex - 1 is the right side, as vertex numbers are CCW
			right := directionForVertexNum(cell, (vertexNum-1+cellNumVerts)%cellNumVerts)
			// This case should be unreachable; invalid verts fail earlier
			if right == invalidDigit {
				return eFailed
			}
			rRotations := int32(0)
			var rightNeighbor h3Index
			rightNeighborError := h3NeighborRotations(cell, right, &rRotations, &rightNeighbor)
			if rightNeighborError != eSuccess {
				return rightNeighborError
			}
			// Set to owner if lowest index
			if rightNeighbor < owner {
				owner = rightNeighbor
				var dir direction
				if isPentagon(owner) {
					dir = _directionForNeighbor(owner, cell)
				} else {
					dir = directions[(revNeighborDirectionsHex[right]+rRotations)%numHexVerts]
				}
				ownerVertexNum = vertexNumForDirection(owner, dir)
			}
		}

		// Determine the vertex number for the left neighbor
		if owner == leftNeighbor {
			ownerIsPentagon := isPentagon(owner)
			var dir direction
			if ownerIsPentagon {
				dir = _directionForNeighbor(owner, cell)
			} else {
				dir = directions[(revNeighborDirectionsHex[left]+lRotations)%numHexVerts]
			}

			// For the left neighbor, we need the second vertex of the
			// edge, which may involve looping around the vertex nums
			ownerVertexNum = vertexNumForDirection(owner, dir) + 1
			if ownerVertexNum == numHexVerts ||
				(ownerIsPentagon && ownerVertexNum == numPentVerts) {
				ownerVertexNum = 0
			}
		}
	}

	// Create the vertex index
	vertex := owner
	vertex = setMode(vertex, h3VertexMode)
	vertex = setReservedBits(vertex, ownerVertexNum)
	*out = vertex

	return eSuccess
}
