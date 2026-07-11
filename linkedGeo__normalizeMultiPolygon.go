package h3

// normalizeMultiPolygon normalizes a linkedGeoPolygon in-place into a structure following GeoJSON
// MultiPolygon rules. Each polygon must have exactly one outer loop, which must be first in the list,
// followed by any holes. Holes are identified by winding order (holes are clockwise), which is
// guaranteed by the h3SetToVertexGraph algorithm.
//
// Technical details:
// - Input is assumed to be a single polygon including all loops to normalize
// - Returns eFailed if root has multiple polygons (root.Next != nil)
// - Early exit with eSuccess if there's only one loop or fewer
// - Separates clockwise loops (holes) from counter-clockwise loops (outer polygons)
// - Assigns each hole to its appropriate parent polygon using findPolygonForHole
// - Outer loops become new polygons, holes are added as inner loops to their containers
// - Memory management: In Go, we don't need explicit cleanup like C's destroyLinkedGeoLoop
//
// Ported from H3 C: linkedGeo.c::normalizeMultiPolygon.
func normalizeMultiPolygon(root *linkedGeoPolygon) h3Error {
	// Handle nil input
	if root == nil {
		return eFailed
	}

	// We assume that the input is a single polygon with loops;
	// if it has multiple polygons, don't touch it
	if root.Next != nil {
		return eFailed
	}

	// Count loops, exiting early if there's only one
	loopCount := countLinkedLoops(root)
	if loopCount <= 1 {
		return eSuccess
	}

	resultCode := eSuccess
	var polygon *linkedGeoPolygon
	var next *linkedGeoLoop
	innerCount := 0
	outerCount := 0

	// Create an array to hold all of the inner loops. Note that
	// this array will never be full, as there will always be fewer
	// inner loops than outer loops.
	innerLoops := make([]*linkedGeoLoop, loopCount)

	// Create an array to hold the bounding boxes for the outer loops
	bboxes := make([]bbox, loopCount)

	// Get the first loop and unlink it from root
	loop := root.First
	*root = linkedGeoPolygon{} // Zero out the root

	// Iterate over all loops, moving inner loops into an array and
	// assigning outer loops to new polygons
	for loop != nil {
		if isClockwiseLinkedGeoLoop(loop) {
			innerLoops[innerCount] = loop
			innerCount++
		} else {
			if polygon == nil {
				polygon = root
			} else {
				polygon = addNewLinkedPolygon(polygon)
			}
			addLinkedLoop(polygon, loop)
			bboxFromLinkedGeoLoop(loop, &bboxes[outerCount])
			outerCount++
		}
		// get the next loop and unlink it from this one
		next = loop.Next
		loop.Next = nil
		loop = next
	}

	// Find polygon for each inner loop and assign the hole to it
	for i := 0; i < innerCount; i++ {
		// Create slice views for bboxes to pass to findPolygonForHole
		bboxPtrs := make([]*bbox, outerCount)
		for j := 0; j < outerCount; j++ {
			bboxPtrs[j] = &bboxes[j]
		}

		// Note: findPolygonForHole expects the head of a linked list of polygons
		// After the normalization, root should contain the first polygon in the chain
		polygon = findPolygonForHole(innerLoops[i], root, bboxPtrs)
		if polygon != nil {
			addLinkedLoop(polygon, innerLoops[i])
		} else {
			// If we can't find a polygon (possible with invalid input), then
			// we need to handle the orphaned hole. In Go, we don't have explicit
			// memory management like destroyLinkedGeoLoop + free, but we still
			// set the result code to indicate failure.
			// The GC will clean up the unlinked loop when it goes out of scope.
			resultCode = eFailed
		}
	}

	return resultCode
}
