package h3

import (
	"math"
)

// bboxFromLinkedGeoLoop creates a bounding box from a simple polygon loop.
// Known limitations:
//   - Does not support polygons with two adjacent points > 180 degrees of
//     longitude apart. These will be interpreted as crossing the antimeridian.
//   - Does not currently support polygons containing a pole.
//
// Ported from H3 C: polygonAlgos.h::GENERIC_LOOP_ALGO(bboxFrom) -> bboxFromLinkedGeoLoop.
func bboxFromLinkedGeoLoop(loop *linkedGeoLoop, bb *bbox) {
	// Early exit if there are no vertices
	if loop.First == nil {
		*bb = bbox{}
		return
	}

	bb.South = Rad(math.MaxFloat64)
	bb.West = Rad(math.MaxFloat64)
	bb.North = Rad(-math.MaxFloat64)
	bb.East = Rad(-math.MaxFloat64)
	minPosLng := Rad(math.MaxFloat64)
	maxNegLng := Rad(-math.MaxFloat64)
	isTransmeridian := false

	// Iterate through linked coordinates
	currentCoord := loop.First
	for currentCoord != nil {
		coord := currentCoord.Vertex

		// Get next coordinate, wrapping to first if at end
		var next LatLng
		if currentCoord.Next != nil {
			next = currentCoord.Next.Vertex
		} else {
			// Wrap to first coordinate to close the loop
			if loop.First != nil {
				next = loop.First.Vertex
			} else {
				break // Safety check
			}
		}

		lat := coord.Lat
		lng := coord.Lng

		if lat < bb.South {
			bb.South = lat
		}
		if lng < bb.West {
			bb.West = lng
		}
		if lat > bb.North {
			bb.North = lat
		}
		if lng > bb.East {
			bb.East = lng
		}

		// Save the min positive and max negative longitude for
		// use in the transmeridian case
		if lng > 0 && lng < minPosLng {
			minPosLng = lng
		}
		if lng < 0 && lng > maxNegLng {
			maxNegLng = lng
		}

		// Check for arcs > 180 degrees longitude, flagging as transmeridian
		if (lng - next.Lng).Abs() > Pi {
			isTransmeridian = true
		}

		currentCoord = currentCoord.Next

		// Break when we've completed the loop
		if currentCoord == loop.First {
			break
		}
	}

	// Swap east and west if transmeridian
	if isTransmeridian {
		bb.East = maxNegLng
		bb.West = minPosLng
	}
}
