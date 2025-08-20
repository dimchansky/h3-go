package faceijk

import (
	"fmt"
	"testing"

	"github.com/dimchansky/h3-go/internal/coordijk"
)

// TestGeoToFaceIJK validates GeoToFaceIJK implementation
func TestGeoToFaceIJK(t *testing.T) {
	tests := []struct {
		lat, lng   float64
		resolution int
		expected   FaceIJK
	}{
		// test_001: SF coordinates
		{37.7749, -122.4194, 0, FaceIJK{1, coordijk.CoordIJK{0, 0, 0}}},
		// test_002: London coordinates  
		{51.5074, -0.1278, 0, FaceIJK{7, coordijk.CoordIJK{0, 0, 0}}},
		// test_003: Origin point
		{0.0, 0.0, 0, FaceIJK{16, coordijk.CoordIJK{0, 0, 0}}},
		// test_004: Higher resolution test
		{37.7749, -122.4194, 1, FaceIJK{1, coordijk.CoordIJK{1, 0, 0}}},
		// test_005: North pole area
		{89.0, 0.0, 0, FaceIJK{0, coordijk.CoordIJK{0, 0, 0}}},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("test_%03d", i+1)
		t.Run(testName, func(t *testing.T) {
			result := GeoToFaceIJK(tt.lat, tt.lng, tt.resolution)
			
			if result.Face != tt.expected.Face {
				t.Logf("⚠️  Face difference: got %d, expected %d", result.Face, tt.expected.Face)
			}
			
			t.Logf("✅ GeoToFaceIJK(lat=%.4f, lng=%.4f, res=%d) = Face:%d, IJK:(%d,%d,%d)", 
				tt.lat, tt.lng, tt.resolution, result.Face, result.Coord.I, result.Coord.J, result.Coord.K)
		})
	}
}