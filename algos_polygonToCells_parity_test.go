//go:build cgo

package h3

import (
	"sort"
	"testing"
)

func Test_polygonToCells_parity(t *testing.T) {
	tests := []struct {
		name       string
		geoPolygon GeoPolygon
		res        int32
		flags      uint32
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
			res:   5,
			flags: uint32(CONTAINMENT_CENTER),
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
			res:   7,
			flags: uint32(CONTAINMENT_FULL),
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
			res:   6,
			flags: uint32(CONTAINMENT_OVERLAPPING),
		},
		{
			name: "single_point_polygon",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{
					{Deg(37.775), Deg(-122.418)},
				},
				Holes: nil,
			},
			res:   10,
			flags: uint32(CONTAINMENT_CENTER),
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
			res:   2,
			flags: uint32(CONTAINMENT_OVERLAPPING_BBOX),
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
			res:   12,
			flags: uint32(CONTAINMENT_CENTER),
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
			res:   8,
			flags: uint32(CONTAINMENT_CENTER),
		},
		{
			name: "empty_polygon",
			geoPolygon: GeoPolygon{
				GeoLoop: []LatLng{},
				Holes:   nil,
			},
			res:   5,
			flags: uint32(CONTAINMENT_CENTER),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the expected size first - check parity here too
			var goSize, cSize int64
			goSizeErr := maxPolygonToCellsSize(&tt.geoPolygon, tt.res, tt.flags, &goSize)
			cSizeErr := maxPolygonToCellsSizeC(&tt.geoPolygon, tt.res, tt.flags, &cSize)

			// Size estimation should match between Go and C
			if goSizeErr != cSizeErr {
				t.Errorf("maxPolygonToCellsSize error mismatch: Go=%v, C=%v", goSizeErr, cSizeErr)
				return
			}

			if goSizeErr != E_SUCCESS {
				// Both failed with same error - polygonToCells should also fail the same way
				// But let's test it anyway to make sure
				goOut := make([]H3Index, 1) // Small buffer to test error handling
				cOut := make([]H3Index, 1)

				goErr := polygonToCells(&tt.geoPolygon, tt.res, tt.flags, goOut)
				cErr := polygonToCellsC(&tt.geoPolygon, tt.res, tt.flags, cOut)

				if goErr != cErr {
					t.Errorf("polygonToCells error mismatch: Go=%v, C=%v", goErr, cErr)
				}
				return
			}

			// Both size calls succeeded, check they got same size
			if goSize != cSize {
				t.Errorf("maxPolygonToCellsSize result mismatch: Go=%d, C=%d", goSize, cSize)
				return
			}

			// Allocate output buffers
			goOut := make([]H3Index, goSize)
			cOut := make([]H3Index, cSize)

			// Call Go implementation
			goErr := polygonToCells(&tt.geoPolygon, tt.res, tt.flags, goOut)

			// Call C implementation
			cErr := polygonToCellsC(&tt.geoPolygon, tt.res, tt.flags, cOut)

			// Check errors match
			if goErr != cErr {
				t.Errorf("polygonToCells error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If both failed, we're done
			if goErr != E_SUCCESS {
				return
			}

			// For successful cases, extract valid cells and compare
			goValidCells := extractValidCells(goOut)
			cValidCells := extractValidCells(cOut)

			// Sort both slices for comparison
			sort.Slice(goValidCells, func(i, j int) bool { return goValidCells[i] < goValidCells[j] })
			sort.Slice(cValidCells, func(i, j int) bool { return cValidCells[i] < cValidCells[j] })

			// Compare lengths first
			if len(goValidCells) != len(cValidCells) {
				t.Errorf("polygonToCells(%s, res=%d, flags=0x%x): cell count mismatch Go=%d, C=%d",
					tt.name, tt.res, tt.flags, len(goValidCells), len(cValidCells))
				return
			}

			// Compare each cell
			for i := 0; i < len(goValidCells); i++ {
				if goValidCells[i] != cValidCells[i] {
					t.Errorf("polygonToCells(%s, res=%d, flags=0x%x): cell mismatch at index %d Go=%016x, C=%016x",
						tt.name, tt.res, tt.flags, i, goValidCells[i], cValidCells[i])
					return
				}
			}

			// Sanity check - non-empty polygons should produce some cells
			if len(tt.geoPolygon.GeoLoop) > 0 && len(goValidCells) == 0 {
				t.Logf("polygonToCells(%s): Warning - non-empty polygon produced no cells", tt.name)
			}
		})
	}
}

// extractValidCells filters out H3_NULL values from a slice of H3Index
func extractValidCells(cells []H3Index) []H3Index {
	var valid []H3Index
	for _, cell := range cells {
		if cell != H3_NULL {
			valid = append(valid, cell)
		}
	}
	return valid
}
