//go:build cgo && c2go

package h3

import (
	"fmt"
	"math"
	"testing"
)

func Test_vertexToLatLng_parity(t *testing.T) {
	tests := []struct {
		name   string
		vertex H3Index
	}{
		// Test vertices from various resolutions and cell types
		{"hex vertex r0", 0x20283080bffffff},
		{"hex vertex r1", 0x21283080bffffff},
		{"hex vertex r2", 0x22283080bffffff},
		{"hex vertex r5", 0x25283080bffffff},
		{"pent vertex r0", 0x20080800fffffff},
		{"pent vertex r1", 0x21080800fffffff},
		{"pent vertex r2", 0x22080800fffffff},
		{"pent vertex r5", 0x25080800fffffff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goCoord, cCoord LatLng

			// Get results from both implementations
			goErr := vertexToLatLng(tt.vertex, &goCoord)
			cErr := vertexToLatLngC(tt.vertex, &cCoord)

			// Compare error codes
			if goErr != cErr {
				t.Fatalf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If there was an error, skip coordinate comparison
			if goErr != E_SUCCESS {
				return
			}

			// Compare coordinates with appropriate tolerance for lat/lng
			const tolerance = 1e-12 // ~1 nanometer precision

			if !goCoord.Lat.EqualApprox(cCoord.Lat, tolerance) {
				t.Errorf("Latitude mismatch: Go=%.15f, C=%.15f, diff=%.15f",
					goCoord.Lat.Rad(), cCoord.Lat.Rad(), math.Abs(goCoord.Lat.Rad()-cCoord.Lat.Rad()))
			}

			if !goCoord.Lng.EqualApprox(cCoord.Lng, tolerance) {
				t.Errorf("Longitude mismatch: Go=%.15f, C=%.15f, diff=%.15f",
					goCoord.Lng.Rad(), cCoord.Lng.Rad(), math.Abs(goCoord.Lng.Rad()-cCoord.Lng.Rad()))
			}
		})
	}
}

func Test_vertexToLatLng_invalidVertex_parity(t *testing.T) {
	invalidVertices := []H3Index{
		H3_NULL,            // null vertex
		0x20283080bfffffff, // invalid reserved bits
		0x10283080bffffff,  // wrong mode (cell mode)
		0x30283080bffffff,  // wrong mode (edge mode)
	}

	for i, vertex := range invalidVertices {
		t.Run(fmt.Sprintf("invalid_%d", i), func(t *testing.T) {
			var goCoord, cCoord LatLng

			goErr := vertexToLatLng(vertex, &goCoord)
			cErr := vertexToLatLngC(vertex, &cCoord)

			if goErr != cErr {
				t.Fatalf("Error mismatch for invalid vertex 0x%x: Go=%v, C=%v",
					vertex, goErr, cErr)
			}
		})
	}
}
