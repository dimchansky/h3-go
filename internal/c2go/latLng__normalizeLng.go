package c2go

import "github.com/dimchansky/h3-go/angle"

// LongitudeNormalization mirrors the C enum in latLng.h
type LongitudeNormalization int

const (
	NORMALIZE_NONE LongitudeNormalization = 0
	NORMALIZE_EAST LongitudeNormalization = 1
	NORMALIZE_WEST LongitudeNormalization = 2
)

// normalizeLng normalizes an input longitude according to the strategy.
// Ported from H3 C: latLng.c::normalizeLng
func normalizeLng(lng angle.Angle, normalization LongitudeNormalization) angle.Angle {
	switch normalization {
	case NORMALIZE_EAST:
		if lng < 0 {
			return lng + angle.TwoPi
		}
		return lng
	case NORMALIZE_WEST:
		if lng > 0 {
			return lng - angle.TwoPi
		}
		return lng
	default:
		return lng
	}
}
