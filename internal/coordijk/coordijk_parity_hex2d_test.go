//go:build oracle

package coordijk

import (
    "math"
    "testing"
)

func TestOracle_Hex2d_IJKToHex2D(t *testing.T) {
    o := newOracle(t)
    cases := []CoordIJK{
        {0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0},
        {-1, 0, 0}, {0, -1, 0}, {2, 1, -1}, {3, -2, 1}, {10, 5, -3},
    }
    const tol = 1e-12
    for _, v := range cases {
        got := v.ToHex2d()
        want := o.IJKToHex2D(v)
        if math.Abs(got.X-want.X) > tol || math.Abs(got.Y-want.Y) > tol {
            t.Fatalf("ToHex2d(%v) = (%.15g, %.15g), want (%.15g, %.15g)", v, got.X, got.Y, want.X, want.Y)
        }
    }
}

func TestOracle_Hex2d_Hex2DToIJK(t *testing.T) {
    o := newOracle(t)
    type pt struct{ X, Y float64 }
    cases := []pt{
        {0, 0}, {1, 0}, {0, 1}, {1, 1}, {-1, 0}, {0, -1}, {-1, -1},
        {0.5, 0.866025}, {2.5, 1.299038}, {-2.5, -1.299038},
    }
    for _, p := range cases {
        got := Hex2dToCoordIJK(Vec2d{X: p.X, Y: p.Y})
        want := o.Hex2DToIJK(p.X, p.Y)
        if got != want {
            t.Fatalf("Hex2dToCoordIJK(%.6f,%.6f) = %v, want %v", p.X, p.Y, got, want)
        }
    }
}

