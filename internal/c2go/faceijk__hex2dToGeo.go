package c2go

import "math"

// _hex2dToGeo converts 2D hexagonal coordinates to geographic coordinates.
// Mirrors _hex2dToGeo in faceijk.c
// Ported from H3 C: faceijk.c::hex2dToGeo
func _hex2dToGeo(v *Vec2d, face int, res int, substrate int, g *LatLng) {
	// calculate (r, theta) in hex2d
	r := _v2dMag(v)

	if r < EPSILON {
		*g = faceCenterGeo[face]
		return
	}

	theta := math.Atan2(v.Y, v.X)

	// scale for current resolution length u
	for i := 0; i < res; i++ {
		r *= M_RSQRT7
	}

	// scale accordingly if this is a substrate grid
	if substrate != 0 {
		r *= M_ONETHIRD
		if isResolutionClassIII(res) {
			r *= M_RSQRT7
		}
	}

	r *= RES0_U_GNOMONIC

	// perform inverse gnomonic scaling of r
	r = math.Atan(r)

	// adjust theta for Class III
	// if a substrate grid, then it's already been adjusted for Class III
	if substrate == 0 && isResolutionClassIII(res) {
		theta = _posAngleRads(theta + M_AP7_ROT_RADS)
	}

	// find theta as an azimuth
	theta = _posAngleRads(faceAxesAzRadsCII[face][0] - theta)

	// now find the point at (r,theta) from the face center
	_geoAzDistanceRads(&faceCenterGeo[face], theta, r, g)
}
