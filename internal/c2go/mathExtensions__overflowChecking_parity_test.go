//go:build cgo

package c2go

import "testing"

func Test_addInt32sOverflows_parity(t *testing.T) {
	tests := []struct {
		name string
		a    int32
		b    int32
	}{
		// Basic cases
		{"zeros", 0, 0},
		{"small positive", 1, 2},
		{"small negative", -1, -2},
		{"mixed small", 5, -3},
		{"mixed small reverse", -5, 3},

		// Edge cases around zero
		{"positive zero", 1, 0},
		{"negative zero", -1, 0},
		{"zero positive", 0, 1},
		{"zero negative", 0, -1},

		// Boundary cases - no overflow
		{"max minus 1", INT32_MAX - 1, 1},
		{"min plus 1", INT32_MIN + 1, -1},
		{"max minus max", INT32_MAX, -INT32_MAX},
		{"min plus max", INT32_MIN, INT32_MAX},

		// Overflow cases - positive overflow
		{"max plus 1", INT32_MAX, 1},
		{"max plus max", INT32_MAX, INT32_MAX},
		{"large positive sum", INT32_MAX - 100, 200},
		{"medium positive overflow", 1000000000, 1500000000},

		// Overflow cases - negative overflow
		{"min minus 1", INT32_MIN, -1},
		{"min plus min", INT32_MIN, INT32_MIN},
		{"large negative sum", INT32_MIN + 100, -200},
		{"medium negative overflow", -1000000000, -1500000000},

		// Edge boundary values
		{"exactly max", INT32_MAX, 0},
		{"exactly min", INT32_MIN, 0},
		{"max div 2", INT32_MAX / 2, INT32_MAX / 2},
		{"min div 2", INT32_MIN / 2, INT32_MIN / 2},

		// Random test cases with known outcomes
		{"known safe 1", 123456789, 234567890},
		{"known safe 2", -123456789, 234567890},
		{"known safe 3", 123456789, -234567890},
		{"known safe 4", -123456789, -234567890},

		// Values near boundaries
		{"near max boundary", INT32_MAX - 1000, 999},
		{"near min boundary", INT32_MIN + 1000, -999},
		{"cross zero positive", -1000, 2000},
		{"cross zero negative", 1000, -2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := addInt32sOverflowsC(tt.a, tt.b)

			// Call Go implementation
			gotGo := addInt32sOverflows(tt.a, tt.b)

			// Compare results
			if gotGo != gotC {
				t.Errorf("addInt32sOverflows(%d, %d) mismatch: Go=%v != C=%v",
					tt.a, tt.b, gotGo, gotC)

				// Show the actual sum for debugging (if it fits in int64)
				sum := int64(tt.a) + int64(tt.b)
				wouldOverflow := sum > int64(INT32_MAX) || sum < int64(INT32_MIN)
				t.Logf("  Actual sum as int64: %d, would overflow int32: %v", sum, wouldOverflow)
			}
		})
	}

	// Test commutativity
	t.Run("commutativity", func(t *testing.T) {
		testCases := []struct{ a, b int32 }{
			{100, 200},
			{-100, 200},
			{INT32_MAX - 1, 1},
			{INT32_MIN + 1, -1},
			{1000000000, 500000000},
		}

		for _, tc := range testCases {
			result1Go := addInt32sOverflows(tc.a, tc.b)
			result2Go := addInt32sOverflows(tc.b, tc.a)
			result1C := addInt32sOverflowsC(tc.a, tc.b)
			result2C := addInt32sOverflowsC(tc.b, tc.a)

			if result1Go != result2Go {
				t.Errorf("Go addInt32sOverflows not commutative: (%d,%d)=%v != (%d,%d)=%v",
					tc.a, tc.b, result1Go, tc.b, tc.a, result2Go)
			}

			if result1C != result2C {
				t.Errorf("C addInt32sOverflows not commutative: (%d,%d)=%v != (%d,%d)=%v",
					tc.a, tc.b, result1C, tc.b, tc.a, result2C)
			}
		}
	})
}

