//go:build cgo

package h3

import "testing"

func Test_geoToClosestFace_ParityWithC(t *testing.T) {
	testCases := []LatLng{
		{0.0, 0.0},                                 // equator/prime meridian
		{0.8035826497189899, 1.2483974196173961},   // face 0 center (from faceCenterGeo)
		{1.3077478834556382, 2.5369450098779212},   // face 1 center
		{1.0547512535239521, -1.3475173589003966},  // face 2 center
		{0.6001915955381868, -0.4506039094697557},  // face 3 center
		{0.4917154281987739, 0.4019882029113069},   // face 4 center
		{-0.8035826497189899, -1.2483974196173961}, // antipodal point
		{1.5707963267948966, 0.0},                  // north pole
		{-1.5707963267948966, 0.0},                 // south pole
		{0.0, 3.141592653589793},                   // equator/antimeridian
		{0.5, 1.0},                                 // arbitrary mid-latitude
		{-0.5, -1.0},                               // arbitrary southern hemisphere
	}

	for i, tc := range testCases {
		var goFace, cFace int32
		var goSqd, cSqd float64

		// Test Go implementation
		_geoToClosestFace(&tc, &goFace, &goSqd)

		// Test C implementation
		_geoToClosestFaceC(&tc, &cFace, &cSqd)

		if goFace != cFace {
			t.Fatalf("case %d (%g, %g): face mismatch: go=%d c=%d", i, tc.Lat, tc.Lng, goFace, cFace)
		}

		// Allow small floating-point differences in squared distance
		const tolerance = 1e-12
		if absf(goSqd-cSqd) > tolerance {
			t.Fatalf("case %d (%g, %g): sqd mismatch: go=%.15f c=%.15f diff=%.15f", i, tc.Lat, tc.Lng, goSqd, cSqd, absf(goSqd-cSqd))
		}
	}
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
