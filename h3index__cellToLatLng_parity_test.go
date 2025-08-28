//go:build cgo

package h3

import (
	"fmt"
	"math"
	"testing"
)

func TestCellToLatLngParity(t *testing.T) {
	testCases := []H3Index{
		// Base resolution cells
		0x8001fffffffffff, // res 0, base cell 1
		0x8007fffffffffff, // res 0, base cell 7 (pentagon)
		0x800dfffffffffff, // res 0, base cell 13 (pentagon)

		// Higher resolution cells
		0x81283ffffffffff, // res 1, regular hexagon
		0x8228bffffffffff, // res 2, regular hexagon
		0x832830fffffffff, // res 3, regular hexagon

		// Pentagon cases
		0x81703ffffffffff, // res 1, pentagon base cell
		0x82734bfffffffff, // res 2, pentagon base cell

		// Known geographic locations
		0x891ea6d6533ffff, // San Francisco area
		0x891f0ab9207ffff, // New York area
		0x8928308280fffff, // London area
		0x89283082803ffff, // Sydney area

		// Edge cases around pentagon boundaries
		0x8427cffffffffff, // res 4, near pentagon
		0x85283473fffffff, // res 5, regular case

		// High resolution cases
		0x8a283082803ffff, // res 10
		0x8b283082800bfff, // res 11
	}

	const tolerance = 1e-10

	for i, h3 := range testCases {
		t.Run(fmt.Sprintf("case_%d_0x%x", i, uint64(h3)), func(t *testing.T) {
			// Test Go implementation
			var goLatLng LatLng
			goErr := cellToLatLng(h3, &goLatLng)

			// Test C implementation
			var cLatLng LatLng
			cErr := cellToLatLngC(h3, &cLatLng)

			// Compare errors
			if goErr != H3Error(cErr) {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
			}

			// If no error, compare coordinates
			if goErr == E_SUCCESS && cErr == 0 {
				if !goLatLng.Lat.EqualApprox(cLatLng.Lat, tolerance) {
					t.Errorf("Lat mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						goLatLng.Lat.Rad(), cLatLng.Lat.Rad(), math.Abs(goLatLng.Lat.Rad()-cLatLng.Lat.Rad()))
				}
				if !goLatLng.Lng.EqualApprox(cLatLng.Lng, tolerance) {
					t.Errorf("Lng mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						goLatLng.Lng.Rad(), cLatLng.Lng.Rad(), math.Abs(goLatLng.Lng.Rad()-cLatLng.Lng.Rad()))
				}
			}
		})
	}
}
