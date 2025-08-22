package c2go

import "math"

// bboxWidthRads returns bbox width in radians. Port of bbox.c::bboxWidthRads
func bboxWidthRads(b BBox) float64 {
    if bboxIsTransmeridian(b) {
        return b.East - b.West + 2*math.Pi
    }
    return b.East - b.West
}

// bboxHeightRads returns bbox height in radians. Port of bbox.c::bboxHeightRads
func bboxHeightRads(b BBox) float64 { return b.North - b.South }

