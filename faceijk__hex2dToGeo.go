package h3

import "math"

// _hex2dToGeo converts 2D hexagonal coordinates to geographic coordinates.
// Mirrors _hex2dToGeo in faceijk.c
// Ported from H3 C: faceijk.c::_hex2dToGeo.
func _hex2dToGeo(v *vec2d, face int32, res int32, substrate int32, g *LatLng) {
	// calculate (r, theta) in hex2d
	r := _v2dMag(v)

	if r < epsilon {
		*g = faceCenterGeo[face]
		return
	}

	theta := math.Atan2(v.Y, v.X)

	// scale for current resolution length u
	for i := int32(0); i < res; i++ {
		r *= mRsqrt7
	}

	// scale accordingly if this is a substrate grid
	if substrate != 0 {
		r *= mOnethird
		if isResolutionClassIII(res) {
			r *= mRsqrt7
		}
	}

	r *= res0UGnomonic

	// perform inverse gnomonic scaling of r
	r = math.Atan(r)

	// adjust theta for Class III
	// if a substrate grid, then it's already been adjusted for Class III
	if substrate == 0 && isResolutionClassIII(res) {
		theta = _posAngleRads(theta + mAp7RotRads)
	}

	// find theta as an azimuth
	theta = _posAngleRads(faceAxesAzRadsCII[face][0] - theta)

	// now find the point at (r,theta) from the face center
	_geoAzDistanceRads(&faceCenterGeo[face], theta, r, g)
}
