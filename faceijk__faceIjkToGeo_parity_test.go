//go:build cgo && c2go && !h3v450

package h3

import "testing"

func Test_faceIjkToGeo_ParityWithC(t *testing.T) {
	testCases := []struct {
		name    string
		faceIjk faceIJK
		res     int32
	}{
		// Face centers (origin coordinates)
		{"face0_origin_res0", faceIJK{0, coordIJK{0, 0, 0}}, 0},
		{"face1_origin_res0", faceIJK{1, coordIJK{0, 0, 0}}, 0},
		{"face5_origin_res0", faceIJK{5, coordIJK{0, 0, 0}}, 0},
		{"face10_origin_res0", faceIJK{10, coordIJK{0, 0, 0}}, 0},
		{"face19_origin_res0", faceIJK{19, coordIJK{0, 0, 0}}, 0},

		// Unit directions from face centers
		{"face0_unit_i", faceIJK{0, coordIJK{1, 0, 0}}, 0},
		{"face0_unit_j", faceIJK{0, coordIJK{0, 1, 0}}, 0},
		{"face0_unit_k", faceIJK{0, coordIJK{0, 0, 1}}, 0},
		{"face5_unit_i", faceIJK{5, coordIJK{1, 0, 0}}, 0},
		{"face5_unit_j", faceIJK{5, coordIJK{0, 1, 0}}, 0},
		{"face5_unit_k", faceIJK{5, coordIJK{0, 0, 1}}, 0},

		// Combined directions
		{"face0_combined_ij", faceIJK{0, coordIJK{1, 1, 0}}, 0},
		{"face0_combined_ik", faceIJK{0, coordIJK{1, 0, 1}}, 0},
		{"face0_combined_jk", faceIJK{0, coordIJK{0, 1, 1}}, 0},
		{"face0_combined_ijk", faceIJK{0, coordIJK{1, 1, 1}}, 0},

		// Negative coordinates
		{"face0_negative_i", faceIJK{0, coordIJK{-1, 0, 0}}, 0},
		{"face0_negative_j", faceIJK{0, coordIJK{0, -1, 0}}, 0},
		{"face0_negative_k", faceIJK{0, coordIJK{0, 0, -1}}, 0},
		{"face0_mixed_signs", faceIJK{0, coordIJK{1, -1, 0}}, 0},

		// Larger coordinates
		{"face0_large_coords", faceIJK{0, coordIJK{5, 3, 2}}, 0},
		{"face0_large_coords_2", faceIJK{0, coordIJK{-3, 7, -1}}, 0},

		// Different resolutions
		{"face0_origin_res1", faceIJK{0, coordIJK{0, 0, 0}}, 1},
		{"face0_origin_res5", faceIJK{0, coordIJK{0, 0, 0}}, 5},
		{"face0_origin_res10", faceIJK{0, coordIJK{0, 0, 0}}, 10},
		{"face0_origin_res15", faceIJK{0, coordIJK{0, 0, 0}}, 15},

		// Different faces with various resolutions
		{"face3_unit_res3", faceIJK{3, coordIJK{1, 0, 0}}, 3},
		{"face8_unit_res7", faceIJK{8, coordIJK{0, 1, 0}}, 7},
		{"face15_combined_res9", faceIJK{15, coordIJK{2, 1, 0}}, 9},

		// Class III resolutions (odd)
		{"face0_origin_res3", faceIJK{0, coordIJK{0, 0, 0}}, 3},
		{"face0_unit_res3", faceIJK{0, coordIJK{1, 0, 0}}, 3},
		{"face5_unit_res7", faceIJK{5, coordIJK{0, 1, 0}}, 7},
		{"face10_combined_res11", faceIJK{10, coordIJK{1, 1, 0}}, 11},

		// Edge cases
		{"face0_asymmetric", faceIJK{0, coordIJK{10, 5, 2}}, 5},
		{"face19_asymmetric", faceIJK{19, coordIJK{-2, 8, -3}}, 8},

		// Various faces for coverage
		{"face7_medium", faceIJK{7, coordIJK{3, 2, 1}}, 4},
		{"face12_medium", faceIJK{12, coordIJK{-1, 4, 2}}, 6},
		{"face16_medium", faceIJK{16, coordIJK{2, -3, 1}}, 12},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var goGeo, cGeo LatLng

			// Test Go implementation
			_faceIjkToGeo(&tc.faceIjk, tc.res, &goGeo)

			// Test C implementation
			_faceIjkToGeoC(&tc.faceIjk, tc.res, &cGeo)

			// Allow small floating-point differences in geographic coordinates
			const tolerance = 1e-12
			if !goGeo.Lat.EqualApprox(cGeo.Lat, tolerance) {
				t.Fatalf("Lat mismatch: go=%.15f c=%.15f diff=%.15f", goGeo.Lat.Rad(), cGeo.Lat.Rad(), absf(goGeo.Lat.Rad()-cGeo.Lat.Rad()))
			}

			if !goGeo.Lng.EqualApprox(cGeo.Lng, tolerance) {
				t.Fatalf("Lng mismatch: go=%.15f c=%.15f diff=%.15f", goGeo.Lng.Rad(), cGeo.Lng.Rad(), absf(goGeo.Lng.Rad()-cGeo.Lng.Rad()))
			}
		})
	}
}

// Test that face centers match the known values
func Test_faceIjkToGeo_FaceCenters(t *testing.T) {
	for face := int32(0); face < numIcosaFaces; face++ {
		t.Run("face_"+string(rune('0'+face)), func(t *testing.T) {
			var result LatLng
			faceijk := faceIJK{face, coordIJK{0, 0, 0}}

			_faceIjkToGeo(&faceijk, 0, &result)

			// Compare with known face center values
			expected := faceCenterGeo[face]
			const tolerance = 1e-12

			if !result.Lat.EqualApprox(expected.Lat, tolerance) {
				t.Fatalf("Face %d center Lat mismatch: got=%.15f expected=%.15f diff=%.15f",
					face, result.Lat.Rad(), expected.Lat.Rad(), absf(result.Lat.Rad()-expected.Lat.Rad()))
			}

			if !result.Lng.EqualApprox(expected.Lng, tolerance) {
				t.Fatalf("Face %d center Lng mismatch: got=%.15f expected=%.15f diff=%.15f",
					face, result.Lng.Rad(), expected.Lng.Rad(), absf(result.Lng.Rad()-expected.Lng.Rad()))
			}
		})
	}
}
