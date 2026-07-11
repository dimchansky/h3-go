//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_maxPolygonToCellsSize_parity(t *testing.T) {
	tests := []struct {
		name        string
		geoPolygon  GeoPolygon
		res         int32
		flags       uint32
		expectError H3Error
	}{
		{
			name: "simple_triangle_sf",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
					{Deg(37.780), Deg(-122.420)},
					{Deg(37.770), Deg(-122.415)},
				},
				Holes: nil,
			},
			res:         5,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_SUCCESS,
		},
		{
			name: "square_polygon",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.770), Deg(-122.420)},
					{Deg(37.770), Deg(-122.410)},
					{Deg(37.780), Deg(-122.410)},
					{Deg(37.780), Deg(-122.420)},
				},
				Holes: nil,
			},
			res:         7,
			flags:       uint32(CONTAINMENT_FULL),
			expectError: E_SUCCESS,
		},
		{
			name: "polygon_with_hole",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.760), Deg(-122.430)},
					{Deg(37.760), Deg(-122.400)},
					{Deg(37.790), Deg(-122.400)},
					{Deg(37.790), Deg(-122.430)},
				},
				Holes: []GeoLoop{
					{
						{Deg(37.770), Deg(-122.420)},
						{Deg(37.770), Deg(-122.410)},
						{Deg(37.780), Deg(-122.410)},
						{Deg(37.780), Deg(-122.420)},
					},
				},
			},
			res:         6,
			flags:       uint32(CONTAINMENT_OVERLAPPING),
			expectError: E_SUCCESS,
		},
		{
			name: "single_point_polygon",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
				},
				Holes: nil,
			},
			res:         10,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_SUCCESS,
		},
		{
			name: "large_polygon_low_res",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.70), Deg(-122.50)},
					{Deg(37.70), Deg(-122.30)},
					{Deg(37.90), Deg(-122.30)},
					{Deg(37.90), Deg(-122.50)},
				},
				Holes: nil,
			},
			res:         2,
			flags:       uint32(CONTAINMENT_OVERLAPPING_BBOX),
			expectError: E_SUCCESS,
		},
		{
			name: "tiny_polygon_high_res",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.7750), Deg(-122.4180)},
					{Deg(37.7751), Deg(-122.4181)},
					{Deg(37.7749), Deg(-122.4179)},
				},
				Holes: nil,
			},
			res:         12,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_SUCCESS,
		},
		{
			name: "multiple_holes",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.750), Deg(-122.450)},
					{Deg(37.750), Deg(-122.380)},
					{Deg(37.800), Deg(-122.380)},
					{Deg(37.800), Deg(-122.450)},
				},
				Holes: []GeoLoop{
					{
						{Deg(37.765), Deg(-122.430)},
						{Deg(37.765), Deg(-122.420)},
						{Deg(37.775), Deg(-122.420)},
						{Deg(37.775), Deg(-122.430)},
					},
					{
						{Deg(37.780), Deg(-122.410)},
						{Deg(37.780), Deg(-122.400)},
						{Deg(37.790), Deg(-122.400)},
						{Deg(37.790), Deg(-122.410)},
					},
				},
			},
			res:         8,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_SUCCESS,
		},
		{
			name: "empty_polygon",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{},
				Holes:   nil,
			},
			res:         5,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goOut, cOut int64

			// Call Go implementation
			goErr := maxPolygonToCellsSize(&tt.geoPolygon, tt.res, tt.flags, &goOut)

			// Call C implementation
			cErr := maxPolygonToCellsSizeC(&tt.geoPolygon, tt.res, tt.flags, &cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("maxPolygonToCellsSize error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If there was an expected error, verify it
			if tt.expectError != E_SUCCESS {
				if goErr != tt.expectError {
					t.Errorf("Expected error %v, got Go=%v", tt.expectError, goErr)
				}
				return
			}

			// For successful cases, outputs should match exactly
			if goOut != cOut {
				t.Errorf("maxPolygonToCellsSize(%s, res=%d, flags=0x%x): Go=%d, C=%d",
					tt.name, tt.res, tt.flags, goOut, cOut)
			}

			// Sanity check - size should be positive for polygons with area
			// Single points and empty polygons may return just the buffer size
			if len(tt.geoPolygon.GeoLoop) > 1 && goOut <= POLYGON_TO_CELLS_BUFFER {
				t.Errorf("maxPolygonToCellsSize should return more than buffer (%d) for polygon with area, got %d", POLYGON_TO_CELLS_BUFFER, goOut)
			}

			// Sanity check - for successful cases output should be reasonable
			if goOut < 0 {
				t.Errorf("maxPolygonToCellsSize should not return negative size, got %d", goOut)
			}
		})
	}
}

