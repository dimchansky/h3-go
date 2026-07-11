package h3

import "math"

// bboxWidthRads returns bbox width in radians. Port of bbox.c::bboxWidthRads
// Ported from H3 C: bbox.c::bboxWidthRads.
func bboxWidthRads(b *BBox) float64 {
	if bboxIsTransmeridian(b) {
		return b.East.Rad() - b.West.Rad() + 2*math.Pi
	}
	return b.East.Rad() - b.West.Rad()
}

// bboxHeightRads returns bbox height in radians. Port of bbox.c::bboxHeightRads
// Ported from H3 C: bbox.c::bboxHeightRads.
func bboxHeightRads(b *BBox) float64 { return b.North.Rad() - b.South.Rad() }
