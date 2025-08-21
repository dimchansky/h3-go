package faceijk

import (
	"fmt"
	"testing"
)

// TestIsResolutionClassII validates IsResolutionClassII implementation.
func TestIsResolutionClassII(t *testing.T) {
	tests := []struct {
		resolution int
		expected   bool
	}{
		// test_001 to test_016: resolutions 0-15 (even=true, odd=false for ClassII)
		{0, true}, {1, false}, {2, true}, {3, false}, {4, true}, {5, false},
		{6, true}, {7, false}, {8, true}, {9, false}, {10, true}, {11, false},
		{12, true}, {13, false}, {14, true}, {15, false},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("test_%03d", i+1)
		t.Run(testName, func(t *testing.T) {
			result := IsResolutionClassII(tt.resolution)

			if result != tt.expected {
				t.Errorf("IsResolutionClassII(%d) = %t, expected %t", tt.resolution, result, tt.expected)
			} else {
				t.Logf("✅ PASS: IsResolutionClassII(%d) = %t", tt.resolution, result)
			}
		})
	}
}

// TestIsResolutionClassIII validates IsResolutionClassIII implementation.
func TestIsResolutionClassIII(t *testing.T) {
	tests := []struct {
		resolution int
		expected   bool
	}{
		// test_001 to test_016: resolutions 0-15 (even=false, odd=true for ClassIII)
		{0, false}, {1, true}, {2, false}, {3, true}, {4, false}, {5, true},
		{6, false}, {7, true}, {8, false}, {9, true}, {10, false}, {11, true},
		{12, false}, {13, true}, {14, false}, {15, true},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("test_%03d", i+1)
		t.Run(testName, func(t *testing.T) {
			result := IsResolutionClassIII(tt.resolution)

			if result != tt.expected {
				t.Errorf("IsResolutionClassIII(%d) = %t, expected %t", tt.resolution, result, tt.expected)
			} else {
				t.Logf("✅ PASS: IsResolutionClassIII(%d) = %t", tt.resolution, result)
			}
		})
	}
}

// TestIsResolutionClass tests both resolution classification functions together.
func TestIsResolutionClass(t *testing.T) {
	tests := []struct {
		res         int
		expectedII  bool
		expectedIII bool
	}{
		{0, true, false},
		{1, false, true},
		{2, true, false},
		{3, false, true},
		{4, true, false},
		{5, false, true},
		{6, true, false},
		{7, false, true},
		{8, true, false},
		{9, false, true},
		{10, true, false},
		{11, false, true},
		{12, true, false},
		{13, false, true},
		{14, true, false},
		{15, false, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("res_%d", tt.res), func(t *testing.T) {
			if IsResolutionClassII(tt.res) != tt.expectedII {
				t.Errorf("IsResolutionClassII(%d) = %v, want %v", tt.res, IsResolutionClassII(tt.res), tt.expectedII)
			}
			if IsResolutionClassIII(tt.res) != tt.expectedIII {
				t.Errorf("IsResolutionClassIII(%d) = %v, want %v", tt.res, IsResolutionClassIII(tt.res), tt.expectedIII)
			}
		})
	}
}
