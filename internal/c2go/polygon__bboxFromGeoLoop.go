package c2go

import "math"

// bboxFromGeoLoop computes a bounding box for a loop of coordinates.
// Mirrors H3's polygon.c::bboxFromGeoLoop behavior.
func bboxFromGeoLoop(loop []LatLng, bbox *BBox) {
	if len(loop) == 0 {
		*bbox = BBox{}
		return
	}
	south := math.MaxFloat64
	west := math.MaxFloat64
	north := -math.MaxFloat64
	east := -math.MaxFloat64
	minPosLng := math.MaxFloat64
	maxNegLng := -math.MaxFloat64
	isTrans := false
	for i := 0; i < len(loop); i++ {
		coord := loop[i]
		next := loop[(i+1)%len(loop)]
		lat := coord.Lat
		lng := coord.Lng
		if lat < south {
			south = lat
		}
		if lng < west {
			west = lng
		}
		if lat > north {
			north = lat
		}
		if lng > east {
			east = lng
		}
		if lng > 0 && lng < minPosLng {
			minPosLng = lng
		}
		if lng < 0 && lng > maxNegLng {
			maxNegLng = lng
		}
		if math.Abs(lng-next.Lng) > math.Pi {
			isTrans = true
		}
	}
	if isTrans {
		east = maxNegLng
		west = minPosLng
	}
	*bbox = BBox{North: north, South: south, East: east, West: west}
}
