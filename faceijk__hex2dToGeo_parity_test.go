//go:build cgo && c2go

package h3

import "testing"

func Test_hex2dToGeo_ParityWithC(t *testing.T) {
	testCases := []struct {
		name      string
		vec       Vec2d
		face      int32
		res       int32
		substrate int32
	}{
		{"origin_face0_res0", Vec2d{0.0, 0.0}, 0, 0, 0},
		{"origin_face5_res0", Vec2d{0.0, 0.0}, 5, 0, 0},
		{"origin_face10_res5", Vec2d{0.0, 0.0}, 10, 5, 0},
		{"small_offset_face0", Vec2d{0.1, 0.05}, 0, 0, 0},
		{"small_offset_face0_res3", Vec2d{0.1, 0.05}, 0, 3, 0},
		{"unit_x_face0", Vec2d{1.0, 0.0}, 0, 0, 0},
		{"unit_y_face0", Vec2d{0.0, 1.0}, 0, 0, 0},
		{"diagonal_face0", Vec2d{0.7071, 0.7071}, 0, 0, 0},
		{"negative_x_face0", Vec2d{-1.0, 0.0}, 0, 0, 0},
		{"negative_y_face0", Vec2d{0.0, -1.0}, 0, 0, 0},
		{"various_faces_res0", Vec2d{0.5, 0.3}, 1, 0, 0},
		{"various_faces_res0_2", Vec2d{0.5, 0.3}, 8, 0, 0},
		{"various_faces_res0_3", Vec2d{0.5, 0.3}, 15, 0, 0},
		{"class_III_res1", Vec2d{0.2, 0.1}, 0, 1, 0},
		{"class_III_res3", Vec2d{0.2, 0.1}, 0, 3, 0},
		{"class_III_res9", Vec2d{0.2, 0.1}, 0, 9, 0},
		{"high_res", Vec2d{0.001, 0.002}, 0, 15, 0},
		{"substrate_test", Vec2d{0.5, 0.3}, 0, 2, 1},
		{"substrate_class_III", Vec2d{0.5, 0.3}, 0, 3, 1},
		{"large_coords", Vec2d{10.0, 5.0}, 0, 0, 0},
		{"medium_coords_res7", Vec2d{2.5, 1.8}, 5, 7, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var goGeo, cGeo LatLng

			// Test Go implementation
			_hex2dToGeo(&tc.vec, tc.face, tc.res, tc.substrate, &goGeo)

			// Test C implementation
			_hex2dToGeoC(&tc.vec, tc.face, tc.res, tc.substrate, &cGeo)

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

// Test round-trip consistency: geo -> hex2d -> geo
func Test_hex2dToGeo_RoundTripConsistency(t *testing.T) {
	testCases := []struct {
		name    string
		geo     LatLng
		skipLng bool // Skip longitude comparison for special cases like poles
	}{
		{"equator_prime", LatLng{0.0, 0.0}, false},
		{"face_center", LatLng{0.8035826497189899, 1.2483974196173961}, false},
		{"north_pole", LatLng{1.5707963267948966, 0.0}, true},  // longitude undefined at poles
		{"south_pole", LatLng{-1.5707963267948966, 0.0}, true}, // longitude undefined at poles
		{"arbitrary_point", LatLng{0.5, 1.0}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// geo -> hex2d -> geo
			var face int32
			var hex2d Vec2d
			var result LatLng

			_geoToHex2d(&tc.geo, 5, &face, &hex2d) // Use resolution 5
			_hex2dToGeo(&hex2d, face, 5, 0, &result)

			const tolerance = 1e-10 // Slightly relaxed for round-trip
			if !tc.geo.Lat.EqualApprox(result.Lat, tolerance) {
				t.Fatalf("Round-trip Lat mismatch: orig=%.15f result=%.15f diff=%.15f", tc.geo.Lat.Rad(), result.Lat.Rad(), absf(tc.geo.Lat.Rad()-result.Lat.Rad()))
			}

			if !tc.skipLng {
				if !tc.geo.Lng.EqualApprox(result.Lng, tolerance) {
					t.Fatalf("Round-trip Lng mismatch: orig=%.15f result=%.15f diff=%.15f", tc.geo.Lng.Rad(), result.Lng.Rad(), absf(tc.geo.Lng.Rad()-result.Lng.Rad()))
				}
			}
		})
	}
}
