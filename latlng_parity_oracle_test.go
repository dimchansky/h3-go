//go:build oracle

package h3

import (
    "testing"

    testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

// Parity test against the C oracle for LatLngToCell.
func TestOracle_LatLngToCell_Parity(t *testing.T) {
    o := testoracle.New(t)
    r := testoracle.NewRand()
    for i := 0; i < testoracle.Max(); i++ {
        lat := r.Float64()*180 - 90   // [-90, 90]
        lng := r.Float64()*360 - 180  // [-180, 180)
        for res := 0; res <= MaxResolution; res++ { // full resolution range [0..15]
            got, err := LatLngToCell(LatLng{Lat: lat, Lng: lng}, res)
            if err != nil {
                // Ask oracle what it thinks
                _, code := o.H3FromLatLng(lat, lng, res)
                if code == 0 {
                    t.Fatalf("Go returned error %v but oracle succeeded for latlng(%.6f,%.6f,%d)", err, lat, lng, res)
                }
                continue
            }
            want, code := o.H3FromLatLng(lat, lng, res)
            if code != 0 {
                t.Fatalf("Oracle error code %d for latlng(%.6f,%.6f,%d)", code, lat, lng, res)
            }
            if uint64(got) != want {
                t.Fatalf("LatLngToCell parity: got 0x%x, want 0x%x (lat=%.6f,lng=%.6f,res=%d)", uint64(got), want, lat, lng, res)
            }
        }
    }
}
