package h3

// Index fits within a 64-bit unsigned integer
type Index = uint64

// MaxCellBndryVerts is a maximum number of cell boundary vertices; worst case is pentagon:
// 5 original verts + 5 edge crossings
const MaxCellBndryVerts = 10
