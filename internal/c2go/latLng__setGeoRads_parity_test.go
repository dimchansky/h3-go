//go:build c2go

package c2go

import "testing"

func Test_setGeoRads_ParityWithC(t *testing.T) {
    p := LatLng{Lat: 0.0, Lng: 0.0}
    goP := p
    _setGeoRads(&goP, 1.23, -0.77)

    cP := _setGeoRadsC(p, 1.23, -0.77)

    if goP != cP {
        t.Fatalf("_setGeoRads mismatch: go=%+v c=%+v", goP, cP)
    }
}

