package h3

import "math"

// _hex2dToVec3 determines the 3D coordinates of a cell given by 2D
// hex coordinates on a particular icosahedral face.
// Ported from H3 C: faceijk.c::_hex2dToVec3.
func _hex2dToVec3(v *vec2d, face int32, res int32, substrate int32, v3 *vec3d) {
	// calculate (r, theta) in hex2d
	r := _v2dMag(v)

	if r < epsilon {
		*v3 = faceCenterPoint[face]
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
	var northDir, eastDir vec3d
	_vec3TangentBasis(faceCenterPoint[face], &northDir, &eastDir)

	dir := vec3LinComb(math.Cos(theta), northDir, math.Sin(theta), eastDir)

	*v3 = vec3LinComb(math.Cos(r), faceCenterPoint[face], math.Sin(r), dir)
	vec3Normalize(v3)
}
