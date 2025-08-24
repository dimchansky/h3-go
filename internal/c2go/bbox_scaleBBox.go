package c2go

import "math"

// scaleBBox scales a given bounding box by some factor.
// Scales both width and height by the factor, rather than scaling area, which will scale at scale^2.
// This function may not behave well for extreme values, and should be used
// within a reasonable domain, and does not guarantee reasonable results for extreme values.
// Ported from H3 C: bbox.c::scaleBBox
func scaleBBox(bbox *BBox, scale float64) {
	width := bboxWidthRads(bbox)
	height := bboxHeightRads(bbox)
	widthBuffer := (width*scale - width) * 0.5
	heightBuffer := (height*scale - height) * 0.5

	// Scale north and south, clamping to latitude domain
	bbox.North += heightBuffer
	if bbox.North > math.Pi/2 {
		bbox.North = math.Pi / 2
	}
	bbox.South -= heightBuffer
	if bbox.South < -math.Pi/2 {
		bbox.South = -math.Pi / 2
	}

	// Scale east and west, clamping to longitude domain
	bbox.East += widthBuffer
	if bbox.East > math.Pi {
		bbox.East -= 2 * math.Pi
	}
	if bbox.East < -math.Pi {
		bbox.East += 2 * math.Pi
	}
	bbox.West -= widthBuffer
	if bbox.West > math.Pi {
		bbox.West -= 2 * math.Pi
	}
	if bbox.West < -math.Pi {
		bbox.West += 2 * math.Pi
	}
}
