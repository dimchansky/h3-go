//go:build oracle

package faceijk

import (
    "testing"

    "github.com/dimchansky/h3-go/internal/angles"
    testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

func TestOracle_GeoToFaceIJK_Parity(t *testing.T) {
    o := testoracle.New(t)
    r := testoracle.NewRand()
    for i := 0; i < testoracle.Max(); i++ {
        lat := r.Float64()*180 - 90
        lng := r.Float64()*360 - 180
        res := int(r.Intn(4))
        // our
        got := GeoToFaceIJK(angles.DegreesToRadians(lat), angles.DegreesToRadians(lng), res)
        // oracle
        face, I, J, K := o.GeoToFaceIJK(lat, lng, res)
        if got.Face != face || got.Coord.I != I || got.Coord.J != J || got.Coord.K != K {
            t.Fatalf("GeoToFaceIJK mismatch lat=%.6f lng=%.6f res=%d: got {%d %d %d %d} want {%d %d %d %d}",
                lat, lng, res, got.Face, got.Coord.I, got.Coord.J, got.Coord.K, face, I, J, K)
        }
    }
}
