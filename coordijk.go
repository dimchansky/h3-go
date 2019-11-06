package h3

//lint:file-ignore U1000 Ignore all unused code

// CoordIJK holds IJK hexagon coordinates.
// Each axis is spaced 120 degrees apart.
type CoordIJK struct {
	i int // i component
	j int // j component
	k int // k component
}
