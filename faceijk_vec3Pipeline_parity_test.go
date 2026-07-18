//go:build cgo && c2go && h3v450

package h3

import "testing"

// Direct parity for the seven H3 4.5.0 faceijk Vec3 pipeline functions
// (_vec3ToFaceIjk, _vec3ToHex2d, _vec3ToClosestFace, _faceIjkToVec3,
// _hex2dToVec3, _vec3AzimuthRads, _vec3TangentBasis; file-statics via
// same-TU wrappers). Discrete outputs (faces, IJK coordinates) compare
// exactly; continuous outputs admit a last-ulp difference because the
// pipeline runs through sin/cos/tan/atan2/acos, where Go's math library
// and the platform libm legitimately differ by 1 ulp on some inputs.

func vec3PipelineSamplePoints(t *testing.T) []vec3d {
	t.Helper()
	pts := []vec3d{}
	for _, g := range []LatLng{
		{Lat: Rad(0.5), Lng: Rad(-1.3)},
		{Lat: Rad(0.6), Lng: Rad(-1.2)},
		{Lat: Rad(-0.4), Lng: Rad(0.8)},
		{Lat: Rad(1.5), Lng: Rad(3.1)},
		{Lat: Rad(-1.2), Lng: Rad(-2.7)},
	} {
		pts = append(pts, latLngToVec3(g))
	}
	// Every icosahedron face center exercises the face-selection edges.
	pts = append(pts, faceCenterPoint[:]...)
	return pts
}

func Test_vec3ToFaceIjk_pipeline_parity(t *testing.T) {
	for _, p := range vec3PipelineSamplePoints(t) {
		for _, res := range []int32{0, 1, 5, 9, 15} {
			// _vec3ToClosestFace: face exact, sqd within ulp noise.
			var goFace int32
			var goSqd float64
			_vec3ToClosestFace(&p, &goFace, &goSqd)
			cFace, cSqd := _vec3ToClosestFaceC(p)
			if goFace != cFace || !vec3UlpClose(goSqd, cSqd) {
				t.Fatalf("_vec3ToClosestFace(%v): Go=(%d,%v) C=(%d,%v)", p, goFace, goSqd, cFace, cSqd)
			}

			// _vec3ToHex2d: face exact, hex2d within ulp noise.
			var hFace int32
			var goV vec2d
			_vec3ToHex2d(&p, res, &hFace, &goV)
			cHFace, cV := _vec3ToHex2dC(p, res)
			if hFace != cHFace || !vec3UlpClose(goV.X, cV.X) || !vec3UlpClose(goV.Y, cV.Y) {
				t.Fatalf("_vec3ToHex2d(%v,%d): Go=(%d,%v) C=(%d,%v)", p, res, hFace, goV, cHFace, cV)
			}

			// _vec3ToFaceIjk: fully discrete output, exact.
			var goFijk faceIJK
			_vec3ToFaceIjk(p, res, &goFijk)
			cFijk := _vec3ToFaceIjkC(p, res)
			if goFijk != cFijk {
				t.Fatalf("_vec3ToFaceIjk(%v,%d): Go=%+v C=%+v", p, res, goFijk, cFijk)
			}

			// _faceIjkToVec3 on the discrete address both sides agree on.
			var goBack vec3d
			_faceIjkToVec3(&goFijk, res, &goBack)
			cBack := _faceIjkToVec3C(goFijk, res)
			if !vec3UlpCloseVec(goBack, cBack) {
				t.Fatalf("_faceIjkToVec3(%+v,%d): Go=%v C=%v", goFijk, res, goBack, cBack)
			}
		}

		// _vec3AzimuthRads/_vec3TangentBasis against the point's closest
		// face center — the only pairing the pipeline ever computes.
		// (Arbitrary far pairs are ill-conditioned: for near-antipodal
		// inputs the tangent-plane projection is cancellation-dominated
		// and both implementations return noise.)
		var closest int32
		var sqd float64
		_vec3ToClosestFace(&p, &closest, &sqd)
		center := faceCenterPoint[closest]
		if p != center {
			goAz := _vec3AzimuthRads(center, p)
			cAz := _vec3AzimuthRadsC(center, p)
			if !vec3UlpClose(goAz, cAz) {
				t.Fatalf("_vec3AzimuthRads(face %d, %v): Go=%v C=%v", closest, p, goAz, cAz)
			}
		}
		var goN, goE vec3d
		_vec3TangentBasis(center, &goN, &goE)
		cN, cE := _vec3TangentBasisC(center)
		if !vec3UlpCloseVec(goN, cN) || !vec3UlpCloseVec(goE, cE) {
			t.Fatalf("_vec3TangentBasis(face %d): Go=(%v,%v) C=(%v,%v)", closest, goN, goE, cN, cE)
		}
	}
}

func Test_hex2dToVec3_parity(t *testing.T) {
	points := []vec2d{
		{X: 0, Y: 0},
		{X: 0.25, Y: -0.5},
		{X: 1.75, Y: 0.3},
		{X: -2.5, Y: -1.25},
		{X: 1e-18, Y: 0}, // below EPSILON: face-center fast path
	}
	for _, v := range points {
		for _, face := range []int32{0, 7, 19} {
			for _, res := range []int32{0, 1, 9} {
				for _, substrate := range []int32{0, 1} {
					var goV3 vec3d
					_hex2dToVec3(&v, face, res, substrate, &goV3)
					cV3 := _hex2dToVec3C(v, face, res, substrate)
					if !vec3UlpCloseVec(goV3, cV3) {
						t.Fatalf("_hex2dToVec3(%v,f%d,r%d,s%d): Go=%v C=%v", v, face, res, substrate, goV3, cV3)
					}
				}
			}
		}
	}
}
