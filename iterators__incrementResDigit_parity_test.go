//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_incrementResDigit_parity(t *testing.T) {
	tests := []struct {
		name    string
		h3Index h3Index
		res     int32
	}{
		{
			"increment_at_res_15",
			h3Index(0x85283473fffffff),
			15,
		},
		{
			"increment_at_res_10",
			h3Index(0x8a2834700007fff),
			10,
		},
		{
			"increment_at_res_5",
			h3Index(0x85283400000ffff),
			5,
		},
		{
			"increment_at_res_1",
			h3Index(0x81283fffffffffff),
			1,
		},
		{
			"increment_zero_index_res_15",
			h3Index(0x0),
			15,
		},
		{
			"increment_zero_index_res_10",
			h3Index(0x0),
			10,
		},
		{
			"increment_zero_index_res_5",
			h3Index(0x0),
			5,
		},
		{
			"increment_zero_index_res_1",
			h3Index(0x0),
			1,
		},
		{
			"increment_max_index_res_15",
			h3Index(0xffffffffffffffff),
			15,
		},
		{
			"increment_max_index_res_10",
			h3Index(0xffffffffffffffff),
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Go implementation
			h3Go := tt.h3Index
			incrementResDigit(&h3Go, tt.res)

			// Test C implementation
			h3C := tt.h3Index
			incrementResDigitC(&h3C, tt.res)

			// Compare results
			if h3Go != h3C {
				t.Errorf("incrementResDigit() mismatch: Go=0x%x, C=0x%x", h3Go, h3C)
			}
		})
	}

	// Test deterministic behavior
	t.Run("deterministic", func(t *testing.T) {
		h3Index := h3Index(0x85283473fffffff)
		res := int32(10)

		h1 := h3Index
		h2 := h3Index
		incrementResDigit(&h1, res)
		incrementResDigit(&h2, res)

		if h1 != h2 {
			t.Errorf("incrementResDigit should be deterministic: first=0x%x != second=0x%x", h1, h2)
		}
	})
}
