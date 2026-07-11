package h3

// bbox mirrors the C struct bbox used in bbox.h
// Stores latitude and longitude bounds as angle.Angle (internally in radians).
type bbox struct {
	North Angle
	South Angle
	East  Angle
	West  Angle
}
