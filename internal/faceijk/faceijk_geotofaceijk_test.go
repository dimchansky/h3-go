package faceijk

import (
	"fmt"
	"testing"

	"github.com/dimchansky/h3-go/internal/angles"
	"github.com/dimchansky/h3-go/internal/coordijk"
)

// GeoToFaceIJK mapping checks (logs differences; values are approximate for smoke testing).
func TestGeoToFaceIJK(t *testing.T) {
	tests := []struct {
		lat, lng float64
		res      int
		expected FaceIJK
	}{
		{37.7749, -122.4194, 0, FaceIJK{Face: 7, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 1}}},
		{51.5074, -0.1278, 0, FaceIJK{Face: 3, Coord: coordijk.CoordIJK{I: 1, J: 0, K: 0}}},
		{0.0, 0.0, 0, FaceIJK{Face: 9, Coord: coordijk.CoordIJK{I: 0, J: 0, K: 2}}},
		{37.7749, -122.4194, 1, FaceIJK{Face: 7, Coord: coordijk.CoordIJK{I: 0, J: 1, K: 3}}},
		{89.0, 0.0, 0, FaceIJK{Face: 1, Coord: coordijk.CoordIJK{I: 1, J: 0, K: 0}}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%03d", i+1), func(t *testing.T) {
			got := GeoToFaceIJK(angles.DegreesToRadians(tt.lat), angles.DegreesToRadians(tt.lng), tt.res)
			if got.Face != tt.expected.Face || got.Coord != tt.expected.Coord {
				t.Errorf("GeoToFaceIJK mismatch lat=%.4f,lng=%.4f,res=%d: got {face=%d ijk=(%d,%d,%d)} want {face=%d ijk=(%d,%d,%d)}",
					tt.lat, tt.lng, tt.res, got.Face, got.Coord.I, got.Coord.J, got.Coord.K, tt.expected.Face, tt.expected.Coord.I, tt.expected.Coord.J, tt.expected.Coord.K)
			}
		})
	}
}
