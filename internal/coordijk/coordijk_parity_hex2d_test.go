//go:build oracle

package coordijk

import (
	"math"
	"testing"

	testoracle "github.com/dimchansky/h3-go/internal/testoracle"
	"github.com/dimchansky/h3-go/internal/v2d"
)

func TestOracle_Hex2d_IJKToHex2d(t *testing.T) {
	o := testoracle.New(t)
	cases := []CoordIJK{
		{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0},
		{-1, 0, 0}, {0, -1, 0}, {2, 1, -1}, {3, -2, 1}, {10, 5, -3},
	}
	const tol = 1e-12
	for _, v := range cases {
		got := v.ToHex2d()
		x, y := o.IJKToHex2d([3]int{v.I, v.J, v.K})
		if math.Abs(got.X-x) > tol || math.Abs(got.Y-y) > tol {
			t.Fatalf("ToHex2d(%v) = (%.15g, %.15g), want (%.15g, %.15g)", v, got.X, got.Y, x, y)
		}
	}
}

func TestOracle_Hex2d_Hex2dToCoordIJK(t *testing.T) {
	o := testoracle.New(t)
	type pt struct{ X, Y float64 }
	cases := []pt{
		{0, 0}, {1, 0}, {0, 1}, {1, 1}, {-1, 0}, {0, -1}, {-1, -1},
		{0.5, 0.866025}, {2.5, 1.299038}, {-2.5, -1.299038},
	}
	for _, p := range cases {
		got := Hex2dToCoordIJK(v2d.Vec2d{X: p.X, Y: p.Y})
		arr := o.Hex2dToCoordIJK(p.X, p.Y)
		want := CoordIJK{arr[0], arr[1], arr[2]}
		if got != want {
			t.Fatalf("Hex2dToCoordIJK(%.6f,%.6f) = %v, want %v", p.X, p.Y, got, want)
		}
	}
}
