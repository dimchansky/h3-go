// Tests ported from testMathExtensionsInternal.c
package h3

import "testing"

func TestIpow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     int64
		exp      int64
		expected int64
	}{
		{"7 ** 0 == 1", 7, 0, 1},
		{"7 ** 1 == 7", 7, 1, 7},
		{"7 ** 2 == 49", 7, 2, 49},
		{"1 ** 20 == 1", 1, 20, 1},
		{"2 ** 5 == 32", 2, 5, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := _ipow(tt.base, tt.exp)
			if result != tt.expected {
				t.Errorf("_ipow(%d, %d) = %d, expected %d", tt.base, tt.exp, result, tt.expected)
			}
		})
	}
}

func TestSubOverflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        int32
		b        int32
		overflow bool
	}{
		{"0 - 0", 0, 0, false},
		{"int32Min - 0", int32Min, 0, false},
		{"int32Min - 1", int32Min, 1, true},
		{"int32Min - (-1)", int32Min, -1, false},
		{"(int32Min + 1) - 0", int32Min + 1, 0, false},
		{"(int32Min + 1) - 1", int32Min + 1, 1, false},
		{"(int32Min + 1) - (-1)", int32Min + 1, -1, false},
		{"(int32Min + 1) - 2", int32Min + 1, 2, true},
		{"(int32Min + 1) - (-2)", int32Min + 1, -2, false},
		{"100 - 10", 100, 10, false},
		{"int32Max - 0", int32Max, 0, false},
		{"int32Max - 1", int32Max, 1, false},
		{"int32Max - (-1)", int32Max, -1, true},
		{"(int32Max - 1) - 1", int32Max - 1, 1, false},
		{"(int32Max - 1) - (-1)", int32Max - 1, -1, false},
		{"(int32Max - 1) - (-2)", int32Max - 1, -2, true},
		{"int32Min - int32Max", int32Min, int32Max, true},
		{"int32Max - int32Min", int32Max, int32Min, true},
		{"int32Min - int32Min", int32Min, int32Min, false},
		{"int32Max - int32Max", int32Max, int32Max, false},
		{"(-1) - 0", -1, 0, false},
		{"(-1) - 10", -1, 10, false},
		{"(-1) - (-10)", -1, -10, false},
		{"(-1) - int32Max", -1, int32Max, false},
		{"(-2) - int32Max", -2, int32Max, true},
		{"(-1) - int32Min", -1, int32Min, false},
		{"0 - int32Min", 0, int32Min, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subInt32sOverflows(tt.a, tt.b)
			if result != tt.overflow {
				t.Errorf("subInt32sOverflows(%d, %d) = %t, expected %t", tt.a, tt.b, result, tt.overflow)
			}
		})
	}
}

func TestAddOverflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        int32
		b        int32
		overflow bool
	}{
		{"0 + 0", 0, 0, false},
		{"int32Min + 0", int32Min, 0, false},
		{"int32Min + 1", int32Min, 1, false},
		{"int32Min + (-1)", int32Min, -1, true},
		{"(int32Min + 1) + 0", int32Min + 1, 0, false},
		{"(int32Min + 1) + 1", int32Min + 1, 1, false},
		{"(int32Min + 1) + (-1)", int32Min + 1, -1, false},
		{"(int32Min + 1) + 2", int32Min + 1, 2, false},
		{"(int32Min + 1) + (-2)", int32Min + 1, -2, true},
		{"100 + 10", 100, 10, false},
		{"int32Max + 0", int32Max, 0, false},
		{"int32Max + 1", int32Max, 1, true},
		{"int32Max + (-1)", int32Max, -1, false},
		{"(int32Max - 1) + 1", int32Max - 1, 1, false},
		{"(int32Max - 1) + (-1)", int32Max - 1, -1, false},
		{"(int32Max - 1) + (-2)", int32Max - 1, -2, false},
		{"(int32Max - 1) + 2", int32Max - 1, 2, true},
		{"int32Min + int32Max", int32Min, int32Max, false},
		{"int32Max + int32Min", int32Max, int32Min, false},
		{"int32Max + int32Max", int32Max, int32Max, true},
		{"int32Min + int32Min", int32Min, int32Min, true},
		{"(-1) + 0", -1, 0, false},
		{"(-1) + 10", -1, 10, false},
		{"(-1) + (-10)", -1, -10, false},
		{"(-1) + int32Max", -1, int32Max, false},
		{"(-2) + int32Max", -2, int32Max, false},
		{"(-1) + int32Min", -1, int32Min, true},
		{"0 + int32Min", 0, int32Min, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addInt32sOverflows(tt.a, tt.b)
			if result != tt.overflow {
				t.Errorf("addInt32sOverflows(%d, %d) = %t, expected %t", tt.a, tt.b, result, tt.overflow)
			}
		})
	}
}
