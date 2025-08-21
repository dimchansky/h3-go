package h3

import "testing"

// TestLatLngToCell_Table verifies known conversions against values
// produced by the C oracle (testref/h3ref). These are fixed cases so
// they run without requiring the oracle at test time.
func TestLatLngToCell_Table(t *testing.T) {
	tests := []struct {
		lat, lng float64
		res      int
		want     uint64 // expected H3 index
	}{
		{37.7749, -122.4194, 9, 0x89283082803ffff},  // San Francisco
		{51.5074, -0.1278, 10, 0x8a195da49a2ffff},   // London
		{-33.8688, 151.2093, 12, 0x8cbe0e35cbad1ff}, // Sydney
		{0.0, 0.0, 5, 0x85754e67fffffff},
		{89.0, 0.0, 2, 0x820327fffffffff},           // near North Pole
		{-89.0, 0.0, 2, 0x82f167fffffffff},          // near South Pole
		{35.6762, 139.6503, 9, 0x892f5a363bbffff},   // Tokyo
		{52.5200, 13.4050, 9, 0x891f1d48947ffff},    // Berlin
		{40.7128, -74.0060, 9, 0x892a1072893ffff},   // New York
		{-22.9068, -43.1729, 8, 0x88a8a06a0dfffff},  // Rio de Janeiro
		{30.0444, 31.2357, 7, 0x873e628e6ffffff},    // Cairo
		{21.3069, -157.8583, 10, 0x8a464b96ab2ffff}, // Honolulu
	}

	for _, tt := range tests {
		got, err := LatLngToCell(LatLng{Lat: tt.lat, Lng: tt.lng}, tt.res)
		if err != nil {
			t.Fatalf("LatLngToCell(%+.6f,%+.6f,%d) unexpected error: %v", tt.lat, tt.lng, tt.res, err)
		}
		if uint64(got) != tt.want {
			t.Fatalf("LatLngToCell(%+.6f,%+.6f,%d): got 0x%x want 0x%x", tt.lat, tt.lng, tt.res, uint64(got), tt.want)
		}
	}
}
