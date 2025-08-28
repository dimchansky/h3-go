package c2go

import "github.com/dimchansky/h3-go/angle"

// BBox mirrors the C struct BBox used in bbox.h
// Stores latitude and longitude bounds as angle.Angle (internally in radians).
type BBox struct {
	North angle.Angle
	South angle.Angle
	East  angle.Angle
	West  angle.Angle
}