func Test_maxPolygonToCellsSize_edge_cases_parity(t *testing.T) {
	tests := []struct {
		name        string
		geoPolygon  GeoPolygon
		res         int32
		flags       uint32
		expectError H3Error
	}{
		{
			name: "invalid_flags",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
					{Deg(37.780), Deg(-122.420)},
					{Deg(37.770), Deg(-122.415)},
				},
				Holes: nil,
			},
			res:         5,
			flags:       16, // Invalid flag
			expectError: E_OPTION_INVALID,
		},
		{
			name: "invalid_containment_mode",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
					{Deg(37.780), Deg(-122.420)},
					{Deg(37.770), Deg(-122.415)},
				},
				Holes: nil,
			},
			res:         5,
			flags:       uint32(CONTAINMENT_INVALID), // Invalid containment
			expectError: E_OPTION_INVALID,
		},
		{
			name: "negative_resolution",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
					{Deg(37.780), Deg(-122.420)},
					{Deg(37.770), Deg(-122.415)},
				},
				Holes: nil,
			},
			res:         -1,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_RES_DOMAIN, // This might depend on bboxHexEstimate validation
		},
		{
			name: "resolution_too_high",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
					{Deg(37.780), Deg(-122.420)},
					{Deg(37.770), Deg(-122.415)},
				},
				Holes: nil,
			},
			res:         16, // Above MAX_H3_RES
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_RES_DOMAIN, // This might depend on bboxHexEstimate validation
		},
		{
			name: "polygon_with_many_vertices",
			geoPolygon: GeoPolygon{
				GeoLoop: generateCirclePolygon(37.775, -122.418, 0.01, 100), // 100 vertices
				Holes:   nil,
			},
			res:         5,
			flags:       uint32(CONTAINMENT_CENTER),
			expectError: E_SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goOut, cOut int64

			// Call Go implementation
			goErr := maxPolygonToCellsSize(&tt.geoPolygon, tt.res, tt.flags, &goOut)

			// Call C implementation
			cErr := maxPolygonToCellsSizeC(&tt.geoPolygon, tt.res, tt.flags, &cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("maxPolygonToCellsSize error mismatch: Go=%v, C=%v for case %s", goErr, cErr, tt.name)
				return
			}

			// If there was an expected error, verify it matches
			if tt.expectError != E_SUCCESS {
				if goErr != tt.expectError {
					t.Errorf("Expected error %v, got Go=%v for case %s", tt.expectError, goErr, tt.name)
				}
				return
			}

			// For successful cases, outputs should match exactly
			if goOut != cOut {
				t.Errorf("maxPolygonToCellsSize outputs differ: Go=%d, C=%d for case %s", goOut, cOut, tt.name)
			}
		})
	}
}

// generateCirclePolygon creates a polygon approximating a circle for testing
func generateCirclePolygon(centerLat, centerLng, radius float64, numVertices int) []LatLng {
	polygon := make([]LatLng, numVertices)
	for i := 0; i < numVertices; i++ {
		// Create a simple polygon with varying vertices
		lat := centerLat + radius*1.0*0.5*float64(i%2+1) // Varies radius slightly
		lng := centerLng + radius*1.2*0.5*float64((i+1)%2+1)
		polygon[i] = LatLng{Lat: Deg(lat), Lng: Deg(lng)}
	}
	return polygon
}
