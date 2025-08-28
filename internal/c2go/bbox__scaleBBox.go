package c2go

import (
	"github.com/dimchansky/h3-go/angle"
)

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
	bbox.North = bbox.North + angle.Rad(heightBuffer)
	if bbox.North > angle.PiOver2 {
		bbox.North = angle.PiOver2
	}
	bbox.South = bbox.South - angle.Rad(heightBuffer)
	if bbox.South < -angle.PiOver2 {
		bbox.South = -angle.PiOver2
	}

	// Scale east and west, clamping to longitude domain
	bbox.East = bbox.East + angle.Rad(widthBuffer)
	if bbox.East > angle.Pi {
		bbox.East = bbox.East - angle.TwoPi
	}
	if bbox.East < -angle.Pi {
		bbox.East = bbox.East + angle.TwoPi
	}
	bbox.West = bbox.West - angle.Rad(widthBuffer)
	if bbox.West > angle.Pi {
		bbox.West = bbox.West - angle.TwoPi
	}
	if bbox.West < -angle.Pi {
		bbox.West = bbox.West + angle.TwoPi
	}
}
