// Tests ported from testPolyfillInternal.c

package h3

import (
	"testing"
)

// Fixtures - equivalent to C struct sfGeoPolygon.
var sfGeoPolygonInternal = GeoPolygon{
	GeoLoop: []LatLng{
		{Lat: 0.659966917655, Lng: -2.1364398519396},
		{Lat: 0.6595011102219, Lng: -2.1359434279405},
		{Lat: 0.6583348114025, Lng: -2.1354884206045},
		{Lat: 0.6581220034068, Lng: -2.1382437718946},
		{Lat: 0.6594479998527, Lng: -2.1384597563896},
		{Lat: 0.6599990002976, Lng: -2.1376771158464},
	},
	Holes: []GeoLoop{}, // numHoles = 0
}

func Test_iterInitPolygonCompact_errors(t *testing.T) {
	t.Parallel()

	t.Run("invalid resolution -1", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, -1, uint32(CONTAINMENT_CENTER))
		if iter.Error != E_RES_DOMAIN {
			t.Errorf("Expected E_RES_DOMAIN, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})

	t.Run("invalid resolution 16", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 16, uint32(CONTAINMENT_CENTER))
		if iter.Error != E_RES_DOMAIN {
			t.Errorf("Expected E_RES_DOMAIN, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})

	t.Run("invalid flags 42", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, 42)
		if iter.Error != E_OPTION_INVALID {
			t.Errorf("Expected E_OPTION_INVALID, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})
}

func Test_iterStepPolygonCompact_invalidCellErrors(t *testing.T) {
	t.Parallel()

	t.Run("invalid cell with bad base cell", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_CENTER))
		if iter.Error != E_SUCCESS {
			t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
		}

		// Give the iterator a cell with a bad base cell
		cell := H3Index(0x85283473fffffff)
		cell = setBaseCell(cell, 123) // Invalid base cell
		iter.Cell = cell

		iterStepPolygonCompact(&iter)
		if iter.Error != E_CELL_INVALID {
			t.Errorf("Expected E_CELL_INVALID, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})

	t.Run("invalid cell with bad base cell at target res", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_CENTER))
		if iter.Error != E_SUCCESS {
			t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
		}

		// Give the iterator a cell with a bad base cell, at the target res
		cell := H3Index(0x89283470003ffff)
		cell = setBaseCell(cell, 123) // Invalid base cell
		iter.Cell = cell

		iterStepPolygonCompact(&iter)
		if iter.Error != E_CELL_INVALID {
			t.Errorf("Expected E_CELL_INVALID, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})

	t.Run("invalid cell with bad base cell full containment", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_FULL))
		if iter.Error != E_SUCCESS {
			t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
		}

		// Give the iterator a cell with a bad base cell, at the target res (full containment)
		cell := H3Index(0x89283470003ffff)
		cell = setBaseCell(cell, 123) // Invalid base cell
		iter.Cell = cell

		iterStepPolygonCompact(&iter)
		if iter.Error != E_CELL_INVALID {
			t.Errorf("Expected E_CELL_INVALID, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})

	t.Run("invalid cell with bad base cell overlapping bbox", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_OVERLAPPING_BBOX))
		if iter.Error != E_SUCCESS {
			t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
		}

		// Give the iterator a cell with a bad base cell, at the target res (overlapping bounding box)
		cell := H3Index(0x89283470003ffff)
		cell = setBaseCell(cell, 123) // Invalid base cell
		iter.Cell = cell

		iterStepPolygonCompact(&iter)
		if iter.Error != E_CELL_INVALID {
			t.Errorf("Expected E_CELL_INVALID, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})

	t.Run("cell too fine for child check", func(t *testing.T) {
		t.Parallel()
		iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_CENTER))
		if iter.Error != E_SUCCESS {
			t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
		}

		// Give the iterator a cell that's too fine for a child check,
		// and a target resolution that allows this to run. This cell has
		// to be inside the polygon to reach the error.
		cell := H3Index(0x8f283080dcb019a)
		iter.Cell = cell
		iter.res = 42 // Access private field for testing

		iterStepPolygonCompact(&iter)
		if iter.Error != E_RES_DOMAIN {
			t.Errorf("Expected E_RES_DOMAIN, got %v", iter.Error)
		}
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL, got %v", iter.Cell)
		}
	})
}

func Test_iterDestroyPolygonCompact(t *testing.T) {
	t.Parallel()

	iter := iterInitPolygonCompact(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_CENTER))
	if iter.Error != E_SUCCESS {
		t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
	}

	iterDestroyPolygonCompact(&iter)
	if iter.Error != E_SUCCESS {
		t.Errorf("Expected E_SUCCESS for destroyed iterator, got %v", iter.Error)
	}
	if iter.Cell != H3_NULL {
		t.Errorf("Expected H3_NULL for destroyed iterator, got %v", iter.Cell)
	}

	// Test that subsequent calls to iterStepPolygonCompact return H3_NULL
	for i := 0; i < 3; i++ {
		iterStepPolygonCompact(&iter)
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL for destroyed iterator on step %d, got %v", i, iter.Cell)
		}
	}
}

func Test_iterDestroyPolygon(t *testing.T) {
	t.Parallel()

	iter := iterInitPolygon(&sfGeoPolygonInternal, 9, uint32(CONTAINMENT_CENTER))
	if iter.Error != E_SUCCESS {
		t.Fatalf("Expected E_SUCCESS for init, got %v", iter.Error)
	}

	iterDestroyPolygon(&iter)
	if iter.Error != E_SUCCESS {
		t.Errorf("Expected E_SUCCESS for destroyed iterator, got %v", iter.Error)
	}
	if iter.Cell != H3_NULL {
		t.Errorf("Expected H3_NULL for destroyed iterator, got %v", iter.Cell)
	}

	// Test that subsequent calls to iterStepPolygon return H3_NULL
	for i := 0; i < 3; i++ {
		iterStepPolygon(&iter)
		if iter.Cell != H3_NULL {
			t.Errorf("Expected H3_NULL for destroyed iterator on step %d, got %v", i, iter.Cell)
		}
	}
}

func Test_cellToBBox_noScale(t *testing.T) {
	t.Parallel()

	// arbitrary cell
	cell := H3Index(0x85283473fffffff)
	bbox, err := cellToBBox(cell, false)
	if err != E_SUCCESS {
		t.Fatalf("Expected E_SUCCESS, got %v", err)
	}

	cellArea, err := cellAreaRads2(cell)
	if err != E_SUCCESS {
		t.Fatalf("Expected E_SUCCESS for cellAreaRads2, got %v", err)
	}

	bboxArea := bboxWidthRads(&bbox) * bboxHeightRads(&bbox)
	ratio := bboxArea / cellArea

	var boundary CellBoundary
	err = cellToBoundary(cell, &boundary)
	if err != E_SUCCESS {
		t.Fatalf("Expected E_SUCCESS for cellToBoundary, got %v", err)
	}

	if ratio >= 3 || ratio <= 1 {
		t.Errorf("Expected reasonable area ratio between cell and bbox (1 < ratio < 3), got %v", ratio)
	}
}

func Test_cellToBBox_boundaryError(t *testing.T) {
	t.Parallel()

	// arbitrary cell with invalid base cell
	cell := H3Index(0x85283473fffffff)
	cell = setBaseCell(cell, 123) // Invalid base cell

	_, err := cellToBBox(cell, false)
	if err != E_CELL_INVALID {
		t.Errorf("Expected E_CELL_INVALID for cell with invalid base cell, got %v", err)
	}
}

func Test_cellToBBox_res0boundaryError(t *testing.T) {
	t.Parallel()

	// arbitrary res 0 cell with invalid base cell
	cell := H3Index(0x8001fffffffffff)
	cell = setBaseCell(cell, 123) // Invalid base cell

	_, err := cellToBBox(cell, false)
	if err != E_CELL_INVALID {
		t.Errorf("Expected E_CELL_INVALID for res 0 cell with invalid base cell, got %v", err)
	}
}

func Test_baseCellNumToCell(t *testing.T) {
	t.Parallel()

	for i := int32(0); i < NUM_BASE_CELLS; i++ {
		cell := baseCellNumToCell(i)
		if !isValidCell(cell) {
			t.Errorf("Cell %v should be valid for base cell %d", cell, i)
		}
		if getBaseCell(cell) != i {
			t.Errorf("Expected base cell %d, got %d", i, getBaseCell(cell))
		}
		if getResolution(cell) != 0 {
			t.Errorf("Expected resolution 0, got %d", getResolution(cell))
		}
	}
}

func Test_baseCellNumToCell_boundaryErrors(t *testing.T) {
	t.Parallel()

	if baseCellNumToCell(-1) != H3_NULL {
		t.Error("Expected H3_NULL for base cell -1")
	}
	if baseCellNumToCell(NUM_BASE_CELLS) != H3_NULL {
		t.Error("Expected H3_NULL for base cell NUM_BASE_CELLS")
	}
}
