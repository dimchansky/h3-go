// Tests ported from H3 v4.4.0: src/apps/testapps/testVec2dInternal.c.
package h3

import (
	"math"
	"testing"
)

func Test_v2dMag(t *testing.T) {
	t.Parallel()
	v := vec2d{X: 3.0, Y: 4.0}
	expected := 5.0
	mag := _v2dMag(&v)
	if math.Abs(mag-expected) >= 1e-15 { // Using double precision epsilon equivalent
		t.Errorf("magnitude not as expected: got %g, want %g", mag, expected)
	}
}

func Test_v2dIntersect(t *testing.T) {
	t.Parallel()
	p0 := vec2d{X: 2.0, Y: 2.0}
	p1 := vec2d{X: 6.0, Y: 6.0}
	p2 := vec2d{X: 0.0, Y: 4.0}
	p3 := vec2d{X: 10.0, Y: 4.0}

	intersection := _v2dIntersect(&p0, &p1, &p2, &p3)

	expectedX := 4.0
	expectedY := 4.0

	if math.Abs(intersection.X-expectedX) >= 1e-15 {
		t.Errorf("X coord not as expected: got %g, want %g", intersection.X, expectedX)
	}
	if math.Abs(intersection.Y-expectedY) >= 1e-15 {
		t.Errorf("Y coord not as expected: got %g, want %g", intersection.Y, expectedY)
	}
}

func Test_v2dAlmostEquals(t *testing.T) {
	t.Parallel()
	v1 := vec2d{X: 3.0, Y: 4.0}
	v2 := vec2d{X: 3.0, Y: 4.0}
	v3 := vec2d{X: 3.5, Y: 4.0}
	v4 := vec2d{X: 3.0, Y: 4.5}
	const dblEpsilon = 2.2204460492503131e-16 // dblEpsilon
	v5 := vec2d{X: 3.0 + dblEpsilon, Y: 4.0 - dblEpsilon}

	if !_v2dAlmostEquals(&v1, &v2) {
		t.Error("expected true for equal vectors")
	}
	if _v2dAlmostEquals(&v1, &v3) {
		t.Error("expected false for different x")
	}
	if _v2dAlmostEquals(&v1, &v4) {
		t.Error("expected false for different y")
	}
	if !_v2dAlmostEquals(&v1, &v5) {
		t.Error("expected true for almost equal vectors")
	}
}
