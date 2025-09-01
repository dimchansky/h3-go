package h3

import "fmt"

// bboxPrintln prints a BBox in compact format with a newline.
// Uses bboxPrint internally and adds a newline.
// Ported from H3 C: utility.c::bboxPrintln
func bboxPrintln(bbox *BBox) {
	bboxPrint(bbox)
	fmt.Print("\n")
}
