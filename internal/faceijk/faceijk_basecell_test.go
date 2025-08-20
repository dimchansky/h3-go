package faceijk

import (
    "testing"

    "github.com/dimchansky/h3-go/internal/coordijk"
)

func TestFaceIJKToBaseCell(t *testing.T) {
    tests := []struct{
        name string
        fijk FaceIJK
        expectValid bool
    }{
        {"face0_000", FaceIJK{0, coordijk.CoordIJK{0,0,0}}, true},
        {"face0_100", FaceIJK{0, coordijk.CoordIJK{1,0,0}}, true},
        {"face0_200", FaceIJK{0, coordijk.CoordIJK{2,0,0}}, true},
        {"face0_010", FaceIJK{0, coordijk.CoordIJK{0,1,0}}, true},
        {"face0_020", FaceIJK{0, coordijk.CoordIJK{0,2,0}}, true},
        {"face1_000", FaceIJK{1, coordijk.CoordIJK{0,0,0}}, true},
        {"face1_100", FaceIJK{1, coordijk.CoordIJK{1,0,0}}, true},
        {"face5_110", FaceIJK{5, coordijk.CoordIJK{1,1,0}}, true},
        {"face10_001", FaceIJK{10, coordijk.CoordIJK{0,0,1}}, true},
        {"face19_220", FaceIJK{19, coordijk.CoordIJK{2,2,0}}, true},
        {"invalid_face", FaceIJK{-1, coordijk.CoordIJK{0,0,0}}, false},
        {"invalid_i_coord", FaceIJK{0, coordijk.CoordIJK{3,0,0}}, false},
        {"invalid_j_coord", FaceIJK{0, coordijk.CoordIJK{0,3,0}}, false},
        {"invalid_negative", FaceIJK{0, coordijk.CoordIJK{-1,0,0}}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := faceIJKToBaseCell(tt.fijk)
            if tt.expectValid && got < 0 { t.Errorf("faceIJKToBaseCell(%v) = %d, expected >=0", tt.fijk, got) }
            if !tt.expectValid && got >= 0 { t.Errorf("faceIJKToBaseCell(%v) = %d, expected <0", tt.fijk, got) }
        })
    }
}

func TestFaceIJKToBaseCellCCWrot60(t *testing.T) {
    tests := []struct{
        name string
        fijk FaceIJK
        expectValid bool
    }{
        {"face0_000", FaceIJK{0, coordijk.CoordIJK{0,0,0}}, true},
        {"face0_100", FaceIJK{0, coordijk.CoordIJK{1,0,0}}, true},
        {"face5_110", FaceIJK{5, coordijk.CoordIJK{1,1,0}}, true},
        {"invalid_face", FaceIJK{-1, coordijk.CoordIJK{0,0,0}}, false},
        {"invalid_coord", FaceIJK{0, coordijk.CoordIJK{3,0,0}}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := faceIJKToBaseCellCCWrot60(tt.fijk)
            if tt.expectValid && got < 0 { t.Errorf("faceIJKToBaseCellCCWrot60(%v) = %d, expected >=0", tt.fijk, got) }
            if !tt.expectValid && got >= 0 { t.Errorf("faceIJKToBaseCellCCWrot60(%v) = %d, expected <0", tt.fijk, got) }
        })
    }
}

