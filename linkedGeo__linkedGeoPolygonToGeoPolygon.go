package h3

// linkedGeoPolygonToGeoPolygon converts a single LinkedGeoPolygon
// (outer loop + holes) to a GeoPolygon. The output's geoloop and holes
// are allocated; on failure, all partial allocations are freed (in Go:
// zeroed for the garbage collector) and `out` is left in a clean
// (zeroed) state.
// Ported from H3 C: linkedGeo.c::linkedGeoPolygonToGeoPolygon.
func linkedGeoPolygonToGeoPolygon(linked *linkedGeoPolygon, out *GeoPolygon) h3Error {
	// Convert outer loop
	firstLoop := linked.First
	err := linkedGeoLoopToGeoLoop(firstLoop, &out.GeoLoop)
	if err != eSuccess {
		return err
	}

	// Count and convert holes
	numHoles := countLinkedLoops(linked) - 1
	if numHoles > 0 {
		holes := make([]GeoLoop, numHoles)
		out.Holes = holes

		loop := firstLoop.Next
		for i := int32(0); loop != nil && i < numHoles; i++ {
			// ALWAYS(i < numHoles) in C.
			err = linkedGeoLoopToGeoLoop(loop, &holes[i])
			if err != eSuccess {
				destroyGeoPolygon(out)
				return err
			}
			loop = loop.Next
		}
	}

	return eSuccess
}
