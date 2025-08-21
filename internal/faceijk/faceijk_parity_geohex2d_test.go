//go:build oracle

package faceijk

import (
    "math"
    "testing"

    "github.com/dimchansky/h3-go/internal/angles"
    testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

func TestOracle_GeoToHex2d_Parity(t *testing.T) {
    o := testoracle.New(t)
    r := testoracle.NewRand()
    for i := 0; i < testoracle.Max(); i++ {
        lat := r.Float64()*180 - 90
        lng := r.Float64()*360 - 180
        res := int(r.Intn(4))
        face, v := GeoToHex2d(angles.DegreesToRadians(lat), angles.DegreesToRadians(lng), res)
        oface, ox, oy := o.GeoToHex2d(lat, lng, res)
        if face != oface || math.Abs(v.X-ox) > 1e-12 || math.Abs(v.Y-oy) > 1e-12 {
            t.Fatalf("GeoToHex2d mismatch lat=%.6f lng=%.6f res=%d: got {face=%d (%.17g,%.17g)} want {face=%d (%.17g,%.17g)}",
                lat, lng, res, face, v.X, v.Y, oface, ox, oy)
        }
    }
}
