//go:build c2go

package c2go

import (
    "math"
    "testing"
)

func Test_geoToVec3d_ParityWithC(t *testing.T) {
    cases := []LatLng{
        {Lat: 0, Lng: 0},
        {Lat: math.Pi / 6, Lng: math.Pi / 4},
        {Lat: -math.Pi / 3, Lng: -math.Pi / 2},
        {Lat: 1.2, Lng: -2.3},
    }
    for _, g := range cases {
        goV := _geoToVec3d(g)
        cV := geoToVec3dC(g)
        if math.Abs(goV.X-cV.X) > 1e-15 || math.Abs(goV.Y-cV.Y) > 1e-15 || math.Abs(goV.Z-cV.Z) > 1e-15 {
            t.Fatalf("_geoToVec3d mismatch for geo=%+v: go=(%.17g,%.17g,%.17g) c=(%.17g,%.17g,%.17g)", g, goV.X, goV.Y, goV.Z, cV.X, cV.Y, cV.Z)
        }
    }
}

