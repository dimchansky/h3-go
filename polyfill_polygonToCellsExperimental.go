package h3

// polygonToCellsExperimental takes a given GeoJSON-like data structure and preallocated
// memory, and fills it with the hexagons that are contained by the GeoJSON-like
// data structure. Polygons are considered in Cartesian space.
//
// This function uses an iterator-based approach to efficiently find cells within
// the polygon boundaries using various containment modes.
//
// Ported from H3 C: polyfill.c::polygonToCellsExperimental
func polygonToCellsExperimental(polygon *GeoPolygon, res int32, flags uint32,
	size int64, out []H3Index) H3Error {
	if len(out) == 0 || int64(len(out)) < size {
		return E_MEMORY_BOUNDS
	}

	iter := iterInitPolygon(polygon, res, flags)
	if iter.Error != E_SUCCESS {
		return iter.Error
	}

	var i int64
	for iter.Cell != H3_NULL {
		if i >= size {
			iterDestroyPolygon(&iter)
			return E_MEMORY_BOUNDS
		}
		out[i] = iter.Cell
		i++
		iterStepPolygon(&iter)
	}

	// Clean up the iterator
	iterDestroyPolygon(&iter)
	return iter.Error
}

// iterInitPolygon initializes a IterCellsPolygon struct representing the sequence of
// cells within the target polygon. The test for including edge cells is defined
// by the polyfill mode passed in the flags argument.
//
// Ported from H3 C: polyfill.c::iterInitPolygon
func iterInitPolygon(polygon *GeoPolygon, res int32, flags uint32) IterCellsPolygon {
	// Create the sub-iterator for compact cells
	cellIter := iterInitPolygonCompact(polygon, res, flags)
	// Create the sub-iterator for children
	childIter := IterCellsChildren{}
	iterInitParent(cellIter.Cell, res, &childIter)

	iter := IterCellsPolygon{
		Cell:      childIter.H,
		Error:     cellIter.Error,
		cellIter:  cellIter,
		childIter: childIter,
	}
	return iter
}

// iterStepPolygon increments the polyfill iterator, outputting the latest cell at the
// desired resolution.
//
// Ported from H3 C: polyfill.c::iterStepPolygon
func iterStepPolygon(iter *IterCellsPolygon) {
	if iter.Cell == H3_NULL {
		return
	}

	// See if there are more children to output
	iterStepChild(&(iter.childIter))
	if iter.childIter.H != H3_NULL {
		iter.Cell = iter.childIter.H
		return
	}

	// Otherwise, increment the polyfill iterator
	iterStepPolygonCompact(&(iter.cellIter))
	if iter.cellIter.Cell != H3_NULL {
		iterInitParent(iter.cellIter.Cell, iter.cellIter.res, &(iter.childIter))
		iter.Cell = iter.childIter.H
		return
	}

	// All done, set to null and report errors if any
	iter.Cell = H3_NULL
	iter.Error = iter.cellIter.Error
}

// iterDestroyPolygon destroys an iterator, releasing any allocated memory.
// Iterators destroyed in this manner are safe to use but will always return H3_NULL.
//
// Ported from H3 C: polyfill.c::iterDestroyPolygon
func iterDestroyPolygon(iter *IterCellsPolygon) {
	iterDestroyPolygonCompact(&(iter.cellIter))
	// null out the child iterator by passing H3_NULL
	iterInitParent(H3_NULL, 0, &(iter.childIter))
	iter.Cell = H3_NULL
	iter.Error = E_SUCCESS
}

// iterInitPolygonCompact initializes a IterCellsPolygonCompact struct representing
// the sequence of compact cells within the target polygon.
//
// Ported from H3 C: polyfill.c::iterInitPolygonCompact
func iterInitPolygonCompact(polygon *GeoPolygon, res int32, flags uint32) IterCellsPolygonCompact {
	iter := initIterPolygonCompact(polygon, res, flags)

	// Start the iterator by taking the first step.
	// This is necessary to have a valid value after initialization.
	iterStepPolygonCompact(&iter)

	return iter
}

