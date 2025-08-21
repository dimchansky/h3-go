package v2d

import (
	"math"
	"testing"
)

func TestVec2d_Mag(t *testing.T) {
	tests := []struct {
		name     string
		v        Vec2d
		expected float64
	}{
		{"zero vector", Vec2d{0, 0}, 0},
		{"unit vector x", Vec2d{1, 0}, 1},
		{"unit vector y", Vec2d{0, 1}, 1},
		{"3-4-5 triangle", Vec2d{3, 4}, 5},
		{"negative components", Vec2d{-3, -4}, 5},
		{"mixed signs", Vec2d{3, -4}, 5},
		{"small values", Vec2d{1e-10, 1e-10}, math.Sqrt(2e-20)},
		{"large values", Vec2d{1e10, 1e10}, math.Sqrt(2e20)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Mag()
			if math.Abs(result-tt.expected) > 1e-14 {
				t.Errorf("Mag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2d_Len2(t *testing.T) {
	tests := []struct {
		name     string
		v        Vec2d
		expected float64
	}{
		{"zero vector", Vec2d{0, 0}, 0},
		{"unit vector x", Vec2d{1, 0}, 1},
		{"unit vector y", Vec2d{0, 1}, 1},
		{"3-4-5 triangle", Vec2d{3, 4}, 25},
		{"negative components", Vec2d{-3, -4}, 25},
		{"mixed signs", Vec2d{3, -4}, 25},
		{"decimal values", Vec2d{1.5, 2.5}, 8.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Len2()
			if math.Abs(result-tt.expected) > 1e-14 {
				t.Errorf("Len2() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2d_AlmostEqual(t *testing.T) {
	tests := []struct {
		name     string
		v1, v2   Vec2d
		expected bool
	}{
		{"identical vectors", Vec2d{1, 2}, Vec2d{1, 2}, true},
		{"zero vectors", Vec2d{0, 0}, Vec2d{0, 0}, true},
		{"within epsilon", Vec2d{1.0, 2.0}, Vec2d{1.0 + Float32Epsilon/2, 2.0}, true},
		{"exactly at epsilon boundary", Vec2d{1.0, 2.0}, Vec2d{1.0 + Float32Epsilon*1.1, 2.0}, false},
		{"beyond epsilon X", Vec2d{1.0, 2.0}, Vec2d{1.0 + Float32Epsilon*2, 2.0}, false},
		{"beyond epsilon Y", Vec2d{1.0, 2.0}, Vec2d{1.0, 2.0 + Float32Epsilon*2}, false},
		{"negative epsilon X", Vec2d{1.0, 2.0}, Vec2d{1.0 - Float32Epsilon/2, 2.0}, true},
		{"negative epsilon Y", Vec2d{1.0, 2.0}, Vec2d{1.0, 2.0 - Float32Epsilon/2}, true},
		{"both components at boundary", Vec2d{1.0, 2.0}, Vec2d{1.0 + Float32Epsilon/2, 2.0 + Float32Epsilon/2}, true},
		{"very small values", Vec2d{1e-10, 1e-10}, Vec2d{1e-10 + Float32Epsilon/2, 1e-10}, true},
		{"large values within epsilon", Vec2d{1e10, 1e10}, Vec2d{1e10 + Float32Epsilon/2, 1e10}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.AlmostEqual(tt.v2)
			if result != tt.expected {
				t.Errorf("AlmostEqual() = %v, want %v (v1=%v, v2=%v, diff=(%v,%v))",
					result, tt.expected, tt.v1, tt.v2,
					math.Abs(tt.v1.X-tt.v2.X), math.Abs(tt.v1.Y-tt.v2.Y))
			}

			// Test symmetry
			result2 := tt.v2.AlmostEqual(tt.v1)
			if result2 != tt.expected {
				t.Errorf("AlmostEqual() is not symmetric: v1.AlmostEqual(v2)=%v, v2.AlmostEqual(v1)=%v", result, result2)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name           string
		p0, p1, p2, p3 Vec2d
		expected       Vec2d
	}{
		{
			name: "perpendicular lines intersecting at origin",
			p0:   Vec2d{-1, 0}, p1: Vec2d{1, 0},
			p2: Vec2d{0, -1}, p3: Vec2d{0, 1},
			expected: Vec2d{0, 0},
		},
		{
			name: "diagonal lines intersecting at (1,1)",
			p0:   Vec2d{0, 0}, p1: Vec2d{2, 2},
			p2: Vec2d{0, 2}, p3: Vec2d{2, 0},
			expected: Vec2d{1, 1},
		},
		{
			name: "horizontal and vertical lines",
			p0:   Vec2d{-2, 3}, p1: Vec2d{4, 3},
			p2: Vec2d{1, -1}, p3: Vec2d{1, 5},
			expected: Vec2d{1, 3},
		},
		{
			name: "lines with negative coordinates",
			p0:   Vec2d{-2, -2}, p1: Vec2d{2, 2},
			p2: Vec2d{-2, 2}, p3: Vec2d{2, -2},
			expected: Vec2d{0, 0},
		},
		{
			name: "non-axis-aligned intersection",
			p0:   Vec2d{0, 0}, p1: Vec2d{2, 1},
			p2: Vec2d{0, 2}, p3: Vec2d{2, 0},
			expected: Vec2d{4.0 / 3.0, 2.0 / 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersect(tt.p0, tt.p1, tt.p2, tt.p3)
			if !result.AlmostEqual(tt.expected) {
				t.Errorf("Intersect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFloat32Epsilon(t *testing.T) {
	// Verify that Float32Epsilon is close to the IEEE-754 FLT_EPSILON
	expectedEpsilon := 1.19209290e-07
	if math.Abs(Float32Epsilon-expectedEpsilon) > 1e-15 {
		t.Errorf("Float32Epsilon = %v, expected approximately %v", Float32Epsilon, expectedEpsilon)
	}
}

// Benchmark tests for performance-critical operations.
func BenchmarkVec2d_Mag(b *testing.B) {
	v := Vec2d{3.141592653589793, 2.718281828459045}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.Mag()
	}
}

func BenchmarkVec2d_Len2(b *testing.B) {
	v := Vec2d{3.141592653589793, 2.718281828459045}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.Len2()
	}
}

func BenchmarkVec2d_AlmostEqual(b *testing.B) {
	v1 := Vec2d{3.141592653589793, 2.718281828459045}
	v2 := Vec2d{3.141592653589794, 2.718281828459046}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v1.AlmostEqual(v2)
	}
}

func BenchmarkIntersect(b *testing.B) {
	p0 := Vec2d{0, 1}
	p1 := Vec2d{3, 4}
	p2 := Vec2d{1, 0}
	p3 := Vec2d{4, 3}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Intersect(p0, p1, p2, p3)
	}
}
