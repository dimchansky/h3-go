package h3

import (
	"testing"
	
	"github.com/dimchansky/h3-go/internal/indexbits"
)

func TestCellIsValid(t *testing.T) {
	tests := []struct {
		name string
		cell Cell
		want bool
	}{
		{
			name: "valid resolution 0 cell",
			cell: Cell(0x08015fffffffffff),
			want: true,
		},
		{
			name: "valid resolution 5 cell",
			cell: Cell(0x085540a73fffffff),
			want: true,
		},
		{
			name: "invalid mode",
			cell: Cell(0x000a000000000000),
			want: false,
		},
		{
			name: "invalid reserved bits",
			cell: Cell(0x090a000000000000),
			want: false,
		},
		{
			name: "invalid base cell",
			cell: Cell(0x08fe1fffffffffff),
			want: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cell.IsValid()
			if got != tt.want {
				t.Errorf("Cell.IsValid(%016x) = %v, want %v", tt.cell, got, tt.want)
			}
		})
	}
}

func TestCellResolution(t *testing.T) {
	tests := []struct {
		name    string
		cell    Cell
		want    int
		wantErr error
	}{
		{
			name:    "resolution 0",
			cell:    Cell(indexbits.Pack(1, 0, 10, nil)),
			want:    0,
			wantErr: nil,
		},
		{
			name:    "resolution 5",
			cell:    Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4})),
			want:    5,
			wantErr: nil,
		},
		{
			name:    "resolution 15",
			cell:    Cell(indexbits.Pack(1, 15, 100, []int{0, 1, 2, 3, 4, 5, 6, 0, 1, 2, 3, 4, 5, 6, 0})),
			want:    15,
			wantErr: nil,
		},
		{
			name:    "invalid cell",
			cell:    Cell(0x000a000000000000), // invalid mode
			want:    0,
			wantErr: ErrCellInvalid,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cell.Resolution()
			if err != tt.wantErr {
				t.Errorf("Cell.Resolution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cell.Resolution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellBaseCell(t *testing.T) {
	tests := []struct {
		name    string
		cell    Cell
		want    int
		wantErr error
	}{
		{
			name:    "base cell 0",
			cell:    Cell(indexbits.Pack(1, 0, 0, nil)),
			want:    0,
			wantErr: nil,
		},
		{
			name:    "base cell 42",
			cell:    Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4})),
			want:    42,
			wantErr: nil,
		},
		{
			name:    "base cell 121",
			cell:    Cell(indexbits.Pack(1, 3, 121, []int{5, 5, 5})),
			want:    121,
			wantErr: nil,
		},
		{
			name:    "invalid cell",
			cell:    Cell(0x000a000000000000), // invalid mode
			want:    0,
			wantErr: ErrCellInvalid,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cell.BaseCell()
			if err != tt.wantErr {
				t.Errorf("Cell.BaseCell() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cell.BaseCell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellIsPentagon(t *testing.T) {
	tests := []struct {
		name    string
		cell    Cell
		want    bool
		wantErr error
	}{
		{
			name:    "resolution 0 pentagon base cell 4",
			cell:    Cell(indexbits.Pack(1, 0, 4, nil)),
			want:    true,
			wantErr: nil,
		},
		{
			name:    "resolution 0 pentagon base cell 14",
			cell:    Cell(indexbits.Pack(1, 0, 14, nil)),
			want:    true,
			wantErr: nil,
		},
		{
			name:    "resolution 0 hexagon base cell 0",
			cell:    Cell(indexbits.Pack(1, 0, 0, nil)),
			want:    false,
			wantErr: nil,
		},
		{
			name:    "resolution 1 pentagon (center child of pentagon base)",
			cell:    Cell(indexbits.Pack(1, 1, 4, []int{0})),
			want:    true,
			wantErr: nil,
		},
		{
			name:    "resolution 2 pentagon (all zeros from pentagon base)",
			cell:    Cell(indexbits.Pack(1, 2, 14, []int{0, 0})),
			want:    true,
			wantErr: nil,
		},
		{
			name:    "resolution 1 hexagon (non-center child of pentagon base)",
			cell:    Cell(indexbits.Pack(1, 1, 4, []int{1})),
			want:    false,
			wantErr: nil,
		},
		{
			name:    "resolution 3 hexagon (has non-zero digit)",
			cell:    Cell(indexbits.Pack(1, 3, 24, []int{0, 0, 2})),
			want:    false,
			wantErr: nil,
		},
		{
			name:    "invalid cell",
			cell:    Cell(0x000a000000000000), // invalid mode
			want:    false,
			wantErr: ErrCellInvalid,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cell.IsPentagon()
			if err != tt.wantErr {
				t.Errorf("Cell.IsPentagon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cell.IsPentagon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxKRingSize(t *testing.T) {
	tests := []struct {
		k    int
		want int
	}{
		{0, 1},      // 3*0*(0+1) + 1 = 1
		{1, 7},      // 3*1*(1+1) + 1 = 7
		{2, 19},     // 3*2*(2+1) + 1 = 19
		{3, 37},     // 3*3*(3+1) + 1 = 37
		{10, 331},   // 3*10*(10+1) + 1 = 331
		{-1, 0},     // negative k returns 0
	}
	
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := MaxKRingSize(tt.k)
			if got != tt.want {
				t.Errorf("MaxKRingSize(%d) = %d, want %d", tt.k, got, tt.want)
			}
		})
	}
}

func TestStubFunctions(t *testing.T) {
	// Test that stub functions return expected errors
	
	t.Run("LatLngToCell", func(t *testing.T) {
		_, err := LatLngToCell(LatLng{Lat: 37.775938, Lng: -122.418307}, 9)
		if err != ErrOptionInvalid {
			t.Errorf("LatLngToCell() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.ToLatLng", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.ToLatLng()
		if err != ErrOptionInvalid {
			t.Errorf("Cell.ToLatLng() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.ToBoundary", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.ToBoundary(nil)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.ToBoundary() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.IsNeighborOf", func(t *testing.T) {
		cell1 := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		cell2 := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 5}))
		_, err := cell1.IsNeighborOf(cell2)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.IsNeighborOf() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.DistanceTo", func(t *testing.T) {
		cell1 := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		cell2 := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 5}))
		_, err := cell1.DistanceTo(cell2)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.DistanceTo() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.KRing", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.KRing(nil, 1)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.KRing() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.HexRange", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.HexRange(nil, 1)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.HexRange() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.HexRangeDistances", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.HexRangeDistances(nil, 1)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.HexRangeDistances() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
	
	t.Run("Cell.HexRing", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.HexRing(nil, 1)
		if err != ErrOptionInvalid {
			t.Errorf("Cell.HexRing() error = %v, want %v", err, ErrOptionInvalid)
		}
	})
}

func TestInputValidation(t *testing.T) {
	t.Run("LatLngToCell invalid resolution", func(t *testing.T) {
		_, err := LatLngToCell(LatLng{Lat: 0, Lng: 0}, -1)
		if err != ErrResolutionDomain {
			t.Errorf("LatLngToCell() with res=-1 error = %v, want %v", err, ErrResolutionDomain)
		}
		
		_, err = LatLngToCell(LatLng{Lat: 0, Lng: 0}, 16)
		if err != ErrResolutionDomain {
			t.Errorf("LatLngToCell() with res=16 error = %v, want %v", err, ErrResolutionDomain)
		}
	})
	
	t.Run("LatLngToCell invalid lat/lng", func(t *testing.T) {
		_, err := LatLngToCell(LatLng{Lat: -91, Lng: 0}, 5)
		if err != ErrLatLngDomain {
			t.Errorf("LatLngToCell() with lat=-91 error = %v, want %v", err, ErrLatLngDomain)
		}
		
		_, err = LatLngToCell(LatLng{Lat: 91, Lng: 0}, 5)
		if err != ErrLatLngDomain {
			t.Errorf("LatLngToCell() with lat=91 error = %v, want %v", err, ErrLatLngDomain)
		}
		
		_, err = LatLngToCell(LatLng{Lat: 0, Lng: -180.1}, 5)
		if err != ErrLatLngDomain {
			t.Errorf("LatLngToCell() with lng=-180.1 error = %v, want %v", err, ErrLatLngDomain)
		}
		
		_, err = LatLngToCell(LatLng{Lat: 0, Lng: 180.1}, 5)
		if err != ErrLatLngDomain {
			t.Errorf("LatLngToCell() with lng=180.1 error = %v, want %v", err, ErrLatLngDomain)
		}
	})
	
	t.Run("Cell.KRing invalid k", func(t *testing.T) {
		cell := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		_, err := cell.KRing(nil, -1)
		if err != ErrDomain {
			t.Errorf("Cell.KRing() with k=-1 error = %v, want %v", err, ErrDomain)
		}
	})
	
	t.Run("Cell.IsNeighborOf resolution mismatch", func(t *testing.T) {
		cell1 := Cell(indexbits.Pack(1, 5, 42, []int{0, 1, 2, 3, 4}))
		cell2 := Cell(indexbits.Pack(1, 6, 42, []int{0, 1, 2, 3, 4, 5}))
		_, err := cell1.IsNeighborOf(cell2)
		if err != ErrResolutionMismatch {
			t.Errorf("Cell.IsNeighborOf() with different resolutions error = %v, want %v", err, ErrResolutionMismatch)
		}
	})
}