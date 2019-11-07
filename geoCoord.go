package h3

//lint:file-ignore U1000 Ignore all unused code
// Functions for working with lat/lon coordinates.

import (
	"math"
)

const (
	EPSILON_DEG = .000000001             // epsilon of ~0.1mm in degrees
	EPSILON_RAD = EPSILON_DEG * M_PI_180 // epsilon of ~0.1mm in radians
)

// _posAngleRads normalizes radians to a value between 0.0 and two PI.
// `rads`: The input radians value.
// Returns the normalized radians value.
func _posAngleRads(rads float64) float64 {
	var tmp float64
	if rads < 0.0 {
		tmp = rads + M_2PI
	} else {
		tmp = rads
	}
	if rads >= M_2PI {
		tmp -= M_2PI
	}
	return tmp
}

// geoAlmostEqualThreshold determines if the components of two spherical coordinates are within some
// threshold distance of each other.
// `p1`: The first spherical coordinates.
// `p2`: The second spherical coordinates.
// `threshold`: The threshold distance.
// Returns whether or not the two coordinates are within the threshold distance
// of each other.
func geoAlmostEqualThreshold(p1 *GeoCoord, p2 *GeoCoord, threshold float64) bool {
	return math.Abs(p1.lat-p2.lat) < threshold &&
		math.Abs(p1.lon-p2.lon) < threshold
}

// geoAlmostEqual determines if the components of two spherical coordinates are within our
// standard epsilon distance of each other.
// `p1`: The first spherical coordinates.
// `p2`: The second spherical coordinates.
// Returns whether or not the two coordinates are within the epsilon distance
// of each other.
func geoAlmostEqual(p1 *GeoCoord, p2 *GeoCoord) bool {
	return geoAlmostEqualThreshold(p1, p2, EPSILON_RAD)
}

// setGeoDegs sets the components of spherical coordinates in decimal degrees.
// `p`: The spherical coordinates.
// `latDegs`: The desired latitude in decimal degrees.
// `lonDegs`: The desired longitude in decimal degrees.
func setGeoDegs(p *GeoCoord, latDegs float64, lonDegs float64) {
	_setGeoRads(p, degsToRads(latDegs), degsToRads(lonDegs))
}

// _setGeoRads sets the components of spherical coordinates in radians.
// `p`: The spherical coordinates.
// `latRads`: The desired latitude in decimal radians.
// `lonRads`: The desired longitude in decimal radians.
func _setGeoRads(p *GeoCoord, latRads float64, lonRads float64) {
	p.lat = latRads
	p.lon = lonRads
}

// degsToRads converts from decimal degrees to radians.
// `degrees`: The decimal degrees.
// Returns the corresponding radians.
func degsToRads(degrees float64) float64 { return degrees * M_PI_180 }

// radsToDegs converts from radians to decimal degrees.
// `radians`: The radians.
// Returns the corresponding decimal degrees.
func radsToDegs(radians float64) float64 { return radians * M_180_PI }

// constrainLat makes sure latitudes are in the proper bounds.
// `lat`: The original lat value.
// Returns the corrected lat value
func constrainLat(lat float64) float64 {
	for lat > M_PI_2 {
		lat = lat - M_PI
	}
	return lat
}

// constrainLng makes sure longitudes are in the proper bounds
// `lng`: The origin lng value.
// Returns the corrected lng value.
func constrainLng(lng float64) float64 {
	for lng > M_PI {
		lng = lng - (2 * M_PI)
	}
	for lng < -M_PI {
		lng = lng + (2 * M_PI)
	}
	return lng
}

// _geoDistRads finds the great circle distance in radians between two spherical coordinates.
// `p1`: The first spherical coordinates.
// `p2`: The second spherical coordinates.
// Returns the great circle distance in radians between `p1` and `p2`.
func _geoDistRads(p1 *GeoCoord, p2 *GeoCoord) float64 {
	// use spherical triangle with p1 at A, p2 at B, and north pole at C
	bigC := math.Abs(p2.lon - p1.lon)
	if bigC > M_PI { // assume we want the complement
		// note that in this case they can't both be negative
		lon1 := p1.lon
		if lon1 < 0.0 {
			lon1 += 2.0 * M_PI
		}
		lon2 := p2.lon
		if lon2 < 0.0 {
			lon2 += 2.0 * M_PI
		}

		bigC = math.Abs(lon2 - lon1)
	}

	b := M_PI_2 - p1.lat
	a := M_PI_2 - p2.lat

	// use law of cosines to find c
	cosc := math.Cos(a)*math.Cos(b) + math.Sin(a)*math.Sin(b)*math.Cos(bigC)
	if cosc > 1.0 {
		cosc = 1.0
	}
	if cosc < -1.0 {
		cosc = -1.0
	}

	return math.Acos(cosc)
}

// _geoDistKm finds the great circle distance in kilometers between two spherical
// coordinates.
// `p1`: The first spherical coordinates.
// `p2`: The second spherical coordinates.
// Returns the distance in kilometers between `p1` and `p2`.
func _geoDistKm(p1 *GeoCoord, p2 *GeoCoord) float64 {
	return EARTH_RADIUS_KM * _geoDistRads(p1, p2)
}

