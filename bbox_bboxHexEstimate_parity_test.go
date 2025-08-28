//go:build cgo

package h3

import (
	"testing"
)

func Test_bboxHexEstimate_parity(t *testing.T) {
	tests := []struct {
		name string
		bbox BBox
		res  int32
	}{
		{
			name: "small_box_sf",
			bbox: BBox{North: Deg(37.785), South: Deg(37.765), East: Deg(-122.408), West: Deg(-122.428)},
			res:  5,
		},
		{
			name: "medium_box_bay_area",
			bbox: BBox{North: Deg(37.85), South: Deg(37.70), East: Deg(-122.35), West: Deg(-122.50)},
			res:  4,
		},
		{
			name: "large_box_california",
			bbox: BBox{North: Deg(42.0), South: Deg(32.5), East: Deg(-114.1), West: Deg(-124.4)},
			res:  2,
		},
		{
			name: "very_small_box_high_res",
			bbox: BBox{North: Deg(37.7760), South: Deg(37.7750), East: Deg(-122.4180), West: Deg(-122.4190)},
			res:  9,
		},
		{
			name: "cross_antimeridian",
			bbox: BBox{North: Deg(40.0), South: Deg(30.0), East: Deg(-170.0), West: Deg(170.0)},
			res:  3,
		},
		{
			name: "near_poles",
			bbox: BBox{North: Deg(85.0), South: Deg(80.0), East: Deg(10.0), West: Deg(-10.0)},
			res:  2,
		},
		{
			name: "equatorial_strip",
			bbox: BBox{North: Deg(5.0), South: Deg(-5.0), East: Deg(180.0), West: Deg(-180.0)},
			res:  1,
		},
		{
			name: "minimal_box",
			bbox: BBox{North: Deg(0.001), South: Deg(0.000), East: Deg(0.001), West: Deg(0.000)},
			res:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.res < 0 || tt.res > MAX_H3_RES {
				t.Skipf("Invalid resolution: %d", tt.res)
			}

			var goOut, cOut int64
			goErr := bboxHexEstimate(&tt.bbox, tt.res, &goOut)
			cErr := bboxHexEstimateC(&tt.bbox, tt.res, &cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("bboxHexEstimate error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If there was an error, both should agree on the error
			if goErr != E_SUCCESS {
				return
			}

			// Check outputs match exactly for successful cases
			if goOut != cOut {
				t.Errorf("bboxHexEstimate(%+v, res=%d): Go=%d, C=%d",
					tt.bbox, tt.res, goOut, cOut)
			}

			// Sanity check - estimate should be positive
			if goOut <= 0 {
				t.Errorf("bboxHexEstimate should return positive estimate, got %d", goOut)
			}
		})
	}
}

func Test_bboxHexEstimate_edge_cases_parity(t *testing.T) {
	tests := []struct {
		name      string
		bbox      BBox
		res       int32
		expectErr H3Error
	}{
		{
			name:      "invalid_res_negative",
			bbox:      BBox{North: Deg(1.0), South: Deg(0.0), East: Deg(1.0), West: Deg(0.0)},
			res:       -1,
			expectErr: E_RES_DOMAIN,
		},
		{
			name:      "invalid_res_too_high",
			bbox:      BBox{North: Deg(1.0), South: Deg(0.0), East: Deg(1.0), West: Deg(0.0)},
			res:       16,
			expectErr: E_RES_DOMAIN,
		},
		{
			name:      "zero_width",
			bbox:      BBox{North: Deg(1.0), South: Deg(0.0), East: Deg(1.0), West: Deg(1.0)}, // same as east
			res:       5,
			expectErr: E_FAILED,
		},
		{
			name:      "zero_height",
			bbox:      BBox{North: Deg(1.0), South: Deg(1.0), East: Deg(1.0), West: Deg(0.0)}, // same as north
			res:       5,
			expectErr: E_FAILED,
		},
		{
			name:      "global_bbox",
			bbox:      BBox{North: Deg(90.0), South: Deg(-90.0), East: Deg(180.0), West: Deg(-180.0)},
			res:       0,
			expectErr: E_SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goOut, cOut int64
			goErr := bboxHexEstimate(&tt.bbox, tt.res, &goOut)
			cErr := bboxHexEstimateC(&tt.bbox, tt.res, &cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("bboxHexEstimate error mismatch: Go=%v, C=%v", goErr, cErr)
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
				t.Errorf("bboxHexEstimate outputs differ: Go=%d, C=%d", goOut, cOut)
			}

			// Sanity check for successful estimates
			if goOut <= 0 {
				t.Errorf("bboxHexEstimate should return positive estimate, got %d", goOut)
			}
		})
	}
}
