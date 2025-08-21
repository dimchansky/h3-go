//go:build oracle

package faceijk

import (
    "testing"

    "github.com/dimchansky/h3-go/internal/coordijk"
    testoracle "github.com/dimchansky/h3-go/internal/testoracle"
)

func TestOracle_FaceIJKToH3_Table(t *testing.T) {
    o := testoracle.New(t)
    cases := []FaceIJK{
        {0, coordijk.CoordIJK{0,0,0}},
        {1, coordijk.CoordIJK{0,0,0}},
        {0, coordijk.CoordIJK{1,0,0}},
        {0, coordijk.CoordIJK{0,1,0}},
        {1, coordijk.CoordIJK{1,1,0}},
        {4, coordijk.CoordIJK{0,2,0}},
    }
    for _, f := range cases {
        got := FaceIJKToH3(f, 0)
        want := o.FaceIjkToH3(f.Face, f.Coord.I, f.Coord.J, f.Coord.K, 0)
        if uint64(got) != want { t.Fatalf("FaceIJKToH3(%v,0) = 0x%x, want 0x%x", f, got, want) }
    }
}

func TestOracle_FaceIJKToH3_Randomized(t *testing.T) {
    o := testoracle.New(t)
    r := testoracle.NewRand()
    // Restrict to known-valid unit IJK directions at res 0/1
    units := [][3]int{{0,0,0},{1,0,0},{0,1,0},{0,0,1},{1,1,0},{1,0,1},{0,1,1}}
    for i := 0; i < testoracle.Max(); i++ {
        face := int(r.Intn(20))
        u := units[r.Intn(len(units))]
        res := int(r.Intn(2))
        f := FaceIJK{face, coordijk.CoordIJK{u[0], u[1], u[2]}}
        got := FaceIJKToH3(f, res)
        want := o.FaceIjkToH3(face, u[0], u[1], u[2], res)
        if uint64(got) != want { t.Fatalf("FaceIJKToH3(%v,%d) = 0x%x, want 0x%x", f, res, got, want) }
    }
}
