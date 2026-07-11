package h3

// scaleBBox scales a given bounding box by some factor.
// Scales both width and height by the factor, rather than scaling area, which will scale at scale^2.
// This function may not behave well for extreme values, and should be used
// within a reasonable domain, and does not guarantee reasonable results for extreme values.
// Ported from H3 C: bbox.c::scaleBBox.
func scaleBBox(bbox *bbox, scale float64) {
	width := bboxWidthRads(bbox)
	height := bboxHeightRads(bbox)
	widthBuffer := (width*scale - width) * 0.5
	heightBuffer := (height*scale - height) * 0.5

	// Scale north and south, clamping to latitude domain
	bbox.North = bbox.North + Rad(heightBuffer)
	if bbox.North > PiOver2 {
		bbox.North = PiOver2
	}
	bbox.South = bbox.South - Rad(heightBuffer)
	if bbox.South < -PiOver2 {
		bbox.South = -PiOver2
	}

	// Scale east and west, clamping to longitude domain
	bbox.East = bbox.East + Rad(widthBuffer)
	if bbox.East > Pi {
		bbox.East = bbox.East - TwoPi
	}
	if bbox.East < -Pi {
		bbox.East = bbox.East + TwoPi
	}
	bbox.West = bbox.West - Rad(widthBuffer)
	if bbox.West > Pi {
		bbox.West = bbox.West - TwoPi
	}
	if bbox.West < -Pi {
		bbox.West = bbox.West + TwoPi
	}
}
