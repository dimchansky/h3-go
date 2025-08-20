package faceijk

import (
    "math"
    "testing"
)

// Angle normalization to [0, 2π)
func TestPosAngleRads(t *testing.T) {
    tests := []struct{
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
            got := PosAngleRads(tt.input)
            if math.Abs(got-tt.expected) > 1e-6 {
                t.Errorf("PosAngleRads(%f) = %f, want %f", tt.input, got, tt.expected)
            }
        })
    }
}

// Azimuth between geographic points
func TestGeoAzimuthRads(t *testing.T) {
    tests := []struct{
        name string
        p1Lat, p1Lng, p2Lat, p2Lng float64
        expected float64
    }{
        {"east", 0, 0, 0, 0.1, 1.5707963268},
        {"north", 0, 0, 0.1, 0, 0},
        {"west", 0, 0, 0, -0.1, -1.5707963268},
        {"south", 0, 0, -0.1, 0, 3.1415926536},
        {"northeast", 0.5, 1.0, 0.6, 1.1, 0.6803923948},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := GeoAzimuthRads(tt.p1Lat, tt.p1Lng, tt.p2Lat, tt.p2Lng)
            if math.Abs(got-tt.expected) > 1e-6 {
                t.Errorf("GeoAzimuthRads(%f,%f,%f,%f) = %f, want %f", tt.p1Lat, tt.p1Lng, tt.p2Lat, tt.p2Lng, got, tt.expected)
            }
        })
    }
}