// _geoAzimuthRads determines the azimuth to p2 from p1 in radians.
// `p1`: The first spherical coordinates.
// `p2`: The second spherical coordinates.
// Returns the azimuth in radians from `p1` to `p2`.
func _geoAzimuthRads(p1 *GeoCoord, p2 *GeoCoord) float64 {
	return math.Atan2(math.Cos(p2.lat)*math.Sin(p2.lon-p1.lon),
		math.Cos(p1.lat)*math.Sin(p2.lat)-
			math.Sin(p1.lat)*math.Cos(p2.lat)*math.Cos(p2.lon-p1.lon))
}

// _geoAzDistanceRads computes the point on the sphere a specified azimuth and distance from
// another point.
// `p1`: The first spherical coordinates.
// `az`: The desired azimuth from `p1`.
// `distance`: The desired distance from p1, must be non-negative.
// `p2`: The spherical coordinates at the desired azimuth and distance from `p1`.
func _geoAzDistanceRads(p1 *GeoCoord, az float64, distance float64, p2 *GeoCoord) {
	if distance < EPSILON {
		*p2 = *p1
		return
	}

	az = _posAngleRads(az)

	// check for due north/south azimuth
	if az < EPSILON || math.Abs(az-M_PI) < EPSILON {
		if az < EPSILON { // due north
			p2.lat = p1.lat + distance
		} else { // due south
			p2.lat = p1.lat - distance
		}

		if math.Abs(p2.lat-M_PI_2) < EPSILON { // north pole
			p2.lat = M_PI_2
			p2.lon = 0.0
		} else if math.Abs(p2.lat+M_PI_2) < EPSILON { // south pole
			p2.lat = -M_PI_2
			p2.lon = 0.0
		} else {
			p2.lon = constrainLng(p1.lon)
		}
	} else { // not due north or south
		sinlat := math.Sin(p1.lat)*math.Cos(distance) +
			math.Cos(p1.lat)*math.Sin(distance)*math.Cos(az)
		if sinlat > 1.0 {
			sinlat = 1.0
		}
		if sinlat < -1.0 {
			sinlat = -1.0
		}
		p2.lat = math.Asin(sinlat)
		if math.Abs(p2.lat-M_PI_2) < EPSILON { // north pole
			p2.lat = M_PI_2
			p2.lon = 0.0
		} else if math.Abs(p2.lat+M_PI_2) < EPSILON { // south pole
			p2.lat = -M_PI_2
			p2.lon = 0.0
		} else {
			sinlon := math.Sin(az) * math.Sin(distance) / math.Cos(p2.lat)
			coslon := (math.Cos(distance) - math.Sin(p1.lat)*math.Sin(p2.lat)) /
				math.Cos(p1.lat) / math.Cos(p2.lat)
			if sinlon > 1.0 {
				sinlon = 1.0
			}
			if sinlon < -1.0 {
				sinlon = -1.0
			}
			if coslon > 1.0 {
				coslon = 1.0
			}
			if coslon < -1.0 {
				coslon = -1.0
			}
			p2.lon = constrainLng(p1.lon + math.Atan2(sinlon, coslon))
		}
	}
}

// The following functions provide meta information about the H3 hexagons at
// each zoom level. Since there are only 16 total levels, these are current
// handled with hardwired static values, but it may be worthwhile to put these
// static values into another file that can be autogenerated by source code in
// the future.

var (
	areasKm2 = []float64{
		4250546.848, 607220.9782, 86745.85403, 12392.26486,
		1770.323552, 252.9033645, 36.1290521, 5.1612932,
		0.7373276, 0.1053325, 0.0150475, 0.0021496,
		0.0003071, 0.0000439, 0.0000063, 0.0000009,
	}
	areasM2 = []float64{
		4.25055E+12, 6.07221E+11, 86745854035, 12392264862,
		1770323552, 252903364.5, 36129052.1, 5161293.2,
		737327.6, 105332.5, 15047.5, 2149.6,
		307.1, 43.9, 6.3, 0.9}
	lensKm = []float64{
		1107.712591, 418.6760055, 158.2446558, 59.81085794,
		22.6063794, 8.544408276, 3.229482772, 1.220629759,
		0.461354684, 0.174375668, 0.065907807, 0.024910561,
		0.009415526, 0.003559893, 0.001348575, 0.000509713,
	}
	lensM = []float64{
		1107712.591, 418676.0055, 158244.6558, 59810.85794,
		22606.3794, 8544.408276, 3229.482772, 1220.629759,
		461.3546837, 174.3756681, 65.90780749, 24.9105614,
		9.415526211, 3.559893033, 1.348574562, 0.509713273,
	}
	hexaNums = []int64{
		122,
		842,
		5882,
		41162,
		288122,
		2016842,
		14117882,
		98825162,
		691776122,
		4842432842,
		33897029882,
		237279209162,
		1660954464122,
		11626681248842,
		81386768741882,
		569707381193162,
	}
)

func hexAreaKm2(res int32) float64   { return areasKm2[res] }
func hexAreaM2(res int32) float64    { return areasM2[res] }
func edgeLengthKm(res int32) float64 { return lensKm[res] }
func edgeLengthM(res int32) float64  { return lensM[res] }
func numHexagons(res int32) int64    { return hexaNums[res] }