func Test_subInt32sOverflows_parity(t *testing.T) {
	tests := []struct {
		name string
		a    int32
		b    int32
	}{
		// Basic cases
		{"zeros", 0, 0},
		{"small positive", 5, 2},
		{"small negative", -5, -2},
		{"mixed small", 5, -3},
		{"mixed small reverse", -5, 3},

		// Edge cases around zero
		{"positive zero", 1, 0},
		{"negative zero", -1, 0},
		{"zero positive", 0, 1},
		{"zero negative", 0, -1},

		// Boundary cases - no overflow
		{"max minus 1", INT32_MAX, 1},
		{"min plus 1", INT32_MIN, -1},
		{"max minus negative", INT32_MAX, -1},
		{"min minus positive", INT32_MIN, 1},

		// Overflow cases - positive overflow (a - b where b is negative)
		{"max minus min", INT32_MAX, INT32_MIN},
		{"max minus large negative", INT32_MAX, -1000},
		{"positive minus large negative", 1000000000, -1500000000},
		{"medium pos minus neg overflow", 1000000000, -1000000000},

		// Overflow cases - negative overflow (a - b where b is positive)
		{"min minus max", INT32_MIN, INT32_MAX},
		{"min minus large positive", INT32_MIN, 1000},
		{"negative minus large positive", -1000000000, 1500000000},
		{"medium neg minus pos overflow", -1000000000, 1000000000},

		// Edge boundary values
		{"exactly max", INT32_MAX, 0},
		{"exactly min", INT32_MIN, 0},
		{"max div 2", INT32_MAX / 2, -(INT32_MAX / 2)},
		{"min div 2", INT32_MIN / 2, INT32_MAX / 2},

		// Random test cases with known outcomes
		{"known safe 1", 123456789, 23456789},
		{"known safe 2", -123456789, -23456789},
		{"known safe 3", 123456789, -23456789},
		{"known safe 4", -123456789, 23456789},

		// Values near boundaries
		{"near max boundary pos", INT32_MAX - 1000, -999},
		{"near min boundary neg", INT32_MIN + 1000, 999},
		{"cross zero positive", -1000, -2000},
		{"cross zero negative", 1000, 2000},

		// Same values (should not overflow)
		{"same positive", 12345, 12345},
		{"same negative", -12345, -12345},
		{"same max", INT32_MAX, INT32_MAX},
		{"same min", INT32_MIN, INT32_MIN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := subInt32sOverflowsC(tt.a, tt.b)

			// Call Go implementation
			gotGo := subInt32sOverflows(tt.a, tt.b)

			// Compare results
			if gotGo != gotC {
				t.Errorf("subInt32sOverflows(%d, %d) mismatch: Go=%v != C=%v",
					tt.a, tt.b, gotGo, gotC)

				// Show the actual difference for debugging (if it fits in int64)
				diff := int64(tt.a) - int64(tt.b)
				wouldOverflow := diff > int64(INT32_MAX) || diff < int64(INT32_MIN)
				t.Logf("  Actual diff as int64: %d, would overflow int32: %v", diff, wouldOverflow)
			}
		})
	}

	// Test that a - a = 0 never overflows (except for MIN - MIN edge case)
	t.Run("self_subtraction", func(t *testing.T) {
		testCases := []int32{0, 1, -1, 100, -100, 123456, -123456, INT32_MAX}

		for _, val := range testCases {
			resultGo := subInt32sOverflows(val, val)
			resultC := subInt32sOverflowsC(val, val)

			// Should never overflow except for INT32_MIN case potentially
			if resultGo != resultC {
				t.Errorf("Self subtraction %d - %d: Go=%v != C=%v", val, val, resultGo, resultC)
			}

			// For most values, subtracting from itself should not overflow
			if val != INT32_MIN {
				if resultGo || resultC {
					t.Errorf("Self subtraction %d - %d should not overflow, but Go=%v, C=%v",
						val, val, resultGo, resultC)
				}
			}
		}
	})
}
