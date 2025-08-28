package h3

// LongitudeNormalization mirrors the C enum in latLng.h
type LongitudeNormalization int

const (
	NORMALIZE_NONE LongitudeNormalization = 0
	NORMALIZE_EAST LongitudeNormalization = 1
	NORMALIZE_WEST LongitudeNormalization = 2
)

// normalizeLng normalizes an input longitude according to the strategy.
// Ported from H3 C: latLng.c::normalizeLng
func normalizeLng(lng Angle, normalization LongitudeNormalization) Angle {
	switch normalization {
	case NORMALIZE_EAST:
		if lng < 0 {
			return lng + TwoPi
		}
		return lng
	case NORMALIZE_WEST:
		if lng > 0 {
			return lng - TwoPi
		}
		return lng
	default:
		return lng
	}
}