// initIterPolygonCompact is the internal initialization function without stepping
//
// Ported from H3 C: polyfill.c::_iterInitPolygonCompact
func initIterPolygonCompact(polygon *GeoPolygon, res int32, flags uint32) IterCellsPolygonCompact {
	iter := IterCellsPolygonCompact{
		// Initialize output properties. The first valid cell will be set in iterStep
		Cell:  baseCellNumToCell(0),
		Error: E_SUCCESS,
		// Save input arguments
		polygon: polygon,
		res:     res,
		flags:   flags,
		bboxes:  nil,
		started: false,
	}

	if res < 0 || res > MAX_H3_RES {
		iterErrorPolygonCompact(&iter, E_RES_DOMAIN)
		return iter
	}

	flagErr := validatePolygonFlags(flags)
	if flagErr != E_SUCCESS {
		iterErrorPolygonCompact(&iter, flagErr)
		return iter
	}

	// Initialize bounding boxes for polygon and any holes
	numBBoxes := 1 + len(polygon.Holes)
	iter.bboxes = make([]BBox, numBBoxes)
	bboxesFromGeoPolygon(polygon, iter.bboxes)

	return iter
}

// iterErrorPolygonCompact sets an error state and cleans up the iterator
func iterErrorPolygonCompact(iter *IterCellsPolygonCompact, error H3Error) {
	iterDestroyPolygonCompact(iter)
	iter.Error = error
}

// iterDestroyPolygonCompact destroys an iterator, releasing any allocated memory.
//
// Ported from H3 C: polyfill.c::iterDestroyPolygonCompact
func iterDestroyPolygonCompact(iter *IterCellsPolygonCompact) {
	iter.Cell = H3_NULL
	iter.Error = E_SUCCESS
	iter.polygon = nil
	iter.res = -1
	iter.flags = 0
	iter.bboxes = nil
}

// nextCell finds the next cell in the sequence of all cells to check in the iteration.
//
// Ported from H3 C: polyfill.c::nextCell
func nextCell(cell H3Index) H3Index {
	res := getResolution(cell)
	for {
		// If this is a base cell, set to next base cell (or H3_NULL if done)
		if res == 0 {
			return baseCellNumToCell(getBaseCell(cell) + 1)
		}

		// Faster cellToParent when we know the resolution is valid
		// and we're only moving up one level
		parent := cell
		parent = setResolution(parent, res-1)
		parent = setIndexDigit(parent, res, int32(H3_DIGIT_MASK))

		// If not the last sibling of parent, return next sibling
		digit := Direction(getIndexDigit(cell, res))
		if digit < INVALID_DIGIT-1 {
			nextDigit := digit + 1
			if isPentagon(parent) && digit == CENTER_DIGIT {
				nextDigit = digit + 2 // Skip missing pentagon child
			}
			cell = setIndexDigit(cell, res, int32(nextDigit))
			return cell
		}
		// Move up to the parent for the next loop iteration
		res--
		cell = parent
	}
}

