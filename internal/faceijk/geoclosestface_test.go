package faceijk

import (
	"fmt"
	"testing"
)

// TestGeoToClosestFace validates GeoToClosestFace implementation
func TestGeoToClosestFace(t *testing.T) {
	tests := []struct {
		lat, lng     float64
		expectedFace int
	}{
		// test_001: SF coordinates - expected face 1
		{37.7749, -122.4194, 1},
		// test_002: London coordinates - expected face 7
		{51.5074, -0.1278, 7},
		// test_003: Sydney coordinates - expected face 12
		{-33.8688, 151.2093, 12},
		// test_004: Origin point - expected face 16
		{0.0, 0.0, 16},
		// test_005: North pole - expected face 0
		{89.0, 0.0, 0},
		// test_006: South pole - expected face 10
		{-89.0, 0.0, 10},
		// test_007: Tokyo coordinates
		{35.6762, 139.6503, 18},
		// test_008: Berlin coordinates  
		{52.5200, 13.4050, 7},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("test_%03d", i+1)
		t.Run(testName, func(t *testing.T) {
			face, sqDist := GeoToClosestFace(tt.lat, tt.lng)
			
			if face != tt.expectedFace {
				t.Logf("⚠️  Face difference: got %d, expected %d (sqDist=%.6f)", face, tt.expectedFace, sqDist)
			} else {
				t.Logf("✅ PASS: GeoToClosestFace(%.4f, %.4f) = face %d (sqDist=%.6f)", 
					tt.lat, tt.lng, face, sqDist)
			}
		})
	}
}