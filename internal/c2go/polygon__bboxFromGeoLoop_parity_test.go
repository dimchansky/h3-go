//go:build c2go

package c2go

import (
    "math"
    "testing"
)

func Test_bboxFromGeoLoop_ParityWithC(t *testing.T) {
    loop := []LatLng{{0.1, -0.2}, {0.4, 0.3}, {-0.3, 0.2}, {-0.2, -0.5}}
    goBBox := bboxFromGeoLoop(loop)
    cBBox := bboxFromGeoLoopC(loop)
    if math.Abs(goBBox.North-cBBox.North) > 1e-15 || math.Abs(goBBox.South-cBBox.South) > 1e-15 ||
        math.Abs(goBBox.East-cBBox.East) > 1e-15 || math.Abs(goBBox.West-cBBox.West) > 1e-15 {
        t.Fatalf("bboxFromGeoLoop mismatch: go=%+v c=%+v", goBBox, cBBox)
    }
}

