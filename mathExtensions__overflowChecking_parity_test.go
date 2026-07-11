//go:build cgo && c2go

package h3

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
		{"max minus 1", int32Max - 1, 1},
		{"min plus 1", int32Min + 1, -1},
		{"max minus max", int32Max, -int32Max},
		{"min plus max", int32Min, int32Max},

		// Overflow cases - positive overflow
		{"max plus 1", int32Max, 1},
		{"max plus max", int32Max, int32Max},
		{"large positive sum", int32Max - 100, 200},
		{"medium positive overflow", 1000000000, 1500000000},

		// Overflow cases - negative overflow
		{"min minus 1", int32Min, -1},
		{"min plus min", int32Min, int32Min},
		{"large negative sum", int32Min + 100, -200},
		{"medium negative overflow", -1000000000, -1500000000},

		// Edge boundary values
		{"exactly max", int32Max, 0},
		{"exactly min", int32Min, 0},
		{"max div 2", int32Max / 2, int32Max / 2},
		{"min div 2", int32Min / 2, int32Min / 2},

		// Random test cases with known outcomes
		{"known safe 1", 123456789, 234567890},
		{"known safe 2", -123456789, 234567890},
		{"known safe 3", 123456789, -234567890},
		{"known safe 4", -123456789, -234567890},

		// Values near boundaries
		{"near max boundary", int32Max - 1000, 999},
		{"near min boundary", int32Min + 1000, -999},
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
				wouldOverflow := sum > int64(int32Max) || sum < int64(int32Min)
				t.Logf("  Actual sum as int64: %d, would overflow int32: %v", sum, wouldOverflow)
			}
		})
	}

	// Test commutativity
	t.Run("commutativity", func(t *testing.T) {
		testCases := []struct{ a, b int32 }{
			{100, 200},
			{-100, 200},
			{int32Max - 1, 1},
			{int32Min + 1, -1},
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
		{"max minus 1", int32Max, 1},
		{"min plus 1", int32Min, -1},
		{"max minus negative", int32Max, -1},
		{"min minus positive", int32Min, 1},

		// Overflow cases - positive overflow (a - b where b is negative)
		{"max minus min", int32Max, int32Min},
		{"max minus large negative", int32Max, -1000},
		{"positive minus large negative", 1000000000, -1500000000},
		{"medium pos minus neg overflow", 1000000000, -1000000000},

		// Overflow cases - negative overflow (a - b where b is positive)
		{"min minus max", int32Min, int32Max},
		{"min minus large positive", int32Min, 1000},
		{"negative minus large positive", -1000000000, 1500000000},
		{"medium neg minus pos overflow", -1000000000, 1000000000},

		// Edge boundary values
		{"exactly max", int32Max, 0},
		{"exactly min", int32Min, 0},
		{"max div 2", int32Max / 2, -(int32Max / 2)},
		{"min div 2", int32Min / 2, int32Max / 2},

		// Random test cases with known outcomes
		{"known safe 1", 123456789, 23456789},
		{"known safe 2", -123456789, -23456789},
		{"known safe 3", 123456789, -23456789},
		{"known safe 4", -123456789, 23456789},

		// Values near boundaries
		{"near max boundary pos", int32Max - 1000, -999},
		{"near min boundary neg", int32Min + 1000, 999},
		{"cross zero positive", -1000, -2000},
		{"cross zero negative", 1000, 2000},

		// Same values (should not overflow)
		{"same positive", 12345, 12345},
		{"same negative", -12345, -12345},
		{"same max", int32Max, int32Max},
		{"same min", int32Min, int32Min},
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
				wouldOverflow := diff > int64(int32Max) || diff < int64(int32Min)
				t.Logf("  Actual diff as int64: %d, would overflow int32: %v", diff, wouldOverflow)
			}
		})
	}

	// Test that a - a = 0 never overflows (except for MIN - MIN edge case)
	t.Run("self_subtraction", func(t *testing.T) {
		testCases := []int32{0, 1, -1, 100, -100, 123456, -123456, int32Max}

		for _, val := range testCases {
			resultGo := subInt32sOverflows(val, val)
			resultC := subInt32sOverflowsC(val, val)

			// Should never overflow except for int32Min case potentially
			if resultGo != resultC {
				t.Errorf("Self subtraction %d - %d: Go=%v != C=%v", val, val, resultGo, resultC)
			}

			// For most values, subtracting from itself should not overflow
			if val != int32Min {
				if resultGo || resultC {
					t.Errorf("Self subtraction %d - %d should not overflow, but Go=%v, C=%v",
						val, val, resultGo, resultC)
				}
			}
		}
	})
}
