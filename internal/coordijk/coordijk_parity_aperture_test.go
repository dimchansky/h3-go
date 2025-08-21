//go:build oracle

package coordijk

import (
	"math"
	"testing"

	testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

func TestOracle_ApertureTransforms(t *testing.T) {
	o := testoracle.New(t)
	inputs := []CoordIJK{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {2, 1, 0}, {7, 0, 0}, {0, 7, 0}}
	for _, v := range inputs {
		got := v
		got.UpAp7()
		a := o.UpAp7([3]int{v.I, v.J, v.K})
		want := CoordIJK{a[0], a[1], a[2]}
		if got != want {
			t.Fatalf("UpAp7(%v)=%v, want %v", v, got, want)
		}
		got = v
		got.UpAp7r()
		a = o.UpAp7r([3]int{v.I, v.J, v.K})
		want = CoordIJK{a[0], a[1], a[2]}
		if got != want {
			t.Fatalf("UpAp7r(%v)=%v, want %v", v, got, want)
		}
		got = v
		got.DownAp7()
		a = o.DownAp7([3]int{v.I, v.J, v.K})
		want = CoordIJK{a[0], a[1], a[2]}
		if got != want {
			t.Fatalf("DownAp7(%v)=%v, want %v", v, got, want)
		}
		got = v
		got.DownAp7r()
		a = o.DownAp7r([3]int{v.I, v.J, v.K})
		want = CoordIJK{a[0], a[1], a[2]}
		if got != want {
			t.Fatalf("DownAp7r(%v)=%v, want %v", v, got, want)
		}
		got = v
		got.DownAp3()
		a = o.DownAp3([3]int{v.I, v.J, v.K})
		want = CoordIJK{a[0], a[1], a[2]}
		if got != want {
			t.Fatalf("DownAp3(%v)=%v, want %v", v, got, want)
		}
		got = v
		got.DownAp3r()
		a = o.DownAp3r([3]int{v.I, v.J, v.K})
		want = CoordIJK{a[0], a[1], a[2]}
		if got != want {
			t.Fatalf("DownAp3r(%v)=%v, want %v", v, got, want)
		}
	}

	// Extremes & round-trip sanity using hex2d equivalence
	extremes := []int{0, 1, 7, 21}
	for _, i := range extremes {
		v := CoordIJK{i, 0, 0}
		// Down then up (CCW)
		g := v
		g.DownAp7()
		g.UpAp7()
		arr := o.DownAp7([3]int{v.I, v.J, v.K})
		arr = o.UpAp7(arr)
		og := CoordIJK{arr[0], arr[1], arr[2]}
		if g.ToHex2d() != og.ToHex2d() {
			t.Fatalf("roundtrip CCW mismatch for %v", v)
		}
		// Down then up (CW)
		g = v
		g.DownAp7r()
		g.UpAp7r()
		arr = o.DownAp7r([3]int{v.I, v.J, v.K})
		arr = o.UpAp7r(arr)
		og = CoordIJK{arr[0], arr[1], arr[2]}
		h1, h2 := g.ToHex2d(), og.ToHex2d()
		if math.Abs(h1.X-h2.X) > 1e-12 || math.Abs(h1.Y-h2.Y) > 1e-12 {
			t.Fatalf("roundtrip CW mismatch for %v", v)
		}
	}

	// Randomized IJK+ sweep
	r := testoracle.NewRand()
	n := testoracle.Max()
	for i := 0; i < n; i++ {
		arr := testoracle.RandIJKPlus(r, 30)
		v := CoordIJK{arr[0], arr[1], arr[2]}
		if g, w := func() (CoordIJK, CoordIJK) {
			x := v
			x.UpAp7()
			a := o.UpAp7([3]int{v.I, v.J, v.K})
			return x, CoordIJK{a[0], a[1], a[2]}
		}(); g != w {
			t.Fatalf("UpAp7 parity %v", v)
		}
		if g, w := func() (CoordIJK, CoordIJK) {
			x := v
			x.UpAp7r()
			a := o.UpAp7r([3]int{v.I, v.J, v.K})
			return x, CoordIJK{a[0], a[1], a[2]}
		}(); g != w {
			t.Fatalf("UpAp7r parity %v", v)
		}
		if g, w := func() (CoordIJK, CoordIJK) {
			x := v
			x.DownAp7()
			a := o.DownAp7([3]int{v.I, v.J, v.K})
			return x, CoordIJK{a[0], a[1], a[2]}
		}(); g != w {
			t.Fatalf("DownAp7 parity %v", v)
		}
		if g, w := func() (CoordIJK, CoordIJK) {
			x := v
			x.DownAp7r()
			a := o.DownAp7r([3]int{v.I, v.J, v.K})
			return x, CoordIJK{a[0], a[1], a[2]}
		}(); g != w {
			t.Fatalf("DownAp7r parity %v", v)
		}
		if g, w := func() (CoordIJK, CoordIJK) {
			x := v
			x.DownAp3()
			a := o.DownAp3([3]int{v.I, v.J, v.K})
			return x, CoordIJK{a[0], a[1], a[2]}
		}(); g != w {
			t.Fatalf("DownAp3 parity %v", v)
		}
		if g, w := func() (CoordIJK, CoordIJK) {
			x := v
			x.DownAp3r()
			a := o.DownAp3r([3]int{v.I, v.J, v.K})
			return x, CoordIJK{a[0], a[1], a[2]}
		}(); g != w {
			t.Fatalf("DownAp3r parity %v", v)
		}
	}
}
