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
		{"INT32_MIN - 0", INT32_MIN, 0, false},
		{"INT32_MIN - 1", INT32_MIN, 1, true},
		{"INT32_MIN - (-1)", INT32_MIN, -1, false},
		{"(INT32_MIN + 1) - 0", INT32_MIN + 1, 0, false},
		{"(INT32_MIN + 1) - 1", INT32_MIN + 1, 1, false},
		{"(INT32_MIN + 1) - (-1)", INT32_MIN + 1, -1, false},
		{"(INT32_MIN + 1) - 2", INT32_MIN + 1, 2, true},
		{"(INT32_MIN + 1) - (-2)", INT32_MIN + 1, -2, false},
		{"100 - 10", 100, 10, false},
		{"INT32_MAX - 0", INT32_MAX, 0, false},
		{"INT32_MAX - 1", INT32_MAX, 1, false},
		{"INT32_MAX - (-1)", INT32_MAX, -1, true},
		{"(INT32_MAX - 1) - 1", INT32_MAX - 1, 1, false},
		{"(INT32_MAX - 1) - (-1)", INT32_MAX - 1, -1, false},
		{"(INT32_MAX - 1) - (-2)", INT32_MAX - 1, -2, true},
		{"INT32_MIN - INT32_MAX", INT32_MIN, INT32_MAX, true},
		{"INT32_MAX - INT32_MIN", INT32_MAX, INT32_MIN, true},
		{"INT32_MIN - INT32_MIN", INT32_MIN, INT32_MIN, false},
		{"INT32_MAX - INT32_MAX", INT32_MAX, INT32_MAX, false},
		{"(-1) - 0", -1, 0, false},
		{"(-1) - 10", -1, 10, false},
		{"(-1) - (-10)", -1, -10, false},
		{"(-1) - INT32_MAX", -1, INT32_MAX, false},
		{"(-2) - INT32_MAX", -2, INT32_MAX, true},
		{"(-1) - INT32_MIN", -1, INT32_MIN, false},
		{"0 - INT32_MIN", 0, INT32_MIN, true},
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
		{"INT32_MIN + 0", INT32_MIN, 0, false},
		{"INT32_MIN + 1", INT32_MIN, 1, false},
		{"INT32_MIN + (-1)", INT32_MIN, -1, true},
		{"(INT32_MIN + 1) + 0", INT32_MIN + 1, 0, false},
		{"(INT32_MIN + 1) + 1", INT32_MIN + 1, 1, false},
		{"(INT32_MIN + 1) + (-1)", INT32_MIN + 1, -1, false},
		{"(INT32_MIN + 1) + 2", INT32_MIN + 1, 2, false},
		{"(INT32_MIN + 1) + (-2)", INT32_MIN + 1, -2, true},
		{"100 + 10", 100, 10, false},
		{"INT32_MAX + 0", INT32_MAX, 0, false},
		{"INT32_MAX + 1", INT32_MAX, 1, true},
		{"INT32_MAX + (-1)", INT32_MAX, -1, false},
		{"(INT32_MAX - 1) + 1", INT32_MAX - 1, 1, false},
		{"(INT32_MAX - 1) + (-1)", INT32_MAX - 1, -1, false},
		{"(INT32_MAX - 1) + (-2)", INT32_MAX - 1, -2, false},
		{"(INT32_MAX - 1) + 2", INT32_MAX - 1, 2, true},
		{"INT32_MIN + INT32_MAX", INT32_MIN, INT32_MAX, false},
		{"INT32_MAX + INT32_MIN", INT32_MAX, INT32_MIN, false},
		{"INT32_MAX + INT32_MAX", INT32_MAX, INT32_MAX, true},
		{"INT32_MIN + INT32_MIN", INT32_MIN, INT32_MIN, true},
		{"(-1) + 0", -1, 0, false},
		{"(-1) + 10", -1, 10, false},
		{"(-1) + (-10)", -1, -10, false},
		{"(-1) + INT32_MAX", -1, INT32_MAX, false},
		{"(-2) + INT32_MAX", -2, INT32_MAX, false},
		{"(-1) + INT32_MIN", -1, INT32_MIN, true},
		{"0 + INT32_MIN", 0, INT32_MIN, false},
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