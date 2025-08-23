//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_bboxIsTransmeridian_ParityWithC(t *testing.T) {
	cases := []BBox{{North: 1, South: 0, East: 1, West: 0}, {North: 1, South: 0, East: -3, West: 3}}
	for _, b := range cases {
		goVal := bboxIsTransmeridian(&b)
		cVal := bboxIsTransmeridianC(b)
		if goVal != cVal {
			t.Fatalf("bboxIsTransmeridian mismatch: go=%v c=%v", goVal, cVal)
		}
	}
}

func Test_bboxWidthHeightRads_ParityWithC(t *testing.T) {
	b := BBox{North: 1, South: -1, East: 2, West: -2}
	if math.Abs(bboxWidthRads(&b)-bboxWidthRadsC(b)) > 1e-15 {
		t.Fatalf("bboxWidthRads mismatch")
	}
	if math.Abs(bboxHeightRads(&b)-bboxHeightRadsC(b)) > 1e-15 {
		t.Fatalf("bboxHeightRads mismatch")
	}
}

func Test_bboxEquals_ParityWithC(t *testing.T) {
	b1 := BBox{North: 1, South: -1, East: 2, West: -2}
	b2 := b1
	goVal := bboxEquals(&b1, &b2)
	cVal := bboxEqualsC(b1, b2)
	if goVal != cVal {
		t.Fatalf("bboxEquals mismatch")
	}
}

func Test_bboxCenter_ParityWithC(t *testing.T) {
	b := BBox{North: 1.2, South: -0.4, East: 3.0, West: -2.5}
	goVal := bboxCenter(&b)
	cVal := bboxCenterC(b)
	if math.Abs(goVal.Lat-cVal.Lat) > 1e-15 || math.Abs(goVal.Lng-cVal.Lng) > 1e-15 {
		t.Fatalf("bboxCenter mismatch: go=(%.17g,%.17g) c=(%.17g,%.17g)", goVal.Lat, goVal.Lng, cVal.Lat, cVal.Lng)
	}
}

func Test_bboxContains_ParityWithC(t *testing.T) {
	b := BBox{North: 1.2, South: -0.4, East: 3.0, West: -2.5}
	pts := []LatLng{{Lat: 0, Lng: 0}, {Lat: -0.3, Lng: 2.9}, {Lat: 2, Lng: 0}}
	for _, p := range pts {
		goVal := bboxContains(&b, &p)
		cVal := bboxContainsC(b, p)
		if goVal != cVal {
			t.Fatalf("bboxContains mismatch for p=%+v: go=%v c=%v", p, goVal, cVal)
		}
	}
}
