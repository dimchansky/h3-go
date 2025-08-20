package coordijk

import "testing"

// TestCoordIJKDistance tests Distance function based on H3 C testGridDistanceInternal.c
func TestCoordIJKDistance(t *testing.T) {
    tests := []struct {
        name string
        a, b CoordIJK
        expected int
    }{
        // Identity tests - same coordinate should have distance 0
        {"identity distance 0,0,0", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, 0},
        {"identity distance 1,0,0", CoordIJK{1, 0, 0}, CoordIJK{1, 0, 0}, 0},
        {"identity distance 1,0,1", CoordIJK{1, 0, 1}, CoordIJK{1, 0, 1}, 0},
        {"identity distance 1,1,0", CoordIJK{1, 1, 0}, CoordIJK{1, 1, 0}, 0},
        {"identity distance 0,2,0", CoordIJK{0, 2, 0}, CoordIJK{0, 2, 0}, 0},
        // Directly from H3 C expected results
        {"0,0,0 to 1,0,0", CoordIJK{0, 0, 0}, CoordIJK{1, 0, 0}, 1},
        {"0,0,0 to 0,2,0", CoordIJK{0, 0, 0}, CoordIJK{0, 2, 0}, 2},
        {"0,0,0 to 1,0,1", CoordIJK{0, 0, 0}, CoordIJK{1, 0, 1}, 1},
        {"1,0,0 to 1,0,1", CoordIJK{1, 0, 0}, CoordIJK{1, 0, 1}, 1},
        {"1,0,1 to 0,2,0", CoordIJK{1, 0, 1}, CoordIJK{0, 2, 0}, 3},
        {"1,0,1 to 1,1,0", CoordIJK{1, 0, 1}, CoordIJK{1, 1, 0}, 2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Distance(tt.a, tt.b)
            if result != tt.expected {
                t.Logf("Distance calculation debug: (%v) - (%v) = %v", tt.a, tt.b, tt.a.Sub(tt.b))
                t.Errorf("Distance(%v, %v) = %d, want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

