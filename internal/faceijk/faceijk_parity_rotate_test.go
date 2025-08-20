//go:build oracle

package faceijk

import (
    "testing"

    testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

func TestOracle_H3Rotate_Parity(t *testing.T) {
    o := testoracle.New(t)
    r := testoracle.NewRand()
    // Randomized H3 indices via oracle latlng
    for i := 0; i < testoracle.Max(); i++ {
        lat := r.Float64()*180 - 90
        lng := r.Float64()*360 - 180
        res := int(r.Intn(4))
        h3, code := o.H3FromLatLng(lat, lng, res)
        if code != 0 {
            continue
        }
        // CCW parity
        gotCCW := h3Rotate60ccw(h3)
        wantCCW := o.RotateH3CCW(h3)
        if gotCCW != wantCCW {
            t.Fatalf("rotate60ccw parity failed for 0x%x: got 0x%x want 0x%x", h3, gotCCW, wantCCW)
        }
        // CW parity
        gotCW := h3Rotate60cw(h3)
        wantCW := o.RotateH3CW(h3)
        if gotCW != wantCW {
            t.Fatalf("rotate60cw parity failed for 0x%x: got 0x%x want 0x%x", h3, gotCW, wantCW)
        }
    }
}

func TestOracle_H3Rotate_PentagonCases(t *testing.T) {
    o := testoracle.New(t)
    // Use known pentagon base cells to generate H3 at res 0
    pentFaces := []struct{ face, i, j, k int }{
        {0, 2, 0, 0},  // base cell 4
        {11, 2, 0, 0}, // base cell 14
    }
    for _, p := range pentFaces {
        h := o.H3FromFaceIJK(p.face, p.i, p.j, p.k, 0)
        if h == 0 { t.Fatalf("oracle pentagon faceijk produced 0") }
        if got := h3Rotate60ccw(h); got != o.RotateH3CCW(h) {
            t.Fatalf("pent ccw mismatch for 0x%x", h)
        }
        if got := h3Rotate60cw(h); got != o.RotateH3CW(h) {
            t.Fatalf("pent cw mismatch for 0x%x", h)
        }
    }
}

