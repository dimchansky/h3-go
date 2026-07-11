package h3

import "math"

// cellToBBox calculates a bounding box for a given cell. If coverChildren is true,
// the bbox will be guaranteed to contain its children at any finer resolution.
// Note that no guarantee is provided as to the level of accuracy, and the bounding
// box may have a significant margin of error.
//
// Ported from H3 C: polyfill.c::cellToBBox.
func cellToBBox(cell H3Index, coverChildren bool) (BBox, H3Error) {
	res := getResolution(cell)
	var out BBox

	if res == 0 {
		baseCell := getBaseCell(cell)
		if baseCell < 0 || baseCell >= NUM_BASE_CELLS {
			return out, E_CELL_INVALID
		}
		out = RES0_BBOXES[baseCell]
	} else {
		var center LatLng
		centerErr := cellToLatLng(cell, &center)
		if centerErr != E_SUCCESS {
			return out, centerErr
		}
		lngRatio := 1.0 / math.Cos(float64(center.Lat))
		out.North = center.Lat + Angle(MAX_EDGE_LENGTH_RADS[res])
		out.South = center.Lat - Angle(MAX_EDGE_LENGTH_RADS[res])
		out.East = center.Lng + Angle(MAX_EDGE_LENGTH_RADS[res]*lngRatio)
		out.West = center.Lng - Angle(MAX_EDGE_LENGTH_RADS[res]*lngRatio)
	}

	// Buffer the bounding box to cover children. Call this even if no buffering
	// is required in order to normalize the bbox to lat/lng bounds
	scaleFactor := CELL_SCALE_FACTOR
	if coverChildren {
		scaleFactor = CHILD_SCALE_FACTOR
	}
	scaleBBox(&out, scaleFactor)

	// Cell that contains the north pole
	if cell == NORTH_POLE_CELLS[res] {
		out.North = PiOver2
	}

	// Cell that contains the south pole
	if cell == SOUTH_POLE_CELLS[res] {
		out.South = -PiOver2
	}

	// If we contain a pole, expand the longitude to include the full domain,
	// effectively making the bbox a circle around the pole.
	if out.North == PiOver2 || out.South == -PiOver2 {
		out.East = Pi
		out.West = -Pi
	}

	return out, E_SUCCESS
}
