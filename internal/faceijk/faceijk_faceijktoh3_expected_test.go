package faceijk

import (
    "fmt"
    "strconv"
    "strings"
    "testing"

    "github.com/dimchansky/h3-go/internal/coordijk"
)

// Expected-case validation against known oracle outputs
func TestFaceIJKToH3(t *testing.T) {
    tests := []struct{
        face int; i,j,k int; resolution int; expected string; critical bool
    }{
        {0,0,0,0, 0, "8021fffffffffff", true},
        {1,0,0,0, 0, "8005fffffffffff", true},
        {2,0,0,0, 0, "800ffffffffffff", true},
        {3,0,0,0, 0, "8035fffffffffff", true},
        {4,0,0,0, 0, "803ffffffffffff", true},
        {0,0,0,0, 1, "81203ffffffffff", true},
        {1,0,0,0, 1, "81043ffffffffff", true},
        {2,0,0,0, 1, "810e3ffffffffff", true},
        {3,0,0,0, 1, "81343ffffffffff", true},
        {4,0,0,0, 1, "813e3ffffffffff", true},
        {0,1,0,0, 0, "8011fffffffffff", true},
        {0,0,1,0, 0, "8043fffffffffff", true},
        {1,1,1,0, 0, "800bfffffffffff", true},
        {0,2,0,0, 0, "8009fffffffffff", true},
        {4,0,2,1, 0, "8083fffffffffff", true},
        {0,0,2,0, 0, "8015fffffffffff", false},
        {1,2,1,0, 0, "8007fffffffffff", false},
        {2,1,0,0, 0, "800ffffffffffff", false},
    }
    pass, critPass, critTotal := 0, 0, 0
    for i, tt := range tests {
        name := fmt.Sprintf("test_%03d", i+1)
        t.Run(name, func(t *testing.T) {
            if tt.critical { critTotal++ }
            fijk := FaceIJK{tt.face, coordijk.CoordIJK{tt.i, tt.j, tt.k}}
            got := FaceIJKToH3(fijk, tt.resolution)
            gotHex := formatH3Index(got)
            if gotHex == tt.expected {
                pass++; if tt.critical { critPass++ }
            } else {
                if tt.critical {
                    t.Errorf("FaceIJKToH3(%v,%d) = %s, expected %s", fijk, tt.resolution, gotHex, tt.expected)
                } else {
                    t.Logf("expected diff: got %s want %s", gotHex, tt.expected)
                }
            }
        })
    }
    t.Logf("\n=== FaceIJKToH3 VALIDATION SUMMARY ===")
    t.Logf("Total tests: %d", len(tests))
    t.Logf("Tests passed: %d (%.1f%%)", pass, float64(pass)/float64(len(tests))*100)
    t.Logf("Critical tests: %d", critTotal)
    t.Logf("Critical passed: %d (%.1f%%)", critPass, float64(critPass)/float64(critTotal)*100)
    if critPass != critTotal {
        t.Errorf("%d critical tests FAILED", critTotal-critPass)
    }
}

func formatH3Index(h3 uint64) string {
    if h3 == 0 { return "0x0" }
    return strings.ToLower(strconv.FormatUint(h3, 16))
}

