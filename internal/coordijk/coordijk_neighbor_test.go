package coordijk

import "testing"

// TestCoordIJKNeighbor tests Neighbor function based on H3 C tests
func TestCoordIJKNeighbor(t *testing.T) {
    tests := []struct {
        name string
        input CoordIJK
        direction Direction
        expected CoordIJK
    }{
        {"center neighbor from origin", CoordIJK{0, 0, 0}, 0, CoordIJK{0, 0, 0}},
        {"i-axis neighbor from origin", CoordIJK{0, 0, 0}, 4, CoordIJK{1, 0, 0}},
        {"j-axis neighbor from origin", CoordIJK{0, 0, 0}, 2, CoordIJK{0, 1, 0}},
        {"k-axis neighbor from origin", CoordIJK{0, 0, 0}, 1, CoordIJK{0, 0, 1}},
        {"invalid direction", CoordIJK{1, 0, 0}, 7, CoordIJK{1, 0, 0}},
        {"center from non-origin", CoordIJK{2, 1, -1}, 0, CoordIJK{2, 1, -1}},
        {"multiple steps", CoordIJK{1, 0, 0}, 1, CoordIJK{1, 0, 1}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.input
            result.Neighbor(tt.direction)
            if result != tt.expected {
                t.Errorf("Neighbor(%d) = %v, want %v", tt.direction, result, tt.expected)
            }
        })
    }
}

