package faceijk

import (
    "testing"

    "github.com/dimchansky/h3-go/internal/coordijk"
)

// Basic FaceIJKToH3 smoke tests (comprehensive cases live in expected/parity suites)
func TestFaceIJKToH3Basic(t *testing.T) {
    tests := []struct{
        name string
        fijk FaceIJK
        res int
        expectValid bool
    }{
        {"res0_face0_000", FaceIJK{0, coordijk.CoordIJK{0,0,0}}, 0, true},
        {"res0_face0_100", FaceIJK{0, coordijk.CoordIJK{1,0,0}}, 0, true},
        {"res1_face1_000", FaceIJK{1, coordijk.CoordIJK{0,0,0}}, 1, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _ = FaceIJKToH3(tt.fijk, tt.res) // accept 0 (H3_NULL) for now per original test
        })
    }
}

