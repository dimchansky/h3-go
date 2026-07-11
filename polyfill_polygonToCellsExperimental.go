package h3

// polygonToCellsExperimental takes a given GeoJSON-like data structure and preallocated
// memory, and fills it with the hexagons that are contained by the GeoJSON-like
// data structure. Polygons are considered in Cartesian space.
//
// This function uses an iterator-based approach to efficiently find cells within
// the polygon boundaries using various containment modes.
//
// Ported from H3 C: polyfill.c::polygonToCellsExperimental.
func polygonToCellsExperimental(polygon *GeoPolygon, res int32, flags uint32,
	size int64, out []h3Index) h3Error {
	if len(out) == 0 || int64(len(out)) < size {
		return eMemoryBounds
	}

	iter := iterInitPolygon(polygon, res, flags)
	if iter.Error != eSuccess {
		return iter.Error
	}

	var i int64
	for iter.Cell != h3Null {
		if i >= size {
			iterDestroyPolygon(&iter)
			return eMemoryBounds
		}
		out[i] = iter.Cell
		i++
		iterStepPolygon(&iter)
	}

	// Clean up the iterator
	iterDestroyPolygon(&iter)
	return iter.Error
}

// iterInitPolygon initializes a iterCellsPolygon struct representing the sequence of
// cells within the target polygon. The test for including edge cells is defined
// by the polyfill mode passed in the flags argument.
//
// Ported from H3 C: polyfill.c::iterInitPolygon.
func iterInitPolygon(polygon *GeoPolygon, res int32, flags uint32) iterCellsPolygon {
	// Create the sub-iterator for compact cells
	cellIter := iterInitPolygonCompact(polygon, res, flags)
	// Create the sub-iterator for children
	childIter := iterCellsChildren{}
	iterInitParent(cellIter.Cell, res, &childIter)

	iter := iterCellsPolygon{
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
// Ported from H3 C: polyfill.c::iterStepPolygon.
func iterStepPolygon(iter *iterCellsPolygon) {
	if iter.Cell == h3Null {
		return
	}

	// See if there are more children to output
	iterStepChild(&(iter.childIter))
	if iter.childIter.H != h3Null {
		iter.Cell = iter.childIter.H
		return
	}

	// Otherwise, increment the polyfill iterator
	iterStepPolygonCompact(&(iter.cellIter))
	if iter.cellIter.Cell != h3Null {
		iterInitParent(iter.cellIter.Cell, iter.cellIter.res, &(iter.childIter))
		iter.Cell = iter.childIter.H
		return
	}

	// All done, set to null and report errors if any
	iter.Cell = h3Null
	iter.Error = iter.cellIter.Error
}

// iterDestroyPolygon destroys an iterator, releasing any allocated memory.
// Iterators destroyed in this manner are safe to use but will always return h3Null.
//
// Ported from H3 C: polyfill.c::iterDestroyPolygon.
func iterDestroyPolygon(iter *iterCellsPolygon) {
	iterDestroyPolygonCompact(&(iter.cellIter))
	// null out the child iterator by passing h3Null
	iterInitParent(h3Null, 0, &(iter.childIter))
	iter.Cell = h3Null
	iter.Error = eSuccess
}

// iterInitPolygonCompact initializes a iterCellsPolygonCompact struct representing
// the sequence of compact cells within the target polygon.
//
// Ported from H3 C: polyfill.c::iterInitPolygonCompact.
func iterInitPolygonCompact(polygon *GeoPolygon, res int32, flags uint32) iterCellsPolygonCompact {
	iter := initIterPolygonCompact(polygon, res, flags)

	// Start the iterator by taking the first step.
	// This is necessary to have a valid value after initialization.
	iterStepPolygonCompact(&iter)

	return iter
}

// initIterPolygonCompact is the internal initialization function without stepping
//
// Ported from H3 C: polyfill.c::_iterInitPolygonCompact.
func initIterPolygonCompact(polygon *GeoPolygon, res int32, flags uint32) iterCellsPolygonCompact {
	iter := iterCellsPolygonCompact{
		// Initialize output properties. The first valid cell will be set in iterStep
		Cell:  baseCellNumToCell(0),
		Error: eSuccess,
		// Save input arguments
		polygon: polygon,
		res:     res,
		flags:   flags,
		bboxes:  nil,
		started: false,
	}

	if res < 0 || res > maxH3Res {
		iterErrorPolygonCompact(&iter, eResDomain)
		return iter
	}

	flagErr := validatePolygonFlags(flags)
	if flagErr != eSuccess {
		iterErrorPolygonCompact(&iter, flagErr)
		return iter
	}

	// Initialize bounding boxes for polygon and any holes
	numBBoxes := 1 + len(polygon.Holes)
	iter.bboxes = make([]bbox, numBBoxes)
	bboxesFromGeoPolygon(polygon, iter.bboxes)

	return iter
}

// iterErrorPolygonCompact sets an error state and cleans up the iterator.
func iterErrorPolygonCompact(iter *iterCellsPolygonCompact, errCode h3Error) {
	iterDestroyPolygonCompact(iter)
	iter.Error = errCode
}

// iterDestroyPolygonCompact destroys an iterator, releasing any allocated memory.
//
// Ported from H3 C: polyfill.c::iterDestroyPolygonCompact.
func iterDestroyPolygonCompact(iter *iterCellsPolygonCompact) {
	iter.Cell = h3Null
	iter.Error = eSuccess
	iter.polygon = nil
	iter.res = -1
	iter.flags = 0
	iter.bboxes = nil
}

// nextCell finds the next cell in the sequence of all cells to check in the iteration.
//
// Ported from H3 C: polyfill.c::nextCell.
func nextCell(cell h3Index) h3Index {
	res := getResolution(cell)
	for {
		// If this is a base cell, set to next base cell (or h3Null if done)
		if res == 0 {
			return baseCellNumToCell(getBaseCell(cell) + 1)
		}

		// Faster cellToParent when we know the resolution is valid
		// and we're only moving up one level
		parent := cell
		parent = setResolution(parent, res-1)
		parent = setIndexDigit(parent, res, int32(h3DigitMask))

		// If not the last sibling of parent, return next sibling
		digit := direction(getIndexDigit(cell, res))
		if digit < invalidDigit-1 {
			nextDigit := digit + 1
			if isPentagon(parent) && digit == centerDigit {
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
// Ported from H3 C: polyfill.c::iterStepPolygonCompact.
func iterStepPolygonCompact(iter *iterCellsPolygonCompact) {
	cell := iter.Cell

	// once the cell is h3Null, the iterator returns an infinite sequence of h3Null
	if cell == h3Null {
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

	mode := flagGetContainmentMode(iter.flags)

	for cell != h3Null {
		cellRes := getResolution(cell)

		// Target res: Do a fine-grained check
		if cellRes == iter.res {
			if mode == ContainmentCenter || mode == ContainmentOverlapping ||
				mode == ContainmentOverlappingBBox {
				// Check if the cell center is inside the polygon
				var center LatLng
				centerErr := cellToLatLng(cell, &center)
				if centerErr != eSuccess {
					iterErrorPolygonCompact(iter, centerErr)
					return
				}
				if pointInsidePolygon(*iter.polygon, iter.bboxes, &center) {
					// Set to next output
					iter.Cell = cell
					return
				}
			}
			if mode == ContainmentOverlapping ||
				mode == ContainmentOverlappingBBox {
				// For overlapping, we need to do a quick check to determine
				// whether the polygon is wholly contained by the cell. We
				// check the first polygon vertex, which if it is contained
				// could also mean we simply intersect.

				// Dereferencing verts[0] is safe because we check len above
				firstVertex := iter.polygon.GeoLoop[0]

				// We have to check whether the point is in the expected range
				// first, because out-of-bounds values will yield false
				// positives with latLngToCell
				if bboxContains(&validRangeBbox, &firstVertex) {
					var polygonCell h3Index
					polygonCellErr := latLngToCell(&firstVertex, cellRes, &polygonCell)
					if polygonCellErr != eSuccess {
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
			if mode == ContainmentFull || mode == ContainmentOverlapping ||
				mode == ContainmentOverlappingBBox {
				var boundary CellBoundary
				boundaryErr := cellToBoundary(cell, &boundary)
				if boundaryErr != eSuccess {
					iterErrorPolygonCompact(iter, boundaryErr)
					return
				}
				bbox, bboxErr := cellToBBox(cell, false)
				if bboxErr != eSuccess {
					// Should be unreachable - invalid cells would be caught in
					// the previous boundaryErr
					iterErrorPolygonCompact(iter, bboxErr)
					return
				}
				// Check if the cell is fully contained by the polygon
				if (mode == ContainmentFull ||
					mode == ContainmentOverlappingBBox) &&
					cellBoundaryInsidePolygon(*iter.polygon, iter.bboxes,
						&boundary, &bbox) {
					// Set to next output
					iter.Cell = cell
					return
				}
				// For overlap, we've already checked for center point inclusion
				// above; if that failed, we only need to check for line
				// intersection
				if (mode == ContainmentOverlapping ||
					mode == ContainmentOverlappingBBox) &&
					cellBoundaryCrossesPolygon(*iter.polygon, iter.bboxes, &boundary, &bbox) {
					// Set to next output
					iter.Cell = cell
					return
				}
			}
			if mode == ContainmentOverlappingBBox {
				// Get a bounding box containing all the cell's children, so
				// this can work for the max size calculation
				bbox, bboxErr := cellToBBox(cell, true)
				if bboxErr != eSuccess {
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
			if bboxErr != eSuccess {
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
				if childErr != eSuccess {
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
