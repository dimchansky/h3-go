package h3

import "testing"

// Ported from H3 C: testIndexDigits.c (added in 4.4.0).

func TestIndexDigits_getIndexDigitForCell(t *testing.T) {
	t.Parallel()

	anywhere := LatLng{Lat: 0, Lng: 0}
	var h h3Index

	for resCell := int32(0); resCell <= maxH3Res; resCell++ {
		if err := latLngToCell(&anywhere, resCell, &h); err != eSuccess {
			t.Fatalf("latLngToCell(res=%d): %v", resCell, err)
		}
		for resDigit := int32(1); resDigit <= maxH3Res; resDigit++ {
			var digit int32
			if err := getIndexDigit(h, resDigit, &digit); err != eSuccess {
				t.Fatalf("getIndexDigit(%#x, %d): %v", uint64(h), resDigit, err)
			}
			if resDigit <= resCell {
				if digit < int32(centerDigit) || digit >= int32(invalidDigit) {
					t.Errorf("res %d digit %d: digit %d not in valid range", resCell, resDigit, digit)
				}
			} else if digit != int32(invalidDigit) {
				t.Errorf("res %d digit %d: digit %d should be 'invalid'", resCell, resDigit, digit)
			}
		}
	}

	var digitUnused int32
	if err := getIndexDigit(h, -1, &digitUnused); err != eResDomain {
		t.Errorf("negative resolution: got %v, want eResDomain", err)
	}
	if err := getIndexDigit(h, 0, &digitUnused); err != eResDomain {
		t.Errorf("zero resolution: got %v, want eResDomain", err)
	}
	if err := getIndexDigit(h, 16, &digitUnused); err != eResDomain {
		t.Errorf("too high resolution: got %v, want eResDomain", err)
	}
}

func TestIndexDigits_getIndexDigitForSetCell(t *testing.T) {
	t.Parallel()

	var h h3Index
	for expectedDigit := int32(centerDigit); expectedDigit < int32(invalidDigit); expectedDigit++ {
		for resCell := int32(0); resCell <= maxH3Res; resCell++ {
			setH3Index(&h, resCell, 0, expectedDigit)
			for resDigit := int32(1); resDigit <= maxH3Res; resDigit++ {
				var digit int32
				if err := getIndexDigit(h, resDigit, &digit); err != eSuccess {
					t.Fatalf("getIndexDigit: %v", err)
				}
				if resDigit <= resCell {
					if digit != expectedDigit {
						t.Errorf("digit should be expected: got %d, want %d", digit, expectedDigit)
					}
				} else if digit != int32(invalidDigit) {
					t.Errorf("digit should be 'invalid': got %d", digit)
				}
			}
		}
	}
}
