package faceijk

import (
    "fmt"
    "testing"

    "github.com/dimchansky/h3-go/internal/angles"
)

// GeoToClosestFace classification checks (kept as informational logs on differences)
func TestGeoToClosestFace(t *testing.T) {
    tests := []struct{ lat, lng float64; expectedFace int }{
        {37.7749, -122.4194, 7},  // SF
        {51.5074, -0.1278, 3},    // London
        {-33.8688, 151.2093, 15}, // Sydney
        {0.0, 0.0, 9},            // Origin
        {89.0, 0.0, 1},           // North pole
        {-89.0, 0.0, 18},         // South pole
        {35.6762, 139.6503, 6},   // Tokyo
        {52.5200, 13.4050, 4},    // Berlin
    }
    for i, tt := range tests {
        t.Run(fmt.Sprintf("test_%03d", i+1), func(t *testing.T) {
            face, sq := GeoToClosestFace(angles.DegreesToRadians(tt.lat), angles.DegreesToRadians(tt.lng))
            if face != tt.expectedFace {
                t.Errorf("GeoToClosestFace mismatch lat=%.4f,lng=%.4f: got %d want %d (sqDist=%.6f)", tt.lat, tt.lng, face, tt.expectedFace, sq)
            }
        })
    }
}
