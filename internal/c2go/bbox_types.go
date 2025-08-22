package c2go

// BBox mirrors the C struct BBox used in bbox.h (radians)
type BBox struct {
	North float64
	South float64
	East  float64
	West  float64
}
