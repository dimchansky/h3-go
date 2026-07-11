package h3

import (
	"math"
	"testing"
)

func FuzzParseCell(f *testing.F) {
	f.Add("8928308280fffff")
	f.Add("0x8f2830828052d25")
	f.Add("")
	f.Add("zzz")
	f.Add("ffffffffffffffff")
	f.Fuzz(func(t *testing.T, s string) {
		c, err := ParseCell(s)
		if err != nil {
			if c != 0 {
				t.Fatalf("ParseCell(%q) returned %v with error %v", s, c, err)
			}
			return
		}
		if !c.IsValid() {
			t.Fatalf("ParseCell(%q) = %v accepted an invalid cell", s, c)
		}
		// String/Parse round trip on the canonical form.
		back, err := ParseCell(c.String())
		if err != nil || back != c {
			t.Fatalf("round trip %q -> %v -> %v (%v)", s, c, back, err)
		}
	})
}

func FuzzLatLngToCell(f *testing.F) {
	f.Add(37.775938728915946, -122.41795063018799, 9)
	f.Add(0.0, 0.0, 0)
	f.Add(90.0, 180.0, 15)
	f.Add(-90.0, -180.0, 15)
	f.Add(math.Inf(1), 0.0, 5)
	f.Fuzz(func(t *testing.T, latDegs, lngDegs float64, res int) {
		c, err := LatLngToCell(LatLngDegs(latDegs, lngDegs), res)
		if err != nil {
			return
		}
		if !c.IsValid() {
			t.Fatalf("LatLngToCell(%v, %v, %d) = %v is invalid", latDegs, lngDegs, res, c)
		}
		if c.Resolution() != res {
			t.Fatalf("resolution mismatch: %d != %d", c.Resolution(), res)
		}
		// The center of the produced cell must map back to the same cell.
		center, err := c.LatLng()
		if err != nil {
			t.Fatal(err)
		}
		back, err := LatLngToCell(center, res)
		if err != nil || back != c {
			t.Fatalf("center round trip: %v -> %v (%v)", c, back, err)
		}
	})
}

func FuzzHierarchyRoundTrip(f *testing.F) {
	f.Add(uint64(0x8928308280fffff), 5)
	f.Add(uint64(0x8f2830828052d25), 0)
	f.Add(uint64(0), 7)
	f.Add(^uint64(0), 3)
	f.Fuzz(func(t *testing.T, raw uint64, parentRes int) {
		c := Cell(raw)
		pos, err := c.ChildPos(parentRes)
		if err != nil {
			return // invalid input rejected: fine
		}
		parent, err := c.Parent(parentRes)
		if err != nil {
			t.Fatalf("ChildPos succeeded but Parent failed: %v", err)
		}
		back, err := parent.ChildAtPos(pos, c.Resolution())
		if err != nil || back != c {
			t.Fatalf("ChildAtPos(ChildPos(%v)) = %v (%v)", c, back, err)
		}
	})
}
