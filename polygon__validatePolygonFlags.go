package h3

// validatePolygonFlags validates whether the flags for the polyfill operation are valid.
// The function checks if any unsupported flags are set and if the containment mode
// is within the valid range.
// Ported from H3 C: polygon.c::validatePolygonFlags.
func validatePolygonFlags(flags uint32) H3Error {
	if (flags&(^FLAG_CONTAINMENT_MODE_MASK)) != 0 ||
		FLAG_GET_CONTAINMENT_MODE(flags) >= CONTAINMENT_INVALID {
		return E_OPTION_INVALID
	}
	return E_SUCCESS
}
