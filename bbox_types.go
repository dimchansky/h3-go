package h3

// BBox mirrors the C struct BBox used in bbox.h
// Stores latitude and longitude bounds as angle.Angle (internally in radians).
type BBox struct {
	North Angle
	South Angle
	East  Angle
	West  Angle
}
