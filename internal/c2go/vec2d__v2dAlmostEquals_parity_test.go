//go:build c2go

package c2go

import "testing"

func Test_v2dAlmostEquals_ParityWithC(t *testing.T) {
    v1 := Vec2d{1.0, 2.0}
    v2 := Vec2d{1.0 + 1e-8, 2.0 - 1e-8}
    if _v2dAlmostEquals(v1, v2) != v2dAlmostEqualsC(v1, v2) {
        t.Fatalf("_v2dAlmostEquals mismatch")
    }
}

