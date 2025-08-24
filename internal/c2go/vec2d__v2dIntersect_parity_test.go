//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_v2dIntersect_ParityWithC(t *testing.T) {
	p0 := Vec2d{0, 0}
	p1 := Vec2d{1, 1}
	p2 := Vec2d{0, 1}
	p3 := Vec2d{1, 0}
	goVal := _v2dIntersect(&p0, &p1, &p2, &p3)
	cVal := v2dIntersectC(p0, p1, p2, p3)
	if math.Abs(goVal.X-cVal.X) > 1e-15 || math.Abs(goVal.Y-cVal.Y) > 1e-15 {
		t.Fatalf("_v2dIntersect mismatch: go=(%.17g,%.17g) c=(%.17g,%.17g)", goVal.X, goVal.Y, cVal.X, cVal.Y)
	}
}
