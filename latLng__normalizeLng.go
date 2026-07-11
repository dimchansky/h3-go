package h3

// longitudeNormalization mirrors the C enum in latLng.h.
type longitudeNormalization int

const (
	normalizeNone longitudeNormalization = 0
	normalizeEast longitudeNormalization = 1
	normalizeWest longitudeNormalization = 2
)

// normalizeLng normalizes an input longitude according to the strategy.
// Ported from H3 C: latLng.c::normalizeLng.
func normalizeLng(lng Angle, normalization longitudeNormalization) Angle {
	switch normalization {
	case normalizeEast:
		if lng < 0 {
			return lng + TwoPi
		}
		return lng
	case normalizeWest:
		if lng > 0 {
			return lng - TwoPi
		}
		return lng
	default:
		return lng
	}
}
