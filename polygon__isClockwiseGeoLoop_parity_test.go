//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_isClockwiseGeoLoop_parity(t *testing.T) {
	tests := []struct {
		name string
		loop GeoLoop
	}{
		{
			name: "clockwise triangle",
			loop: GeoLoop{{Deg(0), Deg(0)}, {Deg(0.1), Deg(0.1)}, {Deg(0), Deg(0.1)}},
		},
		{
			name: "counter-clockwise square",
			loop: GeoLoop{{Deg(0), Deg(0)}, {Deg(0), Deg(0.4)}, {Deg(0.4), Deg(0.4)}, {Deg(0.4), Deg(0)}},
		},
		{
			name: "transmeridian clockwise",
			loop: GeoLoop{
				{Rad(0.4), Rad(math.Pi - 0.1)},
				{Rad(0.4), Rad(-math.Pi + 0.1)},
				{Rad(-0.4), Rad(-math.Pi + 0.1)},
				{Rad(-0.4), Rad(math.Pi - 0.1)},
			},
		},
		{
			name: "transmeridian counter-clockwise",
			loop: GeoLoop{
				{Rad(0.4), Rad(math.Pi - 0.1)},
				{Rad(-0.4), Rad(math.Pi - 0.1)},
				{Rad(-0.4), Rad(-math.Pi + 0.1)},
				{Rad(0.4), Rad(-math.Pi + 0.1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := isClockwiseGeoLoop(tt.loop)
			cResult := isClockwiseGeoLoopC(tt.loop)

			t.Logf("%s: Go=%v, C=%v", tt.name, goResult, cResult)

			if goResult != cResult {
				t.Errorf("Go/C mismatch for %s: Go=%v, C=%v", tt.name, goResult, cResult)
			}
		})
	}
}
