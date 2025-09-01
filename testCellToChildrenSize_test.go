// Tests ported from testCellToChildrenSize.c
package h3

import (
	"testing"
)

func Test_cellToChildrenSize_hexagon(t *testing.T) {
	t.Parallel()
	
	h := H3Index(0x87283080dffffff) // res 7 *hexagon*
	
	// Test coarser resolution (should return error)
	sz, err := cellToChildrenSize(h, 3)
	if err != E_RES_DOMAIN {
		t.Errorf("got expected size for coarser res: expected E_RES_DOMAIN, got %v", err)
	}
	
	// Test same resolution
	sz, err = cellToChildrenSize(h, 7)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for same res: %v", err)
	}
	if sz != 1 {
		t.Errorf("got expected size for same res: expected 1, got %d", sz)
	}
	
	// Test child resolution
	sz, err = cellToChildrenSize(h, 8)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for child res: %v", err)
	}
	if sz != 7 {
		t.Errorf("got expected size for child res: expected 7, got %d", sz)
	}
	
	// Test grandchild resolution
	sz, err = cellToChildrenSize(h, 9)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for grandchild res: %v", err)
	}
	if sz != 7*7 {
		t.Errorf("got expected size for grandchild res: expected %d, got %d", 7*7, sz)
	}
}

func Test_cellToChildrenSize_pentagon(t *testing.T) {
	t.Parallel()
	
	h := H3Index(0x870800000ffffff) // res 7 *pentagon*
	
	// Test coarser resolution (should return error)
	sz, err := cellToChildrenSize(h, 3)
	if err != E_RES_DOMAIN {
		t.Errorf("got expected size for coarser res: expected E_RES_DOMAIN, got %v", err)
	}
	
	// Test same resolution
	sz, err = cellToChildrenSize(h, 7)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for same res: %v", err)
	}
	if sz != 1 {
		t.Errorf("got expected size for same res: expected 1, got %d", sz)
	}
	
	// Test child resolution
	sz, err = cellToChildrenSize(h, 8)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for child res: %v", err)
	}
	if sz != 6 {
		t.Errorf("got expected size for child res: expected 6, got %d", sz)
	}
	
	// Test grandchild resolution
	sz, err = cellToChildrenSize(h, 9)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for grandchild res: %v", err)
	}
	expected := (5 * 7) + (1 * 6)
	if sz != int64(expected) {
		t.Errorf("got expected size for grandchild res: expected %d, got %d", expected, sz)
	}
}

func Test_cellToChildrenSize_largest_hexagon(t *testing.T) {
	t.Parallel()
	
	h := H3Index(0x806dfffffffffff)     // res 0 *hexagon*
	expected := int64(4747561509943)    // 7^15
	
	out, err := cellToChildrenSize(h, 15)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for largest hexagon: %v", err)
	}
	
	if out != expected {
		t.Errorf("got right size for children 15 levels below: expected %d, got %d", expected, out)
	}
}

func Test_cellToChildrenSize_largest_pentagon(t *testing.T) {
	t.Parallel()
	
	h := H3Index(0x8009fffffffffff)     // res 0 *pentagon*
	expected := int64(3956301258286)    // 1 + 5*(7^15 - 1)/6
	
	out, err := cellToChildrenSize(h, 15)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed for largest pentagon: %v", err)
	}
	
	if out != expected {
		t.Errorf("got right size for children 15 levels below: expected %d, got %d", expected, out)
	}
}