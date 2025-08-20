//go:build oracle

package coordijk

import (
    "math"
    "testing"
)

func TestOracle_CoordIJKDistance(t *testing.T) {
    o := newOracle(t)
    cases := []struct{ a, b CoordIJK }{
        {CoordIJK{0,0,0}, CoordIJK{0,0,0}},
        {CoordIJK{1,0,0}, CoordIJK{0,0,0}},
        {CoordIJK{0,2,0}, CoordIJK{0,0,0}},
        {CoordIJK{1,0,1}, CoordIJK{0,0,0}},
        {CoordIJK{1,0,1}, CoordIJK{0,2,0}},
        {CoordIJK{2,-1,1}, CoordIJK{-1,2,-1}},
    }
    for _, tc := range cases {
        if Distance(tc.a, tc.b) != o.Distance(tc.a, tc.b) {
            t.Fatalf("Distance(%v,%v) mismatch", tc.a, tc.b)
        }
    }
}

func TestOracle_CoordIJKRotate(t *testing.T) {
    o := newOracle(t)
    cases := []CoordIJK{{0,0,0},{1,0,0},{0,1,0},{0,0,1},{1,1,0},{1,0,1},{0,1,1},{2,-1,1},{-1,2,-1}}
    for _, v := range cases {
        got := v; got.Rotate60CCW(); want := o.Rotate60CCW(v)
        if got != want { t.Fatalf("Rotate60CCW(%v) = %v, want %v", v, got, want) }
        got = v; got.Rotate60CW(); want = o.Rotate60CW(v)
        if got != want { t.Fatalf("Rotate60CW(%v) = %v, want %v", v, got, want) }
    }

    // Randomized parity and properties
    r := newRand()
    n := oracleMax()
    for i := 0; i < n; i++ {
        v := randIJK(r, -50, 50)
        // Parity
        got := v; got.Rotate60CCW(); want := o.Rotate60CCW(v)
        if got != want { t.Fatalf("Rotate60CCW parity failed for %v: %v vs %v", v, got, want) }
        got = v; got.Rotate60CW(); want = o.Rotate60CW(v)
        if got != want { t.Fatalf("Rotate60CW parity failed for %v: %v vs %v", v, got, want) }

        // Property: CCW(CW(v)) == v
        tmp := v; tmp.Rotate60CW(); tmp.Rotate60CCW()
        if tmp != v { t.Fatalf("Rotate property failed CW->CCW for %v: got %v", v, tmp) }

        // Property: 6×CCW returns same hex position
        h1 := v.ToHex2d()
        tmp = v
        for k := 0; k < 6; k++ { tmp.Rotate60CCW() }
        h2 := tmp.ToHex2d()
        if math.Abs(h1.X-h2.X) > 1e-12 || math.Abs(h1.Y-h2.Y) > 1e-12 {
            t.Fatalf("6x CCW rotation not identity in hex2d for %v: %v vs %v", v, h1, h2)
        }
    }
}

func TestOracle_CoordIJKNeighbor(t *testing.T) {
    o := newOracle(t)
    v := CoordIJK{2,1,-1}
    for d := Direction(0); d < NumDigits; d++ {
        got := v; got.Neighbor(d); want := o.Neighbor(v, d)
        if got != want { t.Fatalf("Neighbor(%d) = %v, want %v", d, got, want) }
    }

    // Randomized sweep over IJK and digits
    r := newRand()
    n := oracleMax()
    for i := 0; i < n; i++ {
        base := randIJK(r, -50, 50)
        for d := Direction(0); d < NumDigits; d++ {
            got := base; got.Neighbor(d); want := o.Neighbor(base, d)
            if got != want { t.Fatalf("Neighbor parity failed for %v dir %d: %v vs %v", base, d, got, want) }
            if d == CenterDigit {
                if got != base { t.Fatalf("Neighbor center digit changed coord for %v -> %v", base, got) }
            }
        }
    }
}

func TestOracle_CoordIJKDistance_Sweeps(t *testing.T) {
    o := newOracle(t)
    // Exhaustive small grid bounded by ORACLE_MAX
    maxPairs := oracleMax()
    count := 0
    for ia := -3; ia <= 3 && count < maxPairs; ia++ {
        for ja := -3; ja <= 3 && count < maxPairs; ja++ {
            for ka := -3; ka <= 3 && count < maxPairs; ka++ {
                a := CoordIJK{ia, ja, ka}
                for ib := -3; ib <= 3 && count < maxPairs; ib++ {
                    for jb := -3; jb <= 3 && count < maxPairs; jb++ {
                        for kb := -3; kb <= 3 && count < maxPairs; kb++ {
                            b := CoordIJK{ib, jb, kb}
                            if Distance(a, b) != o.Distance(a, b) {
                                t.Fatalf("Distance small-grid mismatch: %v,%v", a, b)
                            }
                            count++
                        }
                    }
                }
            }
        }
    }

    // Random sweep
    r := newRand()
    for i := 0; i < oracleMax(); i++ {
        a := randIJK(r, -50, 50)
        b := randIJK(r, -50, 50)
        if Distance(a, b) != o.Distance(a, b) {
            t.Fatalf("Distance random mismatch: %v,%v", a, b)
        }
    }
}
