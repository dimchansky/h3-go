package c2go

// CoordIJK represents IJK hexagon coordinates.
// Each axis is spaced 120 degrees apart.
type CoordIJK struct {
	I int // i component
	J int // j component
	K int // k component
}

// CoordIJ represents IJ coordinates (axial coordinates).
type CoordIJ struct {
	I int // i component
	J int // j component
}
