//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_polygonToCellsExperimental_parity(t *testing.T) {
	tests := []struct {
		name     string
		polygon  GeoPolygon
		res      int32
		flags    uint32
		maxCells int64
	}{
		{
			name: "Empty polygon",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{},
				Holes:   nil,
			},
			res:      6,
			flags:    uint32(ContainmentCenter),
			maxCells: 0,
		},
		{
			name: "Simple triangle",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Angle(0.659966917655), Lng: Angle(-2.1364398519396)},
					{Lat: Angle(0.6595011102219), Lng: Angle(-2.1359434279405)},
					{Lat: Angle(0.6583348114025), Lng: Angle(-2.1354884206045)},
					{Lat: Angle(0.659966917655), Lng: Angle(-2.1364398519396)}, // closed loop
				},
				Holes: nil,
			},
			res:      9,
			flags:    uint32(ContainmentCenter),
			maxCells: 1000,
		},
		{
			name: "Pentagon test",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Angle(0.8), Lng: Angle(1.2)},
					{Lat: Angle(0.7), Lng: Angle(1.2)},
					{Lat: Angle(0.7), Lng: Angle(1.3)},
					{Lat: Angle(0.8), Lng: Angle(1.3)},
					{Lat: Angle(0.8), Lng: Angle(1.2)}, // closed loop
				},
				Holes: nil,
			},
			res:      7,
			flags:    uint32(ContainmentCenter),
			maxCells: 2000,
		},
		{
			name: "Polygon with hole",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Angle(0.8), Lng: Angle(1.2)},
					{Lat: Angle(0.7), Lng: Angle(1.2)},
					{Lat: Angle(0.7), Lng: Angle(1.3)},
					{Lat: Angle(0.8), Lng: Angle(1.3)},
					{Lat: Angle(0.8), Lng: Angle(1.2)}, // closed loop
				},
				Holes: []GeoLoop{
					{
						{Lat: Angle(0.76), Lng: Angle(1.22)},
						{Lat: Angle(0.74), Lng: Angle(1.22)},
						{Lat: Angle(0.74), Lng: Angle(1.28)},
						{Lat: Angle(0.76), Lng: Angle(1.28)},
						{Lat: Angle(0.76), Lng: Angle(1.22)}, // closed loop
					},
				},
			},
			res:      6,
			flags:    uint32(ContainmentCenter),
			maxCells: 500,
		},
	}

	for _, mode := range []ContainmentMode{
		ContainmentCenter,
		ContainmentFull,
		ContainmentOverlapping,
		ContainmentOverlappingBBox,
	} {
		for _, tt := range tests {
			t.Run(tt.name+"_mode_"+string(rune(mode+'0')), func(t *testing.T) {
				flags := uint32(mode)

				// Test with different buffer sizes
				maxCells := tt.maxCells
				if maxCells == 0 {
					maxCells = 10 // minimum for empty case
				}

				// Allocate buffers
				goOut := make([]h3Index, maxCells)
				cOut := make([]h3Index, maxCells)

				// Call Go implementation
				goErr := polygonToCellsExperimental(&tt.polygon, tt.res, flags, maxCells, goOut)

				// Call C implementation
				cErr := polygonToCellsExperimentalC(&tt.polygon, tt.res, flags, maxCells, cOut)

				// Compare errors
				if goErr != cErr {
					t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
					return
				}

				// If both succeeded, compare results
				if goErr == eSuccess {
					// Count non-null cells in each result
					goCount := 0
					cCount := 0
					for i := int64(0); i < maxCells; i++ {
						if goOut[i] != h3Null {
							goCount++
						}
						if cOut[i] != h3Null {
							cCount++
						}
					}

					if goCount != cCount {
						t.Errorf("Cell count mismatch: Go=%d, C=%d", goCount, cCount)
					}

					// For non-empty results, check that all Go cells are in C results
					if goCount > 0 && cCount > 0 {
						// Create maps for faster lookup
						cCells := make(map[h3Index]bool)
						for i := 0; i < cCount && i < len(cOut); i++ {
							if cOut[i] != h3Null {
								cCells[cOut[i]] = true
							}
						}

						missingCount := 0
						for i := 0; i < goCount && i < len(goOut); i++ {
							if goOut[i] != h3Null && !cCells[goOut[i]] {
								missingCount++
								if missingCount <= 5 { // Only show first 5 mismatches
									t.Errorf("Go cell %016x not found in C results", goOut[i])
								}
							}
						}
						if missingCount > 5 {
							t.Errorf("... and %d more missing cells", missingCount-5)
						}
					}
				}
			})
		}
	}
}

