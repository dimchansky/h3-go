//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_bboxesFromGeoPolygon_parity(t *testing.T) {
	tests := []struct {
		name    string
		polygon GeoPolygon
	}{
		{
			"simple_triangle",
			GeoPolygon{
				Geoloop: GeoLoop{
					{0.1, 0.1},
					{0.1, 0.2},
					{0.2, 0.1},
				},
			},
		},
		{
			"square_no_holes",
			GeoPolygon{
				Geoloop: GeoLoop{
					{0.0, 0.0},
					{0.0, 0.1},
					{0.1, 0.1},
					{0.1, 0.0},
				},
			},
		},
		{
			"polygon_with_one_hole",
			GeoPolygon{
				Geoloop: GeoLoop{
					{0.0, 0.0},
					{0.0, 0.2},
					{0.2, 0.2},
					{0.2, 0.0},
				},
				Holes: []GeoLoop{
					{
						{0.05, 0.05},
						{0.05, 0.15},
						{0.15, 0.15},
						{0.15, 0.05},
					},
				},
			},
		},
		{
			"polygon_with_multiple_holes",
			GeoPolygon{
				Geoloop: GeoLoop{
					{0.0, 0.0},
					{0.0, 0.3},
					{0.3, 0.3},
					{0.3, 0.0},
				},
				Holes: []GeoLoop{
					{
						{0.05, 0.05},
						{0.05, 0.1},
						{0.1, 0.1},
						{0.1, 0.05},
					},
					{
						{0.2, 0.2},
						{0.2, 0.25},
						{0.25, 0.25},
						{0.25, 0.2},
					},
				},
			},
		},
		{
			"empty_polygon",
			GeoPolygon{
				Geoloop: GeoLoop{},
			},
		},
	}

	const epsilon = 1e-15

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate expected number of bboxes (1 for main loop + 1 per hole)
			numBboxes := 1 + len(tt.polygon.Holes)

			// Create arrays for both implementations
			bboxesC := make([]BBox, numBboxes)
			bboxesGo := make([]BBox, numBboxes)

			// Call C implementation
			bboxesFromGeoPolygonC(&tt.polygon, bboxesC)

			// Call Go implementation
			bboxesFromGeoPolygon(&tt.polygon, bboxesGo)

			// Compare results
			for i := 0; i < numBboxes; i++ {
				if math.Abs(bboxesGo[i].North-bboxesC[i].North) > epsilon {
					t.Errorf("bboxesFromGeoPolygon() bbox[%d].North mismatch: Go=%g, C=%g",
						i, bboxesGo[i].North, bboxesC[i].North)
				}
				if math.Abs(bboxesGo[i].South-bboxesC[i].South) > epsilon {
					t.Errorf("bboxesFromGeoPolygon() bbox[%d].South mismatch: Go=%g, C=%g",
						i, bboxesGo[i].South, bboxesC[i].South)
				}
				if math.Abs(bboxesGo[i].East-bboxesC[i].East) > epsilon {
					t.Errorf("bboxesFromGeoPolygon() bbox[%d].East mismatch: Go=%g, C=%g",
						i, bboxesGo[i].East, bboxesC[i].East)
				}
				if math.Abs(bboxesGo[i].West-bboxesC[i].West) > epsilon {
					t.Errorf("bboxesFromGeoPolygon() bbox[%d].West mismatch: Go=%g, C=%g",
						i, bboxesGo[i].West, bboxesC[i].West)
				}
			}
		})
	}

	// Test deterministic behavior
	t.Run("deterministic", func(t *testing.T) {
		polygon := GeoPolygon{
			Geoloop: GeoLoop{
				{0.1, 0.1},
				{0.1, 0.2},
				{0.2, 0.1},
			},
			Holes: []GeoLoop{
				{
					{0.12, 0.12},
					{0.12, 0.13},
					{0.13, 0.12},
				},
			},
		}

		bboxes1 := make([]BBox, 2)
		bboxes2 := make([]BBox, 2)

		bboxesFromGeoPolygon(&polygon, bboxes1)
		bboxesFromGeoPolygon(&polygon, bboxes2)

		for i := 0; i < 2; i++ {
			if bboxes1[i] != bboxes2[i] {
				t.Errorf("bboxesFromGeoPolygon should be deterministic: bbox[%d] first=%+v != second=%+v",
					i, bboxes1[i], bboxes2[i])
			}
		}
	})
}
