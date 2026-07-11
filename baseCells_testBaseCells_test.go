// Tests ported from testBaseCells.c
package h3

import (
	"testing"
)

func TestGetRes0Cells(t *testing.T) {
	t.Parallel()
	count := res0CellCount()
	indexes := make([]h3Index, count)
	err := getRes0Cells(indexes)
	if err != eSuccess {
		t.Fatalf("getRes0Cells failed with error: %v", err)
	}
	if indexes[0] != 0x8001fffffffffff {
		t.Errorf("Expected first basecell to be 0x8001fffffffffff, got 0x%x", indexes[0])
	}
	if indexes[121] != 0x80f3fffffffffff {
		t.Errorf("Expected last basecell to be 0x80f3fffffffffff, got 0x%x", indexes[121])
	}
}
