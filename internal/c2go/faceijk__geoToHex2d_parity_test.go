//go:build cgo

package c2go

import "testing"

func Test_geoToHex2d_ParityWithC(t *testing.T) {
	testCases := []struct {
		name string
		geo  LatLng
		res  int
	}{
		{"equator_prime", LatLng{0.0, 0.0}, 0},
		{"equator_prime_res5", LatLng{0.0, 0.0}, 5},
		{"equator_prime_res10", LatLng{0.0, 0.0}, 10},
		{"face_0_center_res0", LatLng{0.8035826497189899, 1.2483974196173961}, 0},
		{"face_0_center_res5", LatLng{0.8035826497189899, 1.2483974196173961}, 5},
		{"face_1_center_res0", LatLng{1.3077478834556382, 2.5369450098779212}, 0},
		{"face_2_center_res0", LatLng{1.0547512535239521, -1.3475173589003966}, 0},
		{"north_pole_res0", LatLng{1.5707963267948966, 0.0}, 0},
		{"north_pole_res3", LatLng{1.5707963267948966, 0.0}, 3},
		{"south_pole_res0", LatLng{-1.5707963267948966, 0.0}, 0},
		{"antimeridian_res0", LatLng{0.0, 3.141592653589793}, 0},
		{"mid_latitude_res0", LatLng{0.5, 1.0}, 0},
		{"mid_latitude_res7", LatLng{0.5, 1.0}, 7},
		{"southern_hem_res0", LatLng{-0.5, -1.0}, 0},
		{"class_III_res1", LatLng{0.5, 1.0}, 1},
		{"class_III_res3", LatLng{0.5, 1.0}, 3},
		{"class_III_res9", LatLng{0.5, 1.0}, 9},
		{"high_res", LatLng{0.1, 0.1}, 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var goFace, cFace int
			var goVec, cVec Vec2d

			// Test Go implementation
			_geoToHex2d(&tc.geo, tc.res, &goFace, &goVec)

			// Test C implementation
			_geoToHex2dC(&tc.geo, tc.res, &cFace, &cVec)

			if goFace != cFace {
				t.Fatalf("face mismatch: go=%d c=%d", goFace, cFace)
			}

			// Allow small floating-point differences in Vec2d coordinates
			// Use relative tolerance for large coordinates (high resolution cases)
			const baseTolerance = 1e-12

			// Calculate relative tolerance based on coordinate magnitude
			tolerance := baseTolerance
			maxCoord := absf(goVec.X)
			if absf(goVec.Y) > maxCoord {
				maxCoord = absf(goVec.Y)
			}
			if absf(cVec.X) > maxCoord {
				maxCoord = absf(cVec.X)
			}
			if absf(cVec.Y) > maxCoord {
				maxCoord = absf(cVec.Y)
			}

			// For large coordinates, use relative tolerance
			if maxCoord > 1000 {
				tolerance = maxCoord * 1e-12
			}

			if absf(goVec.X-cVec.X) > tolerance {
				t.Fatalf("X coordinate mismatch: go=%.15f c=%.15f diff=%.15f (tolerance=%.15f)", goVec.X, cVec.X, absf(goVec.X-cVec.X), tolerance)
			}

			if absf(goVec.Y-cVec.Y) > tolerance {
				t.Fatalf("Y coordinate mismatch: go=%.15f c=%.15f diff=%.15f (tolerance=%.15f)", goVec.Y, cVec.Y, absf(goVec.Y-cVec.Y), tolerance)
			}
		})
	}
}
