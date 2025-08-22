package c2go

import (
	"math"
)

// _geoAzDistanceRads computes the point p2 at azimuth az and distance from p1.
// Ported from H3 C: latLng.c::_geoAzDistanceRads
func _geoAzDistanceRads(p1 *LatLng, az, distance float64) LatLng {
	const epsilon = 1e-16
	if distance < epsilon {
		return *p1
	}
	az = _posAngleRads(az)

	var p2 LatLng
	if az < epsilon || math.Abs(az-math.Pi) < epsilon {
		// Due north or south
		if az < epsilon {
			p2.Lat = p1.Lat + distance
		} else {
			p2.Lat = p1.Lat - distance
		}
		if math.Abs(p2.Lat-math.Pi/2) < epsilon {
			p2.Lat = math.Pi / 2
			p2.Lng = 0
		} else if math.Abs(p2.Lat+math.Pi/2) < epsilon {
			p2.Lat = -math.Pi / 2
			p2.Lng = 0
		} else {
			p2.Lng = constrainLng(p1.Lng)
		}
		return p2
	}

	sinlat := math.Sin(p1.Lat)*math.Cos(distance) + math.Cos(p1.Lat)*math.Sin(distance)*math.Cos(az)
	if sinlat > 1.0 {
		sinlat = 1.0
	}
	if sinlat < -1.0 {
		sinlat = -1.0
	}
	p2.Lat = math.Asin(sinlat)
	if math.Abs(p2.Lat-math.Pi/2) < epsilon {
		p2.Lat = math.Pi / 2
		p2.Lng = 0
		return p2
	} else if math.Abs(p2.Lat+math.Pi/2) < epsilon {
		p2.Lat = -math.Pi / 2
		p2.Lng = 0
		return p2
	}

	invcosp2lat := 1.0 / math.Cos(p2.Lat)
	sinlng := math.Sin(az) * math.Sin(distance) * invcosp2lat
	coslng := (math.Cos(distance) - math.Sin(p1.Lat)*math.Sin(p2.Lat)) / math.Cos(p1.Lat) * invcosp2lat
	if sinlng > 1.0 {
		sinlng = 1.0
	}
	if sinlng < -1.0 {
		sinlng = -1.0
	}
	if coslng > 1.0 {
		coslng = 1.0
	}
	if coslng < -1.0 {
		coslng = -1.0
	}
	p2.Lng = constrainLng(p1.Lng + math.Atan2(sinlng, coslng))
	return p2
}
