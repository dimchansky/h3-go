//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_scaleBBox_parity(t *testing.T) {
	tests := []struct {
		name  string
		bbox  BBox
		scale float64
	}{
		{"no_scaling", BBox{1.0, -1.0, 1.0, -1.0}, 1.0},
		{"double_size", BBox{0.5, -0.5, 0.5, -0.5}, 2.0},
		{"shrink_half", BBox{1.0, -1.0, 1.0, -1.0}, 0.5},
		{"small_box", BBox{0.1, -0.1, 0.1, -0.1}, 1.5},
		{"large_scale", BBox{0.2, -0.2, 0.2, -0.2}, 5.0},
		{"near_poles", BBox{math.Pi / 2 * 0.9, -math.Pi / 2 * 0.9, 0.5, -0.5}, 1.2},
		{"near_antimeridian", BBox{0.5, -0.5, math.Pi * 0.9, -math.Pi * 0.9}, 1.3},
		{"crosses_antimeridian", BBox{0.5, -0.5, -math.Pi + 0.1, math.Pi - 0.1}, 1.1},
		{"zero_scale", BBox{0.5, -0.5, 0.5, -0.5}, 0.0},
		{"tiny_scale", BBox{0.1, -0.1, 0.1, -0.1}, 0.01},
		{"negative_scale", BBox{0.1, -0.1, 0.1, -0.1}, -1.0},
	}

	const epsilon = 1e-15

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Copy for C implementation
			bboxC := tt.bbox
			scaleBBoxC(&bboxC, tt.scale)

			// Copy for Go implementation
			bboxGo := tt.bbox
			scaleBBox(&bboxGo, tt.scale)

			// Compare results
			if math.Abs(bboxGo.North.Rad()-bboxC.North.Rad()) > epsilon {
				t.Errorf("scaleBBox() North mismatch: Go=%g, C=%g (diff=%g)",
					bboxGo.North.Rad(), bboxC.North.Rad(), math.Abs(bboxGo.North.Rad()-bboxC.North.Rad()))
			}
			if math.Abs(bboxGo.South.Rad()-bboxC.South.Rad()) > epsilon {
				t.Errorf("scaleBBox() South mismatch: Go=%g, C=%g (diff=%g)",
					bboxGo.South.Rad(), bboxC.South.Rad(), math.Abs(bboxGo.South.Rad()-bboxC.South.Rad()))
			}
			if math.Abs(bboxGo.East.Rad()-bboxC.East.Rad()) > epsilon {
				t.Errorf("scaleBBox() East mismatch: Go=%g, C=%g (diff=%g)",
					bboxGo.East.Rad(), bboxC.East.Rad(), math.Abs(bboxGo.East.Rad()-bboxC.East.Rad()))
			}
			if math.Abs(bboxGo.West.Rad()-bboxC.West.Rad()) > epsilon {
				t.Errorf("scaleBBox() West mismatch: Go=%g, C=%g (diff=%g)",
					bboxGo.West.Rad(), bboxC.West.Rad(), math.Abs(bboxGo.West.Rad()-bboxC.West.Rad()))
			}
		})
	}

	// Test that scaling is deterministic
	t.Run("deterministic", func(t *testing.T) {
		bbox := BBox{0.5, -0.5, 0.5, -0.5}
		scale := 1.5

		// Apply scaling twice
		bbox1 := bbox
		scaleBBox(&bbox1, scale)
		bbox2 := bbox
		scaleBBox(&bbox2, scale)

		if bbox1.North != bbox2.North || bbox1.South != bbox2.South ||
			bbox1.East != bbox2.East || bbox1.West != bbox2.West {
			t.Errorf("scaleBBox should be deterministic: first=%+v != second=%+v",
				bbox1, bbox2)
		}
	})

	// Test edge cases
	t.Run("edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			name    string
			bbox    BBox
			scale   float64
			checkFn func(*testing.T, BBox)
		}{
			{
				name:  "clamp_to_north_pole",
				bbox:  BBox{math.Pi / 2 * 0.9, -0.1, 0.1, -0.1},
				scale: 2.0,
				checkFn: func(t *testing.T, result BBox) {
					if result.North > PiOver2 {
						t.Errorf("North should be clamped to π/2, got %g", result.North.Rad())
					}
				},
			},
			{
				name:  "clamp_to_south_pole",
				bbox:  BBox{0.1, -math.Pi / 2 * 0.9, 0.1, -0.1},
				scale: 2.0,
				checkFn: func(t *testing.T, result BBox) {
					if result.South < -PiOver2 {
						t.Errorf("South should be clamped to -π/2, got %g", result.South.Rad())
					}
				},
			},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				result := tc.bbox
				scaleBBox(&result, tc.scale)
				tc.checkFn(t, result)
			})
		}
	})
}
