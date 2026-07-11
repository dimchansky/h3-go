//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_maxPolygonToCellsSizeExperimental_parity(t *testing.T) {
	tests := []struct {
		name    string
		polygon GeoPolygon
		res     int32
		flags   uint32
		wantErr h3Error
	}{
		{
			name: "Empty polygon",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{},
				Holes:   []GeoLoop{},
			},
			res:     5,
			flags:   uint32(ContainmentCenter),
			wantErr: eSuccess,
		},
		{
			name: "Simple triangle polygon",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(37.813318999983238), Lng: Deg(-122.4089866999972145)},
					{Lat: Deg(37.7866302000007224), Lng: Deg(-122.3805436999997056)},
					{Lat: Deg(37.7198061999978478), Lng: Deg(-122.4148903999994803)},
					{Lat: Deg(37.813318999983238), Lng: Deg(-122.4089866999972145)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     9,
			flags:   uint32(ContainmentCenter),
			wantErr: eSuccess,
		},
		{
			name: "Small square polygon at various resolutions",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(0.1), Lng: Deg(0.1)},
					{Lat: Deg(0.1), Lng: Deg(0.2)},
					{Lat: Deg(0.2), Lng: Deg(0.2)},
					{Lat: Deg(0.2), Lng: Deg(0.1)},
					{Lat: Deg(0.1), Lng: Deg(0.1)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     7,
			flags:   uint32(ContainmentCenter),
			wantErr: eSuccess,
		},
		{
			name: "Larger polygon with overlapping bbox mode",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(37.82), Lng: Deg(-122.42)},
					{Lat: Deg(37.82), Lng: Deg(-122.40)},
					{Lat: Deg(37.80), Lng: Deg(-122.40)},
					{Lat: Deg(37.80), Lng: Deg(-122.42)},
					{Lat: Deg(37.82), Lng: Deg(-122.42)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     8,
			flags:   uint32(ContainmentOverlappingBBox),
			wantErr: eSuccess,
		},
		{
			name: "Pentagon test case",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(37.0), Lng: Deg(-122.0)},
					{Lat: Deg(37.1), Lng: Deg(-122.0)},
					{Lat: Deg(37.1), Lng: Deg(-121.9)},
					{Lat: Deg(37.05), Lng: Deg(-121.85)},
					{Lat: Deg(37.0), Lng: Deg(-121.9)},
					{Lat: Deg(37.0), Lng: Deg(-122.0)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     6,
			flags:   uint32(ContainmentCenter),
			wantErr: eSuccess,
		},
		{
			name: "High resolution test",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(0.01), Lng: Deg(0.01)},
					{Lat: Deg(0.01), Lng: Deg(0.02)},
					{Lat: Deg(0.02), Lng: Deg(0.02)},
					{Lat: Deg(0.02), Lng: Deg(0.01)},
					{Lat: Deg(0.01), Lng: Deg(0.01)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     12,
			flags:   uint32(ContainmentCenter),
			wantErr: eSuccess,
		},
		{
			name: "Invalid resolution - too high",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(0.1), Lng: Deg(0.1)},
					{Lat: Deg(0.1), Lng: Deg(0.2)},
					{Lat: Deg(0.2), Lng: Deg(0.2)},
					{Lat: Deg(0.1), Lng: Deg(0.1)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     16, // Invalid: > maxH3Res (15)
			flags:   uint32(ContainmentCenter),
			wantErr: eResDomain,
		},
		{
			name: "Invalid resolution - negative",
			polygon: GeoPolygon{
				GeoLoop: GeoLoop{
					{Lat: Deg(0.1), Lng: Deg(0.1)},
					{Lat: Deg(0.1), Lng: Deg(0.2)},
					{Lat: Deg(0.2), Lng: Deg(0.2)},
					{Lat: Deg(0.1), Lng: Deg(0.1)}, // Close the loop
				},
				Holes: []GeoLoop{},
			},
			res:     -1,
			flags:   uint32(ContainmentCenter),
			wantErr: eResDomain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Go implementation
			goResult, goErr := maxPolygonToCellsSizeExperimental(&tt.polygon, tt.res, tt.flags)

			// Call C implementation
			cResult, cErr := maxPolygonToCellsSizeExperimentalC(&tt.polygon, tt.res, tt.flags)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If we expected an error, don't compare results
			if tt.wantErr != eSuccess {
				if goErr != tt.wantErr {
					t.Errorf("Expected error %v, got %v", tt.wantErr, goErr)
				}
				return
			}

			// Both implementations should return the same result
			if goResult != cResult {
				t.Errorf("Result mismatch: Go=%d, C=%d", goResult, cResult)
			}

			// For successful cases, result should be non-negative
			if goErr == eSuccess && goResult < 0 {
				t.Errorf("Expected non-negative result, got %d", goResult)
			}
		})
	}
}

