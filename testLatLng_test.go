// Tests ported from testLatLng.c

package h3

import (
	"math"
	"testing"
)

// testDecreasingFunction tests a function for all resolutions, where the value should be decreasing as
// resolution increases.
func testDecreasingFunction(t *testing.T, function func(int32, *float64) H3Error, message string) {
	last := 0.0
	var next float64
	for i := int32(MAX_H3_RES); i >= 0; i-- {
		err := function(i, &next)
		if err != E_SUCCESS {
			t.Errorf("Function failed for resolution %d: %v", i, err)
			return
		}
		if next <= last {
			t.Errorf("%s: expected next (%f) > last (%f) at resolution %d", message, next, last, i)
		}
		last = next
	}
}

func TestRadsToDegs(t *testing.T) {
	t.Parallel()

	originalRads := 1.0
	degs := radsToDegs(originalRads)
	rads := degsToRads(degs)
	if math.Abs(rads-originalRads) >= EPSILON_RAD {
		t.Error("radsToDegs/degsToRads should be invertible")
	}
}

func TestDistanceRads(t *testing.T) {
	t.Parallel()

	var p1, p2 LatLng
	setGeoDegs(&p1, 10, 10)
	setGeoDegs(&p2, 0, 10)

	// TODO: Epsilon is relatively large
	if greatCircleDistanceRads(&p1, &p1) >= EPSILON_RAD*1000 {
		t.Error("0 distance as expected")
	}
	
	expectedDist := degsToRads(10)
	actualDist := greatCircleDistanceRads(&p1, &p2)
	if math.Abs(actualDist-expectedDist) >= EPSILON_RAD*1000 {
		t.Error("distance along longitude as expected")
	}
}

func TestDistanceRadsWrappedLongitude(t *testing.T) {
	t.Parallel()

	negativeLongitude := LatLng{Lat: 0, Lng: -(math.Pi + math.Pi/2)}
	zero := LatLng{Lat: 0, Lng: 0}

	dist1 := greatCircleDistanceRads(&negativeLongitude, &zero)
	if math.Abs(math.Pi/2-dist1) >= EPSILON_RAD {
		t.Error("Distance with wrapped longitude")
	}

	dist2 := greatCircleDistanceRads(&zero, &negativeLongitude)
	if math.Abs(math.Pi/2-dist2) >= EPSILON_RAD {
		t.Error("Distance with wrapped longitude and swapped arguments")
	}
}

func TestDoubleConstants(t *testing.T) {
	t.Parallel()

	// Simple checks for ordering of values
	testDecreasingFunction(t, getHexagonAreaAvgKm2, "getHexagonAreaAvgKm2 ordering")
	testDecreasingFunction(t, getHexagonAreaAvgM2, "getHexagonAreaAvgM2 ordering")
	testDecreasingFunction(t, getHexagonEdgeLengthAvgKm, "getHexagonEdgeLengthAvgKm ordering")
	testDecreasingFunction(t, getHexagonEdgeLengthAvgM, "getHexagonEdgeLengthAvgM ordering")
}

func TestDoubleConstantsErrors(t *testing.T) {
	t.Parallel()

	var out float64

	if getHexagonAreaAvgKm2(-1, &out) != E_RES_DOMAIN {
		t.Error("getHexagonAreaAvgKm2 resolution negative")
	}
	if getHexagonAreaAvgKm2(16, &out) != E_RES_DOMAIN {
		t.Error("getHexagonAreaAvgKm2 resolution too high")
	}
	if getHexagonAreaAvgM2(-1, &out) != E_RES_DOMAIN {
		t.Error("getHexagonAreaAvgM2 resolution negative")
	}
	if getHexagonAreaAvgM2(16, &out) != E_RES_DOMAIN {
		t.Error("getHexagonAreaAvgM2 resolution too high")
	}
	if getHexagonEdgeLengthAvgKm(-1, &out) != E_RES_DOMAIN {
		t.Error("getHexagonEdgeLengthAvgKm resolution negative")
	}
	if getHexagonEdgeLengthAvgKm(16, &out) != E_RES_DOMAIN {
		t.Error("getHexagonEdgeLengthAvgKm resolution too high")
	}
	if getHexagonEdgeLengthAvgM(-1, &out) != E_RES_DOMAIN {
		t.Error("getHexagonEdgeLengthAvgM resolution negative")
	}
	if getHexagonEdgeLengthAvgM(16, &out) != E_RES_DOMAIN {
		t.Error("getHexagonEdgeLengthAvgM resolution too high")
	}
}

func TestIntConstants(t *testing.T) {
	t.Parallel()

	// Simple checks for ordering of values
	last := int64(0)
	for i := int32(0); i <= MAX_H3_RES; i++ {
		next, err := getNumCells(i)
		if err != E_SUCCESS {
			t.Errorf("getNumCells failed for resolution %d: %v", i, err)
			return
		}
		if next <= last {
			t.Error("getNumCells ordering")
		}
		last = next
	}
}

func TestIntConstantsErrors(t *testing.T) {
	t.Parallel()

	_, err := getNumCells(-1)
	if err != E_RES_DOMAIN {
		t.Error("getNumCells resolution negative")
	}
	
	_, err = getNumCells(16)
	if err != E_RES_DOMAIN {
		t.Error("getNumCells resolution too high")
	}
}

func TestNumHexagons(t *testing.T) {
	t.Parallel()

	// Test numHexagon counts of the number of *cells* at each resolution
	expected := []int64{
		122, 842, 5882, 41162, 288122, 2016842, 14117882, 98825162,
		691776122, 4842432842, 33897029882, 237279209162, 1660954464122,
		11626681248842, 81386768741882, 569707381193162,
	}

	for r := int32(0); r <= MAX_H3_RES; r++ {
		num, err := getNumCells(r)
		if err != E_SUCCESS {
			t.Errorf("getNumCells failed for resolution %d: %v", r, err)
			continue
		}
		if num != expected[r] {
			t.Errorf("incorrect numHexagons count for resolution %d: got %d, expected %d", r, num, expected[r])
		}
	}
}