package h3

import "math"

// _vec3ToHex2d encodes a coordinate on the sphere to the corresponding
// icosahedral face and containing 2D hex coordinates relative to that
// face center.
//
// Vec3d p is expected to be on the unit sphere.
// Ported from H3 C: faceijk.c::_vec3ToHex2d.
func _vec3ToHex2d(p *vec3d, res int32, face *int32, v *vec2d) {
	// determine the icosahedron face
	var sqd float64
	_vec3ToClosestFace(p, face, &sqd)

	// cos(r) = 1 - 2 * sin^2(r/2) = 1 - 2 * (sqd / 4) = 1 - sqd/2
	r := math.Acos(1 - sqd*0.5)

	if r < epsilon {
		v.X = 0.0
		v.Y = 0.0
		return
	}

	// now have face and r, now find CCW theta from CII i-axis
	theta := _posAngleRads(
		faceAxesAzRadsCII[*face][0] -
			_posAngleRads(_vec3AzimuthRads(faceCenterPoint[*face], *p)))

	// adjust theta for Class III (odd resolutions)
	if isResolutionClassIII(res) {
		theta = _posAngleRads(theta - mAp7RotRads)
	}

	// perform gnomonic scaling of r
	r = math.Tan(r)

	// scale for current resolution length u
	r *= invRes0UGnomonic
	for i := int32(0); i < res; i++ {
		r *= mSqrt7
	}

	// we now have (r, theta) in hex2d with theta ccw from x-axes

	// convert to local x,y
	v.X = r * math.Cos(theta)
	v.Y = r * math.Sin(theta)
}
