//go:build cgo

package c2go

import "testing"

func Test_h3index_latLngToCell_ParityWithC(t *testing.T) {
	testCases := []struct {
		lat, lng float64
		res      int32
	}{
		// Basic test cases
		{37.775938728915946, -122.41795063018799, 9},
		{40.689167, -74.044444, 10},
		{0.0, 0.0, 0},
		{0.0, 0.0, 15},

		// Edge cases
		{90.0, 0.0, 5},   // North pole
		{-90.0, 0.0, 5},  // South pole
		{0.0, 180.0, 5},  // International date line
		{0.0, -180.0, 5}, // International date line

		// Pentagon locations (approximate)
		{58.2, 11.25, 0},    // Base cell 4
		{-58.2, -168.75, 0}, // Base cell 49

		// Various resolutions
		{37.775938, -122.417951, 0},
		{37.775938, -122.417951, 1},
		{37.775938, -122.417951, 5},
		{37.775938, -122.417951, 10},
		{37.775938, -122.417951, 15},

		// Random locations
		{51.5074, -0.1278, 7},   // London
		{35.6762, 139.6503, 8},  // Tokyo
		{-33.8688, 151.2093, 6}, // Sydney
		{-22.9068, -43.1729, 9}, // Rio de Janeiro
	}

	for _, tc := range testCases {
		g := &LatLng{Lat: tc.lat, Lng: tc.lng}

		var goOut H3Index
		var cOut H3Index

		goErr := latLngToCell(g, tc.res, &goOut)
		cErr := latLngToCellC(g, tc.res, &cOut)

		if uint32(goErr) != cErr {
			t.Fatalf("latLngToCell error mismatch for lat=%f lng=%f res=%d: go=%d c=%d",
				tc.lat, tc.lng, tc.res, uint32(goErr), cErr)
		}

		if goOut != cOut {
			t.Fatalf("latLngToCell result mismatch for lat=%f lng=%f res=%d: go=%x c=%x",
				tc.lat, tc.lng, tc.res, uint64(goOut), uint64(cOut))
		}
	}
}

func Test_h3index_latLngToCell_ErrorCases_ParityWithC(t *testing.T) {
	errorCases := []struct {
		lat, lng float64
		res      int32
		desc     string
	}{
		// Invalid resolution
		{37.775938, -122.417951, -1, "negative resolution"},
		{37.775938, -122.417951, 16, "resolution too high"},
	}

	for _, tc := range errorCases {
		g := &LatLng{Lat: tc.lat, Lng: tc.lng}

		var goOut H3Index
		var cOut H3Index

		goErr := latLngToCell(g, tc.res, &goOut)
		cErr := latLngToCellC(g, tc.res, &cOut)

		if uint32(goErr) != cErr {
			t.Fatalf("latLngToCell error mismatch for %s (lat=%f lng=%f res=%d): go=%d c=%d",
				tc.desc, tc.lat, tc.lng, tc.res, uint32(goErr), cErr)
		}

		// Both should return error
		if goErr == E_SUCCESS {
			t.Fatalf("latLngToCell should return error for %s", tc.desc)
		}
	}
}
