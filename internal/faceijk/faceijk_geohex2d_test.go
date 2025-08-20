package faceijk

import (
    "math"
    "testing"
)

func TestGeoToHex2d(t *testing.T) {
    tests := []struct{
        name string
        lat, lng float64
        res int
        expectValidFace bool
    }{
        {"face0_res0", 0.8, 1.2, 0, true},
        {"face1_res1", 1.3, 2.5, 1, true},
        {"face4_res2", 0.5, 0.4, 2, true},
        {"equator_res0", 0.0, 0.0, 0, true},
        {"north_pole_res1", 1.57079632679, 0.0, 1, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            face, v := GeoToHex2d(tt.lat, tt.lng, tt.res)
            if tt.expectValidFace {
                if face < 0 || face >= NumIcosaFaces {
                    t.Errorf("GeoToHex2d(%f,%f,%d) face=%d out of range", tt.lat, tt.lng, tt.res, face)
                }
                if math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) {
                    t.Errorf("GeoToHex2d produced invalid hex2d: %v", v)
                }
            }
        })
    }
}

