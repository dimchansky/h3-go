//go:build oracle

package coordijk

import (
	"math"
	"testing"

	testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

func TestOracle_CoordIJKDistance(t *testing.T) {
	o := testoracle.New(t)
	cases := []struct{ a, b CoordIJK }{
		{CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{CoordIJK{1, 0, 0}, CoordIJK{0, 0, 0}},
		{CoordIJK{0, 2, 0}, CoordIJK{0, 0, 0}},
		{CoordIJK{1, 0, 1}, CoordIJK{0, 0, 0}},
		{CoordIJK{1, 0, 1}, CoordIJK{0, 2, 0}},
		{CoordIJK{2, -1, 1}, CoordIJK{-1, 2, -1}},
	}
	for _, tc := range cases {
		if Distance(tc.a, tc.b) != o.IJKDistance([3]int{tc.a.I, tc.a.J, tc.a.K}, [3]int{tc.b.I, tc.b.J, tc.b.K}) {
			t.Fatalf("Distance(%v,%v) mismatch", tc.a, tc.b)
		}
	}
}

func TestOracle_CoordIJKRotate(t *testing.T) {
	o := testoracle.New(t)
	cases := []CoordIJK{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0}, {1, 0, 1}, {0, 1, 1}, {2, -1, 1}, {-1, 2, -1}}
	for _, v := range cases {
		got := v
		got.Rotate60CCW()
		arr := o.IJKRotate60ccw([3]int{v.I, v.J, v.K})
		want := CoordIJK{arr[0], arr[1], arr[2]}
		if got != want {
			t.Fatalf("Rotate60CCW(%v) = %v, want %v", v, got, want)
		}
		got = v
		got.Rotate60CW()
		arr = o.IJKRotate60cw([3]int{v.I, v.J, v.K})
		want = CoordIJK{arr[0], arr[1], arr[2]}
		if got != want {
			t.Fatalf("Rotate60CW(%v) = %v, want %v", v, got, want)
		}
	}

	// Randomized parity and properties
	r := testoracle.NewRand()
	n := testoracle.Max()
	for i := 0; i < n; i++ {
		vArr := testoracle.RandIJK(r, -50, 50)
		v := CoordIJK{vArr[0], vArr[1], vArr[2]}
		// Parity
		got := v
		got.Rotate60CCW()
		arr := o.IJKRotate60ccw([3]int{v.I, v.J, v.K})
		want := CoordIJK{arr[0], arr[1], arr[2]}
		if got != want {
			t.Fatalf("Rotate60CCW parity failed for %v: %v vs %v", v, got, want)
		}
		got = v
		got.Rotate60CW()
		arr = o.IJKRotate60cw([3]int{v.I, v.J, v.K})
		want = CoordIJK{arr[0], arr[1], arr[2]}
		if got != want {
			t.Fatalf("Rotate60CW parity failed for %v: %v vs %v", v, got, want)
		}

		// Property: CCW(CW(v)) preserves hex position (compare in hex2d)
		tmp := v
		tmp.Rotate60CW()
		tmp.Rotate60CCW()
		h0, hT := v.ToHex2d(), tmp.ToHex2d()
		if math.Abs(h0.X-hT.X) > 1e-12 || math.Abs(h0.Y-hT.Y) > 1e-12 {
			t.Fatalf("Rotate property failed CW->CCW for %v: hex2d %v vs %v", v, h0, hT)
		}

		// Property: 6×CCW returns same hex position
		h1 := v.ToHex2d()
		tmp = v
		for k := 0; k < 6; k++ {
			tmp.Rotate60CCW()
		}
		h2 := tmp.ToHex2d()
		if math.Abs(h1.X-h2.X) > 1e-12 || math.Abs(h1.Y-h2.Y) > 1e-12 {
			t.Fatalf("6x CCW rotation not identity in hex2d for %v: %v vs %v", v, h1, h2)
		}
	}
}

func TestOracle_CoordIJKNeighbor(t *testing.T) {
	o := testoracle.New(t)
	v := CoordIJK{2, 1, -1}
	for d := Direction(0); d < NumDigits; d++ {
		got := v
		got.Neighbor(d)
		arr := o.Neighbor([3]int{v.I, v.J, v.K}, int(d))
		want := CoordIJK{arr[0], arr[1], arr[2]}
		if got != want {
			t.Fatalf("Neighbor(%d) = %v, want %v", d, got, want)
		}
	}

	// Randomized sweep over IJK and digits
	r := testoracle.NewRand()
	n := testoracle.Max()
	for i := 0; i < n; i++ {
		bArr := testoracle.RandIJK(r, -50, 50)
		base := CoordIJK{bArr[0], bArr[1], bArr[2]}
		for d := Direction(0); d < NumDigits; d++ {
			got := base
			got.Neighbor(d)
			arr := o.Neighbor([3]int{base.I, base.J, base.K}, int(d))
			want := CoordIJK{arr[0], arr[1], arr[2]}
			if got != want {
				t.Fatalf("Neighbor parity failed for %v dir %d: %v vs %v", base, d, got, want)
			}
			if d == CenterDigit {
				if got != base {
					t.Fatalf("Neighbor center digit changed coord for %v -> %v", base, got)
				}
			}
		}
	}
}

func TestOracle_CoordIJKDistance_Sweeps(t *testing.T) {
	o := testoracle.New(t)
	// Exhaustive small grid bounded by ORACLE_MAX
	maxPairs := testoracle.Max()
	count := 0
	for ia := -3; ia <= 3 && count < maxPairs; ia++ {
		for ja := -3; ja <= 3 && count < maxPairs; ja++ {
			for ka := -3; ka <= 3 && count < maxPairs; ka++ {
				a := CoordIJK{ia, ja, ka}
				for ib := -3; ib <= 3 && count < maxPairs; ib++ {
					for jb := -3; jb <= 3 && count < maxPairs; jb++ {
						for kb := -3; kb <= 3 && count < maxPairs; kb++ {
							b := CoordIJK{ib, jb, kb}
							if Distance(a, b) != o.IJKDistance([3]int{a.I, a.J, a.K}, [3]int{b.I, b.J, b.K}) {
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
	r := testoracle.NewRand()
	for i := 0; i < testoracle.Max(); i++ {
		aArr := testoracle.RandIJK(r, -50, 50)
		bArr := testoracle.RandIJK(r, -50, 50)
		a := CoordIJK{aArr[0], aArr[1], aArr[2]}
		b := CoordIJK{bArr[0], bArr[1], bArr[2]}
		if Distance(a, b) != o.IJKDistance([3]int{a.I, a.J, a.K}, [3]int{b.I, b.J, b.K}) {
			t.Fatalf("Distance random mismatch: %v,%v", a, b)
		}
	}
}