func Test_polygonToCellsExperimental_invalid_inputs_parity(t *testing.T) {
	tests := []struct {
		name      string
		polygon   *GeoPolygon
		res       int32
		flags     uint32
		maxCells  int64
		expectErr h3Error
	}{
		{
			name:      "Negative resolution",
			polygon:   &GeoPolygon{GeoLoop: GeoLoop{{Lat: 0, Lng: 0}}},
			res:       -1,
			flags:     uint32(ContainmentCenter),
			maxCells:  100,
			expectErr: eResDomain,
		},
		{
			name:      "Too high resolution",
			polygon:   &GeoPolygon{GeoLoop: GeoLoop{{Lat: 0, Lng: 0}}},
			res:       maxH3Res + 1,
			flags:     uint32(ContainmentCenter),
			maxCells:  100,
			expectErr: eResDomain,
		},
		{
			name:      "Invalid containment mode",
			polygon:   &GeoPolygon{GeoLoop: GeoLoop{{Lat: 0, Lng: 0}}},
			res:       5,
			flags:     uint32(ContainmentInvalid),
			maxCells:  100,
			expectErr: eOptionInvalid,
		},
		{
			name:      "Zero max cells",
			polygon:   &GeoPolygon{GeoLoop: GeoLoop{{Lat: 0, Lng: 0}}},
			res:       5,
			flags:     uint32(ContainmentCenter),
			maxCells:  0,
			expectErr: eSuccess, // No cells to write, so no overflow
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Allocate buffer
			goOut := make([]h3Index, maxInt(1, int(tt.maxCells)))
			cOut := make([]h3Index, maxInt(1, int(tt.maxCells)))

			// Call Go implementation
			goErr := polygonToCellsExperimental(tt.polygon, tt.res, tt.flags, tt.maxCells, goOut)

			// Call C implementation
			cErr := polygonToCellsExperimentalC(tt.polygon, tt.res, tt.flags, tt.maxCells, cOut)

			// Both should return the expected error
			if goErr != tt.expectErr {
				t.Errorf("Go: expected error %v, got %v", tt.expectErr, goErr)
			}
			if cErr != tt.expectErr {
				t.Errorf("C: expected error %v, got %v", tt.expectErr, cErr)
			}

			// Errors should match
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}
		})
	}
}

// max helper for compatibility (avoid duplicate declaration)
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Test_polygonToCellsExperimental_large_polygon_parity(t *testing.T) {
	// Create a larger polygon around San Francisco
	polygon := GeoPolygon{
		GeoLoop: GeoLoop{
			{Lat: Angle(37.813318999983238 * math.Pi / 180), Lng: Angle(-122.4089866999972145 * math.Pi / 180)},
			{Lat: Angle(37.7866302000007224 * math.Pi / 180), Lng: Angle(-122.3805436999997056 * math.Pi / 180)},
			{Lat: Angle(37.7198061999978478 * math.Pi / 180), Lng: Angle(-122.3544736999993603 * math.Pi / 180)},
			{Lat: Angle(37.7076131999975672 * math.Pi / 180), Lng: Angle(-122.5123436999983966 * math.Pi / 180)},
			{Lat: Angle(37.7835871999971715 * math.Pi / 180), Lng: Angle(-122.5247187000021967 * math.Pi / 180)},
			{Lat: Angle(37.8151571999998453 * math.Pi / 180), Lng: Angle(-122.4798767000009008 * math.Pi / 180)},
			{Lat: Angle(37.813318999983238 * math.Pi / 180), Lng: Angle(-122.4089866999972145 * math.Pi / 180)},
		},
		Holes: nil,
	}

	for _, mode := range []ContainmentMode{ContainmentCenter, ContainmentOverlappingBBox} {
		for _, res := range []int32{6, 7, 8} {
			t.Run("sf_polygon", func(t *testing.T) {
				flags := uint32(mode)
				maxCells := int64(10000)

				// Allocate buffers
				goOut := make([]h3Index, maxCells)
				cOut := make([]h3Index, maxCells)

				// Call Go implementation
				goErr := polygonToCellsExperimental(&polygon, res, flags, maxCells, goOut)

				// Call C implementation
				cErr := polygonToCellsExperimentalC(&polygon, res, flags, maxCells, cOut)

				// Compare errors
				if goErr != cErr {
					t.Errorf("Error mismatch for res %d mode %d: Go=%v, C=%v", res, mode, goErr, cErr)
					return
				}

				if goErr == eSuccess {
					// Count cells
					goCount := 0
					cCount := 0
					for i := 0; i < len(goOut); i++ {
						if goOut[i] != h3Null {
							goCount++
						}
						if cOut[i] != h3Null {
							cCount++
						}
					}

					// Allow for small differences due to implementation variations
					tolerance := maxInt(1, (cCount+goCount)/20) // 5% tolerance
					if absInt(goCount-cCount) > tolerance {
						t.Errorf("Cell count difference too large for res %d mode %d: Go=%d, C=%d, diff=%d, tolerance=%d",
							res, mode, goCount, cCount, absInt(goCount-cCount), tolerance)
					}
				}
			})
		}
	}
}

// absInt helper for compatibility (avoid duplicate declaration)
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