// iterStepPolygonCompact increments the polyfill iterator, running the polygon to cells algorithm.
//
// Ported from H3 C: polyfill.c::iterStepPolygonCompact
func iterStepPolygonCompact(iter *IterCellsPolygonCompact) {
	cell := iter.Cell

	// once the cell is H3_NULL, the iterator returns an infinite sequence of H3_NULL
	if cell == H3_NULL {
		return
	}

	// For the first step, we need to evaluate the current cell; after that, we
	// should start with the next cell.
	if iter.started {
		cell = nextCell(cell)
	} else {
		iter.started = true
	}

	// Short-circuit iteration for 0-vert polygon
	if len(iter.polygon.GeoLoop) == 0 {
		iterDestroyPolygonCompact(iter)
		return
	}

	mode := FLAG_GET_CONTAINMENT_MODE(iter.flags)

	for cell != H3_NULL {
		cellRes := getResolution(cell)

		// Target res: Do a fine-grained check
		if cellRes == iter.res {
			if mode == CONTAINMENT_CENTER || mode == CONTAINMENT_OVERLAPPING ||
				mode == CONTAINMENT_OVERLAPPING_BBOX {
				// Check if the cell center is inside the polygon
				var center LatLng
				centerErr := cellToLatLng(cell, &center)
				if centerErr != E_SUCCESS {
					iterErrorPolygonCompact(iter, centerErr)
					return
				}
				if pointInsidePolygon(*iter.polygon, iter.bboxes, &center) {
					// Set to next output
					iter.Cell = cell
					return
				}
			}
			if mode == CONTAINMENT_OVERLAPPING ||
				mode == CONTAINMENT_OVERLAPPING_BBOX {
				// For overlapping, we need to do a quick check to determine
				// whether the polygon is wholly contained by the cell. We
				// check the first polygon vertex, which if it is contained
				// could also mean we simply intersect.

				// Dereferencing verts[0] is safe because we check len above
				firstVertex := iter.polygon.GeoLoop[0]

				// We have to check whether the point is in the expected range
				// first, because out-of-bounds values will yield false
				// positives with latLngToCell
				if bboxContains(&VALID_RANGE_BBOX, &firstVertex) {
					var polygonCell H3Index
					polygonCellErr := latLngToCell(&firstVertex, cellRes, &polygonCell)
					if polygonCellErr != E_SUCCESS {
						// This should be unreachable with the bbox check
						iterErrorPolygonCompact(iter, polygonCellErr)
						return
					}
					if polygonCell == cell {
						// Set to next output
						iter.Cell = cell
						return
					}
				}
			}
			if mode == CONTAINMENT_FULL || mode == CONTAINMENT_OVERLAPPING ||
				mode == CONTAINMENT_OVERLAPPING_BBOX {
				var boundary CellBoundary
				boundaryErr := cellToBoundary(cell, &boundary)
				if boundaryErr != E_SUCCESS {
					iterErrorPolygonCompact(iter, boundaryErr)
					return
				}
				bbox, bboxErr := cellToBBox(cell, false)
				if bboxErr != E_SUCCESS {
					// Should be unreachable - invalid cells would be caught in
					// the previous boundaryErr
					iterErrorPolygonCompact(iter, bboxErr)
					return
				}
				// Check if the cell is fully contained by the polygon
				if (mode == CONTAINMENT_FULL ||
					mode == CONTAINMENT_OVERLAPPING_BBOX) &&
					cellBoundaryInsidePolygon(*iter.polygon, iter.bboxes,
						&boundary, &bbox) {
					// Set to next output
					iter.Cell = cell
					return
				}
				// For overlap, we've already checked for center point inclusion
				// above; if that failed, we only need to check for line
				// intersection
				if (mode == CONTAINMENT_OVERLAPPING ||
					mode == CONTAINMENT_OVERLAPPING_BBOX) &&
					cellBoundaryCrossesPolygon(*iter.polygon, iter.bboxes, &boundary, &bbox) {
					// Set to next output
					iter.Cell = cell
					return
				}
			}
			if mode == CONTAINMENT_OVERLAPPING_BBOX {
				// Get a bounding box containing all the cell's children, so
				// this can work for the max size calculation
				bbox, bboxErr := cellToBBox(cell, true)
				if bboxErr != E_SUCCESS {
					iterErrorPolygonCompact(iter, bboxErr)
					return
				}
				if bboxOverlapsBBox(&iter.bboxes[0], &bbox) {
					bboxBoundary := bboxToCellBoundary(&bbox)
					if bboxContainsBBox(&bbox, &iter.bboxes[0]) ||
						pointInsidePolygon(*iter.polygon, iter.bboxes,
							&bboxBoundary.Verts[0]) ||
						cellBoundaryCrossesPolygon(*iter.polygon,
							iter.bboxes, &bboxBoundary, &bbox) {
						iter.Cell = cell
						return
					}
				}
			}
		}

		// Coarser cell: Check the bounding box
		if cellRes < iter.res {
			// Get a bounding box for all of the cell's children
			bbox, bboxErr := cellToBBox(cell, true)
			if bboxErr != E_SUCCESS {
				iterErrorPolygonCompact(iter, bboxErr)
				return
			}
			if bboxOverlapsBBox(&iter.bboxes[0], &bbox) {
				// Quick check for possible containment
				if bboxContainsBBox(&iter.bboxes[0], &bbox) {
					bboxBoundary := bboxToCellBoundary(&bbox)
					// Do a fine-grained, more expensive check on the polygon
					if cellBoundaryInsidePolygon(*iter.polygon, iter.bboxes,
						&bboxBoundary, &bbox) {
						// Bounding box is fully contained, so all children are
						// included. Set to next output.
						iter.Cell = cell
						return
					}
				}
				// Otherwise, the intersecting bbox means we need to test all
				// children, starting with the first child
				child, childErr := cellToCenterChild(cell, cellRes+1)
				if childErr != E_SUCCESS {
					iterErrorPolygonCompact(iter, childErr)
					return
				}
				// Restart the loop with the child cell
				cell = child
				continue
			}
		}

		// Find the next cell in the sequence of all cells and continue
		cell = nextCell(cell)
	}
	// If we make it out of the loop, we're done
	iterDestroyPolygonCompact(iter)
}
