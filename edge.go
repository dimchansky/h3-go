package h3

// DirectedEdge is an H3 directed edge index: a directed connection from an
// origin cell to one of its adjacent cells.
//
// The zero value is not a valid directed edge; IsValid reports false for it.
// A DirectedEdge may be constructed by conversion from a raw uint64 index
// (unchecked — use IsValid to verify), parsed with ParseDirectedEdge, or
// produced by operations such as Cell.DirectedEdgeTo.
type DirectedEdge uint64
