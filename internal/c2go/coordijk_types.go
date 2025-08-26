package c2go

// CoordIJK represents IJK hexagon coordinates.
// Each axis is spaced 120 degrees apart.
// Uses int32 to match H3 C implementation exactly (including overflow behavior).
type CoordIJK struct {
	I int32 // i component
	J int32 // j component
	K int32 // k component
}
