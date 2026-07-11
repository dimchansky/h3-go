package h3

// validatePolygonFlags validates whether the flags for the polyfill operation are valid.
// The function checks if any unsupported flags are set and if the containment mode
// is within the valid range.
// Ported from H3 C: polygon.c::validatePolygonFlags.
func validatePolygonFlags(flags uint32) h3Error {
	if (flags&(^flagContainmentModeMask)) != 0 ||
		flagGetContainmentMode(flags) >= ContainmentInvalid {
		return eOptionInvalid
	}
	return eSuccess
}
