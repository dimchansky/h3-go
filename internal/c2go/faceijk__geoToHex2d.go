package c2go

import "math"

// _geoToHex2d converts geographic coordinates to 2D hex coordinates on a specified face.
// Mirrors _geoToHex2d in faceijk.c
func _geoToHex2d(g *LatLng, res int, face *int, v *Vec2d) {
	// determine the icosahedron face
	var sqd float64
	_geoToClosestFace(g, face, &sqd)

	// cos(r) = 1 - 2 * sin^2(r/2) = 1 - 2 * (sqd / 4) = 1 - sqd/2
	r := math.Acos(1 - sqd*0.5)

	if r < EPSILON {
		v.X = 0.0
		v.Y = 0.0
		return
	}

	// now have face and r, now find CCW theta from CII i-axis
	theta := _posAngleRads(faceAxesAzRadsCII[*face][0] -
		_posAngleRads(_geoAzimuthRads(&faceCenterGeo[*face], g)))

	// adjust theta for Class III (odd resolutions)
	if isResolutionClassIII(res) != 0 {
		theta = _posAngleRads(theta - M_AP7_ROT_RADS)
	}

	// perform gnomonic scaling of r
	r = math.Tan(r)

	// scale for current resolution length u
	r *= INV_RES0_U_GNOMONIC
	for i := 0; i < res; i++ {
		r *= M_SQRT7
	}

	// we now have (r, theta) in hex2d with theta ccw from x-axes

	// convert to local x,y
	v.X = r * math.Cos(theta)
	v.Y = r * math.Sin(theta)
}
