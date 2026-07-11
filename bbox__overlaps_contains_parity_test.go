//go:build cgo && c2go

package h3

import "testing"

func Test_bboxOverlapsBBox_ParityWithC(t *testing.T) {
	a := BBox{North: 1.0, South: -0.5, East: 3.0, West: -2.0}
	b := BBox{North: 0.8, South: -0.8, East: -3.1, West: 2.9} // transmeridian
	if bboxOverlapsBBox(&a, &b) != bboxOverlapsBBoxC(a, b) {
		t.Fatalf("bboxOverlapsBBox mismatch")
	}
}

func Test_bboxContainsBBox_ParityWithC(t *testing.T) {
	a := BBox{North: 1.0, South: -1.0, East: 3.1, West: -3.1}
	b := BBox{North: 0.5, South: -0.5, East: 2.9, West: -2.9}
	if bboxContainsBBox(&a, &b) != bboxContainsBBoxC(a, b) {
		t.Fatalf("bboxContainsBBox mismatch")
	}
}
