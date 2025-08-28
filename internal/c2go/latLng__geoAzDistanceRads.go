package c2go

import (
	"math"

	"github.com/dimchansky/h3-go/angle"
)

// _geoAzDistanceRads computes the point p2 at azimuth az and distance from p1.
// Ported from H3 C: latLng.c::_geoAzDistanceRads
func _geoAzDistanceRads(p1 *LatLng, az, distance float64, p2 *LatLng) {
	if distance < EPSILON {
		*p2 = *p1
		return
	}
	az = _posAngleRads(az)
	if az < EPSILON || math.Abs(az-math.Pi) < EPSILON {
		// Due north or south
		if az < EPSILON {
			p2.Lat = p1.Lat + angle.Rad(distance)
		} else {
			p2.Lat = p1.Lat - angle.Rad(distance)
		}
		if math.Abs((p2.Lat - angle.PiOver2).Rad()) < EPSILON {
			p2.Lat = angle.PiOver2
			p2.Lng = 0
		} else if math.Abs((p2.Lat + angle.PiOver2).Rad()) < EPSILON {
			p2.Lat = -angle.PiOver2
			p2.Lng = 0
		} else {
			// Keep longitude the same; ensure wrapped to [-π, π]
			p2.Lng = constrainLng(p1.Lng)
		}
		return
	}

	// Precompute sin/cos where used multiple times
	sd, cd := math.Sincos(distance)
	sa, ca := math.Sincos(az)
	sl1, cl1 := p1.Lat.SinCos()

	sinlat := sl1*cd + cl1*sd*ca
	if sinlat > 1.0 {
		sinlat = 1.0
	}
	if sinlat < -1.0 {
		sinlat = -1.0
	}
	p2.Lat = angle.Rad(math.Asin(sinlat))
	if math.Abs((p2.Lat - angle.PiOver2).Rad()) < EPSILON {
		p2.Lat = angle.PiOver2
		p2.Lng = 0
		return
	} else if math.Abs((p2.Lat + angle.PiOver2).Rad()) < EPSILON {
		p2.Lat = -angle.PiOver2
		p2.Lng = 0
		return
	}

	invcosp2lat := 1.0 / p2.Lat.Cos()
	sinlng := sa * sd * invcosp2lat
	coslng := (cd - sl1*p2.Lat.Sin()) / cl1 * invcosp2lat
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
	p2.Lng = constrainLng(p1.Lng + angle.Rad(math.Atan2(sinlng, coslng)))
}
