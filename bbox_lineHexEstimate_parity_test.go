//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_lineHexEstimate_parity(t *testing.T) {
	tests := []struct {
		name        string
		origin      LatLng
		destination LatLng
		res         int32
	}{
		{
			name:        "same_point",
			origin:      LatLng{Lat: Deg(37.775), Lng: Deg(-122.418)},
			destination: LatLng{Lat: Deg(37.775), Lng: Deg(-122.418)},
			res:         5,
		},
		{
			name:        "short_distance",
			origin:      LatLng{Lat: Deg(37.775), Lng: Deg(-122.418)},
			destination: LatLng{Lat: Deg(37.785), Lng: Deg(-122.428)},
			res:         5,
		},
		{
			name:        "medium_distance",
			origin:      LatLng{Lat: Deg(37.775), Lng: Deg(-122.418)},
			destination: LatLng{Lat: Deg(40.689), Lng: Deg(-74.044)}, // SF to NYC
			res:         3,
		},
		{
			name:        "long_distance",
			origin:      LatLng{Lat: Deg(37.775), Lng: Deg(-122.418)}, // SF
			destination: LatLng{Lat: Deg(51.507), Lng: Deg(-0.127)},   // London
			res:         2,
		},
		{
			name:        "cross_antimeridian",
			origin:      LatLng{Lat: Deg(37.775), Lng: Deg(179.0)},
			destination: LatLng{Lat: Deg(37.775), Lng: Deg(-179.0)},
			res:         4,
		},
		{
			name:        "high_res",
			origin:      LatLng{Lat: Deg(37.775), Lng: Deg(-122.418)},
			destination: LatLng{Lat: Deg(37.776), Lng: Deg(-122.419)},
			res:         10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.res < 0 || tt.res > MAX_H3_RES {
				t.Skipf("Invalid resolution: %d", tt.res)
			}

			var goOut, cOut int64
			goErr := lineHexEstimate(&tt.origin, &tt.destination, tt.res, &goOut)
			cErr := lineHexEstimateC(&tt.origin, &tt.destination, tt.res, &cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("lineHexEstimate error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If there was an error, both should agree on the error
			if goErr != E_SUCCESS {
				return
			}

			// Check outputs match exactly for successful cases
			if goOut != cOut {
				t.Errorf("lineHexEstimate(%+v -> %+v, res=%d): Go=%d, C=%d",
					tt.origin, tt.destination, tt.res, goOut, cOut)
			}

			// Sanity check - estimate should be positive
			if goOut <= 0 {
				t.Errorf("lineHexEstimate should return positive estimate, got %d", goOut)
			}
		})
	}
}

func Test_lineHexEstimate_edge_cases_parity(t *testing.T) {
	tests := []struct {
		name        string
		origin      LatLng
		destination LatLng
		res         int32
		expectErr   H3Error
	}{
		{
			name:        "invalid_res_negative",
			origin:      LatLng{Lat: 0, Lng: 0},
			destination: LatLng{Lat: Deg(1), Lng: Deg(1)},
			res:         -1,
			expectErr:   E_RES_DOMAIN,
		},
		{
			name:        "invalid_res_too_high",
			origin:      LatLng{Lat: 0, Lng: 0},
			destination: LatLng{Lat: Deg(1), Lng: Deg(1)},
			res:         16,
			expectErr:   E_RES_DOMAIN,
		},
		{
			name:        "north_pole_to_south_pole",
			origin:      LatLng{Lat: Deg(89.9), Lng: 0},
			destination: LatLng{Lat: Deg(-89.9), Lng: 0},
			res:         1,
			expectErr:   E_SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goOut, cOut int64
			goErr := lineHexEstimate(&tt.origin, &tt.destination, tt.res, &goOut)
			cErr := lineHexEstimateC(&tt.origin, &tt.destination, tt.res, &cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("lineHexEstimate error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			if tt.expectErr != E_SUCCESS {
				if goErr != tt.expectErr {
					t.Errorf("Expected error %v, got Go=%v", tt.expectErr, goErr)
				}
				return
			}

			// For success cases, outputs should match
			if goOut != cOut {
				t.Errorf("lineHexEstimate outputs differ: Go=%d, C=%d", goOut, cOut)
			}
		})
	}
}
