//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_bboxFromGeoLoop_ParityWithC(t *testing.T) {
	loop := []LatLng{{0.1, -0.2}, {0.4, 0.3}, {-0.3, 0.2}, {-0.2, -0.5}}
	var goBBox BBox
	bboxFromGeoLoop(loop, &goBBox)
	cBBox := bboxFromGeoLoopC(loop)
	if math.Abs(goBBox.North.Rad()-cBBox.North.Rad()) > 1e-15 || math.Abs(goBBox.South.Rad()-cBBox.South.Rad()) > 1e-15 ||
		math.Abs(goBBox.East.Rad()-cBBox.East.Rad()) > 1e-15 || math.Abs(goBBox.West.Rad()-cBBox.West.Rad()) > 1e-15 {
		t.Fatalf("bboxFromGeoLoop mismatch: go=%+v c=%+v", goBBox, cBBox)
	}
}
