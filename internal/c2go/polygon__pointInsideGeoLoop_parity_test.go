//go:build c2go

package c2go

import "testing"

func Test_pointInsideGeoLoop_ParityWithC(t *testing.T) {
	loop := []LatLng{{0.0, 0.0}, {0.0, 1.0}, {1.0, 1.0}, {1.0, 0.0}}
	bbox := bboxFromGeoLoop(loop)
	pts := []LatLng{{0.5, 0.5}, {1.5, 0.5}, {0.5, -0.5}, {0.1, 0.9}}
	for _, p := range pts {
		goVal := pointInsideGeoLoop(loop, &bbox, &p)
		cVal := pointInsideGeoLoopC(loop, bbox, p)
		if goVal != cVal {
			t.Fatalf("pointInsideGeoLoop mismatch for p=%+v: go=%v c=%v", p, goVal, cVal)
		}
	}
}
