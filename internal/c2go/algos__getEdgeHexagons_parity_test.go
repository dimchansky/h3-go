//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_getEdgeHexagons_parity(t *testing.T) {
	tests := []struct {
		name        string
		geoloop     []LatLng
		numHexagons int64
		res         int32
		description string
	}{
		{
			name:        "simple triangle",
			geoloop:     []LatLng{{0.1, 0.1}, {0.1, 0.2}, {0.2, 0.15}},
			numHexagons: 100,
			res:         9,
			description: "Small triangle at low resolution",
		},
		{
			name:        "square loop",
			geoloop:     []LatLng{{0.0, 0.0}, {0.0, 0.1}, {0.1, 0.1}, {0.1, 0.0}},
			numHexagons: 200,
			res:         8,
			description: "Square polygon",
		},
		{
			name:        "single point loop",
			geoloop:     []LatLng{{0.5, -0.5}},
			numHexagons: 50,
			res:         7,
			description: "Degenerate single point loop",
		},
		{
			name:        "pentagon around center",
			geoloop:     createPentagonLoop(0.0, 0.0, 0.05),
			numHexagons: 300,
			res:         10,
			description: "Pentagon shape at higher resolution",
		},
		{
			name:        "large hexagon at low res",
			geoloop:     createHexagonLoop(0.3, -0.3, 0.2),
			numHexagons: 150,
			res:         6,
			description: "Larger hexagon at lower resolution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize arrays for both implementations
			goSearch := make([]H3Index, tt.numHexagons)
			goFound := make([]H3Index, tt.numHexagons)
			goNumSearchHexes := int64(0)

			cSearch := make([]H3Index, tt.numHexagons)
			cFound := make([]H3Index, tt.numHexagons)
			cNumSearchHexes := int64(0)

			// Call Go implementation
			goErr := _getEdgeHexagons(tt.geoloop, tt.numHexagons, tt.res, &goNumSearchHexes, goSearch, goFound)

			// Call C implementation
			cErr := _getEdgeHexagonsC(tt.geoloop, tt.numHexagons, tt.res, &cNumSearchHexes, cSearch, cFound)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// Compare number of search hexagons found
			if goNumSearchHexes != cNumSearchHexes {
				t.Errorf("numSearchHexes mismatch: Go=%d, C=%d", goNumSearchHexes, cNumSearchHexes)
				return
			}

			// Compare search arrays up to numSearchHexes
			goSearchSet := make(map[H3Index]bool)
			cSearchSet := make(map[H3Index]bool)

			for i := int64(0); i < goNumSearchHexes; i++ {
				goSearchSet[goSearch[i]] = true
			}
			for i := int64(0); i < cNumSearchHexes; i++ {
				cSearchSet[cSearch[i]] = true
			}

			if len(goSearchSet) != len(cSearchSet) {
				t.Errorf("Different number of unique search hexagons: Go=%d, C=%d", len(goSearchSet), len(cSearchSet))
				return
			}

			for hex := range goSearchSet {
				if !cSearchSet[hex] {
					t.Errorf("Go search contains hex %d not found in C search", hex)
					return
				}
			}

			// Compare found arrays (hash table structure)
			goFoundSet := make(map[H3Index]bool)
			cFoundSet := make(map[H3Index]bool)

			for _, hex := range goFound {
				if hex != 0 {
					goFoundSet[hex] = true
				}
			}
			for _, hex := range cFound {
				if hex != 0 {
					cFoundSet[hex] = true
				}
			}

			if len(goFoundSet) != len(cFoundSet) {
				t.Errorf("Different number of found hexagons: Go=%d, C=%d", len(goFoundSet), len(cFoundSet))
				return
			}

			for hex := range goFoundSet {
				if !cFoundSet[hex] {
					t.Errorf("Go found contains hex %d not found in C found", hex)
					return
				}
			}
		})
	}
}

// createPentagonLoop creates a pentagon-shaped loop of coordinates
func createPentagonLoop(centerLat, centerLng, radius float64) []LatLng {
	loop := make([]LatLng, 5)
	for i := 0; i < 5; i++ {
		angle := 2.0 * math.Pi * float64(i) / 5.0
		loop[i] = LatLng{
			Lat: centerLat + radius*math.Cos(angle),
			Lng: centerLng + radius*math.Sin(angle),
		}
	}
	return loop
}

// createHexagonLoop creates a hexagon-shaped loop of coordinates
func createHexagonLoop(centerLat, centerLng, radius float64) []LatLng {
	loop := make([]LatLng, 6)
	for i := 0; i < 6; i++ {
		angle := 2.0 * math.Pi * float64(i) / 6.0
		loop[i] = LatLng{
			Lat: centerLat + radius*math.Cos(angle),
			Lng: centerLng + radius*math.Sin(angle),
		}
	}
	return loop
}
