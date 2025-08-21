//go:build oracle

package testoracle

import (
	"math"
	"testing"

	"github.com/dimchansky/h3-go/internal/indexbits"
)

func TestLatLngToCell(t *testing.T) {
	c := New(t)
	h, code := c.LatLngToCell(37.7749, -122.4194, 9)
	if code != 0 {
		t.Fatalf("LatLngToCell returned error code %d", code)
	}
	if h == 0 {
		t.Fatalf("LatLngToCell returned zero index")
	}
}

func TestFaceIjkToH3(t *testing.T) {
	c := New(t)
	h := c.FaceIjkToH3(0, 0, 0, 0, 0)
	if h == 0 {
		t.Fatalf("FaceIjkToH3 returned zero index")
	}
}

func TestRotateH3(t *testing.T) {
	c := New(t)
	h, code := c.LatLngToCell(51.5074, -0.1278, 10)
	if code != 0 {
		t.Fatalf("LatLngToCell code=%d", code)
	}
	ccw := c.H3Rotate60ccw(h)
	back := c.H3Rotate60cw(ccw)
	if back == 0 {
		t.Fatalf("H3Rotate60cw returned zero")
	}
}

func TestGetResolutionAndBaseCell(t *testing.T) {
	c := New(t)
	h, code := c.LatLngToCell(0, 0, 5)
	if code != 0 {
		t.Fatalf("LatLngToCell code=%d", code)
	}
	if got := c.GetResolution(h); got != 5 {
		t.Fatalf("GetResolution: got %d want 5", got)
	}
	wantBase := indexbits.GetBaseCell(h)
	if gotBase := c.GetBaseCellNumber(h); gotBase != wantBase {
		t.Fatalf("GetBaseCellNumber: got %d want %d", gotBase, wantBase)
	}
}

func TestIsBaseCellPentagon(t *testing.T) {
	c := New(t)
	if !c.IsBaseCellPentagon(4) {
		t.Fatalf("base cell 4 should be pentagon")
	}
	if c.IsBaseCellPentagon(0) {
		t.Fatalf("base cell 0 should be hexagon")
	}
}

func TestIJKDistance(t *testing.T) {
	c := New(t)
	d := c.IJKDistance([3]int{0, 0, 0}, [3]int{1, 0, 1})
	if d != 1 {
		t.Fatalf("IJKDistance: got %d want 1", d)
	}
}

func TestIJKRotate60(t *testing.T) {
	c := New(t)
	got := c.IJKRotate60ccw([3]int{1, 0, 0})
	if got != ([3]int{1, 1, 0}) {
		t.Fatalf("RotateIJKCCW: got %v want %v", got, [3]int{1, 1, 0})
	}
	got2 := c.IJKRotate60cw([3]int{0, 1, 0})
	if got2 != ([3]int{1, 1, 0}) {
		t.Fatalf("RotateIJKCW: got %v want %v", got2, [3]int{1, 1, 0})
	}
}

func TestIJKHex2d(t *testing.T) {
	c := New(t)
	x, y := c.IJKToHex2d([3]int{1, 1, 0})
	if math.Abs(x-0.5) > 1e-12 || math.Abs(y-math.Sqrt(3)/2) > 1e-12 {
		t.Fatalf("IJKToHex2D unexpected: (%.15f,%.15f)", x, y)
	}
	ijk := c.Hex2dToCoordIJK(x, y)
	if ijk != ([3]int{1, 1, 0}) {
		t.Fatalf("Hex2DToIJK roundtrip: got %v want %v", ijk, [3]int{1, 1, 0})
	}
}

func TestNeighbor(t *testing.T) {
	c := New(t)
	got := c.Neighbor([3]int{0, 0, 0}, 4)
	if got != ([3]int{1, 0, 0}) {
		t.Fatalf("Neighbor: got %v want %v", got, [3]int{1, 0, 0})
	}
}

func TestApertureTransforms(t *testing.T) {
	c := New(t)
	// Just smoke-test commands return values
	_ = c.UpAp7([3]int{1, 0, 0})
	_ = c.UpAp7r([3]int{1, 0, 0})
	_ = c.DownAp7([3]int{1, 0, 0})
	_ = c.DownAp7r([3]int{1, 0, 0})
	_ = c.DownAp3([3]int{1, 0, 0})
	_ = c.DownAp3r([3]int{1, 0, 0})
}

func TestGeoToFaceIjk(t *testing.T) {
	c := New(t)
	face, i, j, k := c.GeoToFaceIjk(37.7749, -122.4194, 0)
	if face < 0 || face > 19 {
		t.Fatalf("face out of range: %d", face)
	}
	if i < 0 || j < 0 || k < 0 {
		t.Fatalf("negative IJK: %d %d %d", i, j, k)
	}
}

func TestGeoToHex2d(t *testing.T) {
	c := New(t)
	face, x, y := c.GeoToHex2d(0, 0, 0)
	if face < 0 || face > 19 {
		t.Fatalf("face out of range: %d", face)
	}
	if math.IsNaN(x) || math.IsNaN(y) {
		t.Fatalf("geoToHex2d returned NaN")
	}
}
