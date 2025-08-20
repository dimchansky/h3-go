//go:build oracle

package coordijk

import (
    "math"
    "testing"
)

func TestOracle_ApertureTransforms(t *testing.T) {
    o := newOracle(t)
    inputs := []CoordIJK{{0,0,0},{1,0,0},{0,1,0},{0,0,1},{2,1,0},{7,0,0},{0,7,0}}
    for _, v := range inputs {
        got := v; got.UpAp7();  want := o.UpAp7(v);  if got != want { t.Fatalf("UpAp7(%v)=%v, want %v", v, got, want) }
        got = v; got.UpAp7r(); want = o.UpAp7r(v); if got != want { t.Fatalf("UpAp7r(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp7();  want = o.DownAp7(v);  if got != want { t.Fatalf("DownAp7(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp7r(); want = o.DownAp7r(v); if got != want { t.Fatalf("DownAp7r(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp3();  want = o.DownAp3(v);  if got != want { t.Fatalf("DownAp3(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp3r(); want = o.DownAp3r(v); if got != want { t.Fatalf("DownAp3r(%v)=%v, want %v", v, got, want) }
    }

    // Extremes & round-trip sanity using hex2d equivalence
    extremes := []int{0,1,7,21}
    for _, i := range extremes {
        v := CoordIJK{i, 0, 0}
        // Down then up (CCW)
        g := v; g.DownAp7(); g.UpAp7()
        og := v; og = o.DownAp7(og); og = o.UpAp7(og)
        if g.ToHex2d() != og.ToHex2d() { t.Fatalf("roundtrip CCW mismatch for %v", v) }
        // Down then up (CW)
        g = v; g.DownAp7r(); g.UpAp7r()
        og = v; og = o.DownAp7r(og); og = o.UpAp7r(og)
        h1, h2 := g.ToHex2d(), og.ToHex2d()
        if math.Abs(h1.X-h2.X) > 1e-12 || math.Abs(h1.Y-h2.Y) > 1e-12 { t.Fatalf("roundtrip CW mismatch for %v", v) }
    }

    // Randomized IJK+ sweep
    r := newRand()
    n := oracleMax()
    for i := 0; i < n; i++ {
        v := randIJKPlus(r, 30)
        if g, w := func() (CoordIJK, CoordIJK) { x:=v; x.UpAp7(); return x, o.UpAp7(v) }(); g != w { t.Fatalf("UpAp7 parity %v", v) }
        if g, w := func() (CoordIJK, CoordIJK) { x:=v; x.UpAp7r(); return x, o.UpAp7r(v) }(); g != w { t.Fatalf("UpAp7r parity %v", v) }
        if g, w := func() (CoordIJK, CoordIJK) { x:=v; x.DownAp7(); return x, o.DownAp7(v) }(); g != w { t.Fatalf("DownAp7 parity %v", v) }
        if g, w := func() (CoordIJK, CoordIJK) { x:=v; x.DownAp7r(); return x, o.DownAp7r(v) }(); g != w { t.Fatalf("DownAp7r parity %v", v) }
        if g, w := func() (CoordIJK, CoordIJK) { x:=v; x.DownAp3(); return x, o.DownAp3(v) }(); g != w { t.Fatalf("DownAp3 parity %v", v) }
        if g, w := func() (CoordIJK, CoordIJK) { x:=v; x.DownAp3r(); return x, o.DownAp3r(v) }(); g != w { t.Fatalf("DownAp3r parity %v", v) }
    }
}
