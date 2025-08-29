// Tests ported from testPentagonIndexes.c

package h3

import (
	"fmt"
	"testing"
)

const PADDED_COUNT = 16

func Test_getPentagons_propertyTests(t *testing.T) {
	t.Parallel()

	expectedCount := pentagonCount()

	for res := int32(0); res <= 15; res++ {
		res := res // capture for parallel subtest
		t.Run(fmt.Sprintf("resolution_%d", res), func(t *testing.T) {
			t.Parallel()

			h3Indexes := make([]H3Index, PADDED_COUNT)
			err := getPentagons(res, h3Indexes)
			if err != E_SUCCESS {
				t.Fatalf("getPentagons failed with error: %v", err)
			}

			var numFound int32

			for i := 0; i < PADDED_COUNT; i++ {
				h3Index := h3Indexes[i]
				if h3Index != 0 {
					numFound++
					
					if !isValidCell(h3Index) {
						t.Errorf("index %v should be valid", h3Index)
					}
					
					if !isPentagon(h3Index) {
						t.Errorf("index %v should be pentagon", h3Index)
					}
					
					if getResolution(h3Index) != res {
						t.Errorf("index %v should have resolution %d, got %d", 
							h3Index, res, getResolution(h3Index))
					}

					// verify uniqueness
					for j := i + 1; j < PADDED_COUNT; j++ {
						if h3Indexes[j] == h3Index {
							t.Errorf("index %v should be seen only once", h3Index)
						}
					}
				}
			}

			if numFound != expectedCount {
				t.Errorf("there should be exactly %d pentagons, found %d", expectedCount, numFound)
			}
		})
	}
}

func Test_getPentagons_invalid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		res  int32
	}{
		{"resolution 16", 16},
		{"resolution 100", 100},
		{"resolution -1", -1},
	}

	for _, tc := range testCases {
		tc := tc // capture for parallel subtest
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h3Indexes := make([]H3Index, PADDED_COUNT)
			err := getPentagons(tc.res, h3Indexes)
			if err != E_RES_DOMAIN {
				t.Errorf("getPentagons of invalid resolution %d should return E_RES_DOMAIN, got %v", 
					tc.res, err)
			}
		})
	}
}

func Test_isPentagon_invalid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		index   H3Index
		isPent  bool
	}{
		{"zero index", H3Index(0), false},
		{"all but high bit", H3Index(0x7fffffffffffffff), false},
	}

	for _, tc := range testCases {
		tc := tc // capture for parallel subtest
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := isPentagon(tc.index)
			if result != tc.isPent {
				t.Errorf("isPentagon(%v) = %v, expected %v", tc.index, result, tc.isPent)
			}
		})
	}
}