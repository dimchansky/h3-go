package c2go

// validatePolygonFlags validates polyfill flags (port of polygon.c::validatePolygonFlags)
// Returns 0 on success (E_SUCCESS), non-zero otherwise.
// Ported from H3 C: polygon.c::validatePolygonFlags
func validatePolygonFlags(flags uint32) uint32 {
	if (flags & ^FLAG_CONTAINMENT_MODE_MASK) != 0 || FLAG_GET_CONTAINMENT_MODE(flags) >= CONTAINMENT_INVALID {
		return 1
	}
	return 0
}