func Test_maxPolygonToCellsSizeExperimental_resolution_scaling_parity(t *testing.T) {
	// Test that higher resolutions generally give larger estimates
	polygon := GeoPolygon{
		GeoLoop: GeoLoop{
			{Lat: Deg(37.0), Lng: Deg(-122.0)},
			{Lat: Deg(37.1), Lng: Deg(-122.0)},
			{Lat: Deg(37.1), Lng: Deg(-121.9)},
			{Lat: Deg(37.0), Lng: Deg(-121.9)},
			{Lat: Deg(37.0), Lng: Deg(-122.0)}, // Close the loop
		},
		Holes: []GeoLoop{},
	}

	resolutions := []int32{3, 5, 7, 9}
	for _, res := range resolutions {
		t.Run("Resolution_"+string(rune(res+'0')), func(t *testing.T) {
			// Call both implementations
			goResult, goErr := maxPolygonToCellsSizeExperimental(&polygon, res, uint32(ContainmentCenter))
			cResult, cErr := maxPolygonToCellsSizeExperimentalC(&polygon, res, uint32(ContainmentCenter))

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch at res %d: Go=%v, C=%v", res, goErr, cErr)
				return
			}

			if goErr != eSuccess {
				t.Errorf("Unexpected error at res %d: %v", res, goErr)
				return
			}

			// Compare results
			if goResult != cResult {
				t.Errorf("Result mismatch at res %d: Go=%d, C=%d", res, goResult, cResult)
			}

			// Result should be positive for valid polygons
			if goResult <= 0 {
				t.Errorf("Expected positive result at res %d, got %d", res, goResult)
			}
		})
	}
}

func Test_maxPolygonToCellsSizeExperimental_flags_parity(t *testing.T) {
	// Test different containment flags
	polygon := GeoPolygon{
		GeoLoop: GeoLoop{
			{Lat: Deg(0.05), Lng: Deg(0.05)},
			{Lat: Deg(0.05), Lng: Deg(0.1)},
			{Lat: Deg(0.1), Lng: Deg(0.1)},
			{Lat: Deg(0.1), Lng: Deg(0.05)},
			{Lat: Deg(0.05), Lng: Deg(0.05)}, // Close the loop
		},
		Holes: []GeoLoop{},
	}

	flags := []uint32{
		uint32(ContainmentCenter),
		uint32(ContainmentFull),
		uint32(ContainmentOverlappingBBox),
	}

	res := int32(8)
	for _, flag := range flags {
		t.Run("Flag_"+string(rune(flag+'0')), func(t *testing.T) {
			// Call both implementations
			goResult, goErr := maxPolygonToCellsSizeExperimental(&polygon, res, flag)
			cResult, cErr := maxPolygonToCellsSizeExperimentalC(&polygon, res, flag)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch with flag %d: Go=%v, C=%v", flag, goErr, cErr)
				return
			}

			if goErr != eSuccess {
				t.Errorf("Unexpected error with flag %d: %v", flag, goErr)
				return
			}

			// Compare results
			if goResult != cResult {
				t.Errorf("Result mismatch with flag %d: Go=%d, C=%d", flag, goResult, cResult)
			}
		})
	}
}
