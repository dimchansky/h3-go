package c2go

import "github.com/dimchansky/h3-go/internal/tables"

// isPentagon returns 1 if the index is a pentagon, else 0 (parity with C).
func isPentagon(h H3Index) int {
	basePent := tables.IsPentagonBaseCell(getBaseCellNumber(h))
	leading := _h3LeadingNonZeroDigit(h)
	if basePent && leading == 0 {
		return 1
	}
	return 0
}
