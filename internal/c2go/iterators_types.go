package c2go

// IterCellsChildren mirrors the C IterCellsChildren struct for iterator state.
type IterCellsChildren struct {
	H         H3Index // Current H3 index
	ParentRes int32   // Parent resolution
	SkipDigit int32   // Skip digit for pentagons
}
