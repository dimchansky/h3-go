//go:build c2go

package c2go

import "testing"

func Test_faceIjkToGeo_ParityWithC(t *testing.T) {
	testCases := []struct {
		name    string
		faceIjk FaceIJK
		res     int
	}{
		// Face centers (origin coordinates)
		{"face0_origin_res0", FaceIJK{0, CoordIJK{0, 0, 0}}, 0},
		{"face1_origin_res0", FaceIJK{1, CoordIJK{0, 0, 0}}, 0},
		{"face5_origin_res0", FaceIJK{5, CoordIJK{0, 0, 0}}, 0},
		{"face10_origin_res0", FaceIJK{10, CoordIJK{0, 0, 0}}, 0},
		{"face19_origin_res0", FaceIJK{19, CoordIJK{0, 0, 0}}, 0},

		// Unit directions from face centers
		{"face0_unit_i", FaceIJK{0, CoordIJK{1, 0, 0}}, 0},
		{"face0_unit_j", FaceIJK{0, CoordIJK{0, 1, 0}}, 0},
		{"face0_unit_k", FaceIJK{0, CoordIJK{0, 0, 1}}, 0},
		{"face5_unit_i", FaceIJK{5, CoordIJK{1, 0, 0}}, 0},
		{"face5_unit_j", FaceIJK{5, CoordIJK{0, 1, 0}}, 0},
		{"face5_unit_k", FaceIJK{5, CoordIJK{0, 0, 1}}, 0},

		// Combined directions
		{"face0_combined_ij", FaceIJK{0, CoordIJK{1, 1, 0}}, 0},
		{"face0_combined_ik", FaceIJK{0, CoordIJK{1, 0, 1}}, 0},
		{"face0_combined_jk", FaceIJK{0, CoordIJK{0, 1, 1}}, 0},
		{"face0_combined_ijk", FaceIJK{0, CoordIJK{1, 1, 1}}, 0},

		// Negative coordinates
		{"face0_negative_i", FaceIJK{0, CoordIJK{-1, 0, 0}}, 0},
		{"face0_negative_j", FaceIJK{0, CoordIJK{0, -1, 0}}, 0},
		{"face0_negative_k", FaceIJK{0, CoordIJK{0, 0, -1}}, 0},
		{"face0_mixed_signs", FaceIJK{0, CoordIJK{1, -1, 0}}, 0},

		// Larger coordinates
		{"face0_large_coords", FaceIJK{0, CoordIJK{5, 3, 2}}, 0},
		{"face0_large_coords_2", FaceIJK{0, CoordIJK{-3, 7, -1}}, 0},

		// Different resolutions
		{"face0_origin_res1", FaceIJK{0, CoordIJK{0, 0, 0}}, 1},
		{"face0_origin_res5", FaceIJK{0, CoordIJK{0, 0, 0}}, 5},
		{"face0_origin_res10", FaceIJK{0, CoordIJK{0, 0, 0}}, 10},
		{"face0_origin_res15", FaceIJK{0, CoordIJK{0, 0, 0}}, 15},

		// Different faces with various resolutions
		{"face3_unit_res3", FaceIJK{3, CoordIJK{1, 0, 0}}, 3},
		{"face8_unit_res7", FaceIJK{8, CoordIJK{0, 1, 0}}, 7},
		{"face15_combined_res9", FaceIJK{15, CoordIJK{2, 1, 0}}, 9},

		// Class III resolutions (odd)
		{"face0_origin_res3", FaceIJK{0, CoordIJK{0, 0, 0}}, 3},
		{"face0_unit_res3", FaceIJK{0, CoordIJK{1, 0, 0}}, 3},
		{"face5_unit_res7", FaceIJK{5, CoordIJK{0, 1, 0}}, 7},
		{"face10_combined_res11", FaceIJK{10, CoordIJK{1, 1, 0}}, 11},

		// Edge cases
		{"face0_asymmetric", FaceIJK{0, CoordIJK{10, 5, 2}}, 5},
		{"face19_asymmetric", FaceIJK{19, CoordIJK{-2, 8, -3}}, 8},

		// Various faces for coverage
		{"face7_medium", FaceIJK{7, CoordIJK{3, 2, 1}}, 4},
		{"face12_medium", FaceIJK{12, CoordIJK{-1, 4, 2}}, 6},
		{"face16_medium", FaceIJK{16, CoordIJK{2, -3, 1}}, 12},
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
			if absf(goGeo.Lat-cGeo.Lat) > tolerance {
				t.Fatalf("Lat mismatch: go=%.15f c=%.15f diff=%.15f", goGeo.Lat, cGeo.Lat, absf(goGeo.Lat-cGeo.Lat))
			}

			if absf(goGeo.Lng-cGeo.Lng) > tolerance {
				t.Fatalf("Lng mismatch: go=%.15f c=%.15f diff=%.15f", goGeo.Lng, cGeo.Lng, absf(goGeo.Lng-cGeo.Lng))
			}
		})
	}
}

// Test that face centers match the known values
func Test_faceIjkToGeo_FaceCenters(t *testing.T) {
	for face := 0; face < NUM_ICOSA_FACES; face++ {
		t.Run("face_"+string(rune('0'+face)), func(t *testing.T) {
			var result LatLng
			faceijk := FaceIJK{face, CoordIJK{0, 0, 0}}

			_faceIjkToGeo(&faceijk, 0, &result)

			// Compare with known face center values
			expected := faceCenterGeo[face]
			const tolerance = 1e-12

			if absf(result.Lat-expected.Lat) > tolerance {
				t.Fatalf("Face %d center Lat mismatch: got=%.15f expected=%.15f diff=%.15f",
					face, result.Lat, expected.Lat, absf(result.Lat-expected.Lat))
			}

			if absf(result.Lng-expected.Lng) > tolerance {
				t.Fatalf("Face %d center Lng mismatch: got=%.15f expected=%.15f diff=%.15f",
					face, result.Lng, expected.Lng, absf(result.Lng-expected.Lng))
			}
		})
	}
}
