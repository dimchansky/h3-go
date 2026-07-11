package h3

// cellToVertex gets a single vertex for a given cell, as an H3 index, or
// H3_NULL if the vertex is invalid.
// Ported from H3 C: vertex.c::cellToVertex.
func cellToVertex(cell H3Index, vertexNum int32, out *H3Index) H3Error {
	cellIsPentagon := isPentagon(cell)
	var cellNumVerts int32
	if cellIsPentagon {
		cellNumVerts = NUM_PENT_VERTS
	} else {
		cellNumVerts = NUM_HEX_VERTS
	}
	res := int32(getResolution(cell))

	// Check for invalid vertexes
	if vertexNum < 0 || vertexNum > cellNumVerts-1 {
		return E_DOMAIN
	}

	// Default the owner and vertex number to the input cell
	owner := cell
	ownerVertexNum := vertexNum

	// Determine the owner, looking at the three cells that share the vertex.
	// By convention, the owner is the cell with the lowest numerical index.

	// If the cell is the center child of its parent, it will always have
	// the lowest index of any neighbor, so we can skip determining the owner
	if res == 0 || getIndexDigit(cell, res) != int32(CENTER_DIGIT) {
		// Get the left neighbor of the vertex, with its rotations
		left := directionForVertexNum(cell, vertexNum)
		if left == INVALID_DIGIT {
			return E_FAILED
		}
		lRotations := int32(0)
		var leftNeighbor H3Index
		leftNeighborError := h3NeighborRotations(cell, left, &lRotations, &leftNeighbor)
		if leftNeighborError != E_SUCCESS {
			return leftNeighborError
		}
		// Set to owner if lowest index
		if leftNeighbor < owner {
			owner = leftNeighbor
		}

		// As above, skip the right neighbor if the left is known lowest
		if res == 0 || getIndexDigit(leftNeighbor, res) != int32(CENTER_DIGIT) {
			// Get the right neighbor of the vertex, with its rotations
			// Note that vertex - 1 is the right side, as vertex numbers are CCW
			right := directionForVertexNum(cell, (vertexNum-1+cellNumVerts)%cellNumVerts)
			// This case should be unreachable; invalid verts fail earlier
			if right == INVALID_DIGIT {
				return E_FAILED
			}
			rRotations := int32(0)
			var rightNeighbor H3Index
			rightNeighborError := h3NeighborRotations(cell, right, &rRotations, &rightNeighbor)
			if rightNeighborError != E_SUCCESS {
				return rightNeighborError
			}
			// Set to owner if lowest index
			if rightNeighbor < owner {
				owner = rightNeighbor
				var dir Direction
				if isPentagon(owner) {
					dir = _directionForNeighbor(owner, cell)
				} else {
					dir = directions[(revNeighborDirectionsHex[right]+rRotations)%NUM_HEX_VERTS]
				}
				ownerVertexNum = vertexNumForDirection(owner, dir)
			}
		}

		// Determine the vertex number for the left neighbor
		if owner == leftNeighbor {
			ownerIsPentagon := isPentagon(owner)
			var dir Direction
			if ownerIsPentagon {
				dir = _directionForNeighbor(owner, cell)
			} else {
				dir = directions[(revNeighborDirectionsHex[left]+lRotations)%NUM_HEX_VERTS]
			}

			// For the left neighbor, we need the second vertex of the
			// edge, which may involve looping around the vertex nums
			ownerVertexNum = vertexNumForDirection(owner, dir) + 1
			if ownerVertexNum == NUM_HEX_VERTS ||
				(ownerIsPentagon && ownerVertexNum == NUM_PENT_VERTS) {
				ownerVertexNum = 0
			}
		}
	}

	// Create the vertex index
	vertex := owner
	vertex = setMode(vertex, H3_VERTEX_MODE)
	vertex = setReservedBits(vertex, ownerVertexNum)
	*out = vertex

	return E_SUCCESS
}
