package c2go

import (
	"math"

	"github.com/dimchansky/h3-go/angle"
)

// bboxFromGeoLoop computes a bounding box for a loop of coordinates.
// Mirrors H3's polygon.c::bboxFromGeoLoop behavior.
// Ported from H3 C: polygon.c::bboxFromGeoLoop
func bboxFromGeoLoop(loop []LatLng, bbox *BBox) {
	if len(loop) == 0 {
		*bbox = BBox{}
		return
	}
	south := angle.Rad(math.MaxFloat64)
	west := angle.Rad(math.MaxFloat64)
	north := angle.Rad(-math.MaxFloat64)
	east := angle.Rad(-math.MaxFloat64)
	minPosLng := angle.Rad(math.MaxFloat64)
	maxNegLng := angle.Rad(-math.MaxFloat64)
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
		if (lng - next.Lng).Abs() > angle.Pi {
			isTrans = true
		}
	}
	if isTrans {
		east = maxNegLng
		west = minPosLng
	}
	*bbox = BBox{North: north, South: south, East: east, West: west}
}
