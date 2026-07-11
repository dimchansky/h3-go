package h3

// Vertex is an H3 vertex index: a single topological vertex of an H3 cell,
// shared by three cells. One of the three neighboring cells is arbitrarily
// designated the vertex's "owner" and determines its canonical index.
//
// The zero value is not a valid vertex; IsValid reports false for it. A
// Vertex may be constructed by conversion from a raw uint64 index (unchecked
// — use IsValid to verify), parsed with ParseVertex, or produced by
// operations such as Cell.Vertex.
type Vertex uint64
