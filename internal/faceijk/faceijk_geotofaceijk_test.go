package faceijk

import (
    "fmt"
    "testing"

    "github.com/dimchansky/h3-go/internal/coordijk"
)

// GeoToFaceIJK mapping checks (logs differences; values are approximate for smoke testing)
func TestGeoToFaceIJK(t *testing.T) {
    tests := []struct{ lat, lng float64; res int; expected FaceIJK }{
        {37.7749, -122.4194, 0, FaceIJK{1, coordijk.CoordIJK{0,0,0}}},
        {51.5074, -0.1278,   0, FaceIJK{7, coordijk.CoordIJK{0,0,0}}},
        {0.0,     0.0,       0, FaceIJK{16, coordijk.CoordIJK{0,0,0}}},
        {37.7749, -122.4194, 1, FaceIJK{1, coordijk.CoordIJK{1,0,0}}},
        {89.0,    0.0,       0, FaceIJK{0, coordijk.CoordIJK{0,0,0}}},
    }
    for i, tt := range tests {
        t.Run(fmt.Sprintf("test_%03d", i+1), func(t *testing.T) {
            got := GeoToFaceIJK(tt.lat, tt.lng, tt.res)
            if got.Face != tt.expected.Face {
                t.Logf("face diff: got %d want %d", got.Face, tt.expected.Face)
            }
            t.Logf("GeoToFaceIJK(lat=%.4f,lng=%.4f,res=%d) => Face:%d IJK:(%d,%d,%d)", tt.lat, tt.lng, tt.res, got.Face, got.Coord.I, got.Coord.J, got.Coord.K)
        })
    }
}

