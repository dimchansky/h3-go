package c2go

// LongitudeNormalization mirrors the C enum in latLng.h
type LongitudeNormalization int

const (
	NORMALIZE_NONE LongitudeNormalization = 0
	NORMALIZE_EAST LongitudeNormalization = 1
	NORMALIZE_WEST LongitudeNormalization = 2
)

// normalizeLng normalizes an input longitude according to the strategy.
// Ported from H3 C: latLng.c::normalizeLng
func normalizeLng(lng float64, normalization LongitudeNormalization) float64 {
	switch normalization {
	case NORMALIZE_EAST:
		if lng < 0 {
			return lng + (2.0 * 3.14159265358979323846)
		}
		return lng
	case NORMALIZE_WEST:
		if lng > 0 {
			return lng - (2.0 * 3.14159265358979323846)
		}
		return lng
	default:
		return lng
	}
}
