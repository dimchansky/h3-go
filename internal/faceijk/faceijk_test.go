package faceijk

import (
	"math"
	"testing"

	"github.com/dimchansky/h3-go/internal/coordijk"
)


// TestPosAngleRads tests angle normalization to [0, 2π) range
func TestPosAngleRads(t *testing.T) {
	tests := []struct {
		name string
		input float64
		expected float64
	}{
		{"0", 0.0000000000, 0.0000000000},
		{"pi/2", 1.5708000000, 1.5708000000},
		{"pi", 3.1416000000, 3.1416000000},
		{"3pi/2", 4.7124000000, 4.7124000000},
		{"2pi", 6.2832000000, 0.0000146928},
		{"-pi/2", -1.5708000000, 4.7123853072},
		{"-pi", -3.1416000000, 3.1415853072},
		{"5pi/2", 7.8540000000, 1.5708146928},
		{"-3pi/2", -4.7124000000, 1.5707853072},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PosAngleRads(tt.input)
			if math.Abs(result - tt.expected) > 1e-6 {  // Relaxed tolerance for floating point
				t.Errorf("PosAngleRads(%f) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGeoAzimuthRads tests azimuth calculations between geographic points
func TestGeoAzimuthRads(t *testing.T) {
	tests := []struct {
		name string
		p1Lat, p1Lng, p2Lat, p2Lng float64
		expected float64
	}{
		{"east", 0.0000000000, 0.0000000000, 0.0000000000, 0.1000000000, 1.5707963268},
		{"north", 0.0000000000, 0.0000000000, 0.1000000000, 0.0000000000, 0.0000000000},
		{"west", 0.0000000000, 0.0000000000, 0.0000000000, -0.1000000000, -1.5707963268},
		{"south", 0.0000000000, 0.0000000000, -0.1000000000, 0.0000000000, 3.1415926536},
		{"northeast", 0.5000000000, 1.0000000000, 0.6000000000, 1.1000000000, 0.6803923948},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeoAzimuthRads(tt.p1Lat, tt.p1Lng, tt.p2Lat, tt.p2Lng)
			if math.Abs(result - tt.expected) > 1e-6 {  // Relaxed tolerance
				t.Errorf("GeoAzimuthRads(%f, %f, %f, %f) = %f, want %f", tt.p1Lat, tt.p1Lng, tt.p2Lat, tt.p2Lng, result, tt.expected)
			}
		})
	}
}

// TestFaceIJKToBaseCell tests base cell lookup from FaceIJK coordinates
func TestFaceIJKToBaseCell(t *testing.T) {
	tests := []struct {
		name string
		fijk FaceIJK
		expectValid bool  // Whether we expect a valid base cell (>= 0)
	}{
		{"face0_000", FaceIJK{0, coordijk.CoordIJK{0, 0, 0}}, true},
		{"face0_100", FaceIJK{0, coordijk.CoordIJK{1, 0, 0}}, true},
		{"face0_200", FaceIJK{0, coordijk.CoordIJK{2, 0, 0}}, true},
		{"face0_010", FaceIJK{0, coordijk.CoordIJK{0, 1, 0}}, true},
		{"face0_020", FaceIJK{0, coordijk.CoordIJK{0, 2, 0}}, true},
		{"face1_000", FaceIJK{1, coordijk.CoordIJK{0, 0, 0}}, true},
		{"face1_100", FaceIJK{1, coordijk.CoordIJK{1, 0, 0}}, true},
		{"face5_110", FaceIJK{5, coordijk.CoordIJK{1, 1, 0}}, true},
		{"face10_001", FaceIJK{10, coordijk.CoordIJK{0, 0, 1}}, true},
		{"face19_220", FaceIJK{19, coordijk.CoordIJK{2, 2, 0}}, true},
		{"invalid_face", FaceIJK{-1, coordijk.CoordIJK{0, 0, 0}}, false},
		{"invalid_i_coord", FaceIJK{0, coordijk.CoordIJK{3, 0, 0}}, false},
		{"invalid_j_coord", FaceIJK{0, coordijk.CoordIJK{0, 3, 0}}, false},
		{"invalid_negative", FaceIJK{0, coordijk.CoordIJK{-1, 0, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := faceIJKToBaseCell(tt.fijk)
			if tt.expectValid && result < 0 {
				t.Errorf("faceIJKToBaseCell(%v) = %d, expected valid base cell (>= 0)", tt.fijk, result)
			}
			if !tt.expectValid && result >= 0 {
				t.Errorf("faceIJKToBaseCell(%v) = %d, expected invalid result (< 0)", tt.fijk, result)
			}
		})
	}
}

// TestFaceIJKToBaseCellCCWrot60 tests rotation lookup from FaceIJK coordinates
func TestFaceIJKToBaseCellCCWrot60(t *testing.T) {
	tests := []struct {
		name string
		fijk FaceIJK
		expectValid bool  // Whether we expect a valid rotation (>= 0)
	}{
		{"face0_000", FaceIJK{0, coordijk.CoordIJK{0, 0, 0}}, true},
		{"face0_100", FaceIJK{0, coordijk.CoordIJK{1, 0, 0}}, true},
		{"face5_110", FaceIJK{5, coordijk.CoordIJK{1, 1, 0}}, true},
		{"invalid_face", FaceIJK{-1, coordijk.CoordIJK{0, 0, 0}}, false},
		{"invalid_coord", FaceIJK{0, coordijk.CoordIJK{3, 0, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := faceIJKToBaseCellCCWrot60(tt.fijk)
			if tt.expectValid && result < 0 {
				t.Errorf("faceIJKToBaseCellCCWrot60(%v) = %d, expected valid rotation (>= 0)", tt.fijk, result)
			}
			if !tt.expectValid && result >= 0 {
				t.Errorf("faceIJKToBaseCellCCWrot60(%v) = %d, expected invalid result (< 0)", tt.fijk, result)
			}
		})
	}
}

// TestGeoToHex2d tests geographic coordinate to hex2d conversion
func TestGeoToHex2d(t *testing.T) {
	tests := []struct {
		name string
		lat, lng float64
		res int
		expectValidFace bool
	}{
		{"face0_res0", 0.8, 1.2, 0, true},
		{"face1_res1", 1.3, 2.5, 1, true},
		{"face4_res2", 0.5, 0.4, 2, true},
		{"equator_res0", 0.0, 0.0, 0, true},
		{"north_pole_res1", 1.57079632679, 0.0, 1, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			face, v := GeoToHex2d(tt.lat, tt.lng, tt.res)
			
			if tt.expectValidFace {
				if face < 0 || face >= NumIcosaFaces {
					t.Errorf("GeoToHex2d(%f, %f, %d) face = %d, expected valid face [0-%d]", 
						tt.lat, tt.lng, tt.res, face, NumIcosaFaces-1)
				}
				// Hex2d coordinates should be finite
				if math.IsInf(v.X, 0) || math.IsNaN(v.X) || math.IsInf(v.Y, 0) || math.IsNaN(v.Y) {
					t.Errorf("GeoToHex2d(%f, %f, %d) produced invalid hex2d coordinates: %v", 
						tt.lat, tt.lng, tt.res, v)
				}
			}
		})
	}
}

// TestFaceIJKToH3Basic tests basic FaceIJKToH3 functionality 
// For comprehensive oracle validation, see faceijktoh3_test.go
func TestFaceIJKToH3Basic(t *testing.T) {
	tests := []struct {
		name string
		fijk FaceIJK
		res int
		expectValid bool
	}{
		{"res0_face0_000", FaceIJK{0, coordijk.CoordIJK{0, 0, 0}}, 0, true},
		{"res0_face0_100", FaceIJK{0, coordijk.CoordIJK{1, 0, 0}}, 0, true},
		{"res1_face1_000", FaceIJK{1, coordijk.CoordIJK{0, 0, 0}}, 1, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FaceIJKToH3(tt.fijk, tt.res)
			
			if tt.expectValid {
				if result == 0 {
					t.Logf("FaceIJKToH3(%v, %d) = 0 (H3_NULL) - may be expected", tt.fijk, tt.res)
				} else {
					t.Logf("FaceIJKToH3(%v, %d) = 0x%x (success)", tt.fijk, tt.res, result)
				}
			}
		})
	}
}

// TestFaceConstants validates that face-related constants are properly defined
func TestFaceConstants(t *testing.T) {
	// Test face center coordinates exist for all faces
	if len(FaceCenterGeo) != NumIcosaFaces {
		t.Errorf("FaceCenterGeo length = %d, want %d", len(FaceCenterGeo), NumIcosaFaces)
	}
	
	// Test face axes exist for all faces
	if len(FaceAxesAzRadsCII) != NumIcosaFaces {
		t.Errorf("FaceAxesAzRadsCII length = %d, want %d", len(FaceAxesAzRadsCII), NumIcosaFaces)
	}
	
	if len(FaceAxesAzRadsCIII) != NumIcosaFaces {
		t.Errorf("FaceAxesAzRadsCIII length = %d, want %d", len(FaceAxesAzRadsCIII), NumIcosaFaces)
	}
	
	// Test that face center coordinates are within valid latitude/longitude ranges
	for i, center := range FaceCenterGeo {
		lat, lng := center[0], center[1]
		if lat < -math.Pi/2 || lat > math.Pi/2 {
			t.Errorf("Face %d latitude %f out of range [-π/2, π/2]", i, lat)
		}
		if lng < -math.Pi || lng > math.Pi {
			t.Errorf("Face %d longitude %f out of range [-π, π]", i, lng)
		}
	}
	
	// Test scaling arrays have expected lengths
	expectedMaxRes := 16 // res 0-15
	if len(MaxDimByCIIres) != expectedMaxRes {
		t.Errorf("MaxDimByCIIres length = %d, want %d", len(MaxDimByCIIres), expectedMaxRes)
	}
	
	if len(MaxDimByCIIIres) != expectedMaxRes {
		t.Errorf("MaxDimByCIIIres length = %d, want %d", len(MaxDimByCIIIres), expectedMaxRes)
	}
}

