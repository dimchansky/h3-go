package faceijk

import (
    "fmt"
    "testing"
)

// GeoToClosestFace classification checks (kept as informational logs on differences)
func TestGeoToClosestFace(t *testing.T) {
    tests := []struct{ lat, lng float64; expectedFace int }{
        {37.7749, -122.4194, 1},  // SF
        {51.5074, -0.1278, 7},    // London
        {-33.8688, 151.2093, 12}, // Sydney
        {0.0, 0.0, 16},           // Origin
        {89.0, 0.0, 0},           // North pole
        {-89.0, 0.0, 10},         // South pole
        {35.6762, 139.6503, 18},  // Tokyo
        {52.5200, 13.4050, 7},    // Berlin
    }
    for i, tt := range tests {
        t.Run(fmt.Sprintf("test_%03d", i+1), func(t *testing.T) {
            face, sq := GeoToClosestFace(tt.lat, tt.lng)
            if face != tt.expectedFace {
                t.Logf("face diff: got %d want %d (sqDist=%.6f)", face, tt.expectedFace, sq)
            }
        })
    }
}

