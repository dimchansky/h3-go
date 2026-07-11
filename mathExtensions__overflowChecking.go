package h3

// addInt32sOverflows evaluates to true if a + b would overflow for int32.
// Mirrors H3's mathExtensions.h::ADD_INT32S_OVERFLOWS.
// Ported from H3 C: mathExtensions.h::ADD_INT32S_OVERFLOWS.
func addInt32sOverflows(a, b int32) bool {
	if a > 0 {
		return int32Max-int32(a) < b
	} else {
		return int32Min-int32(a) > b
	}
}

// subInt32sOverflows evaluates to true if a - b would overflow for int32.
// Mirrors H3's mathExtensions.h::SUB_INT32S_OVERFLOWS.
// Ported from H3 C: mathExtensions.h::SUB_INT32S_OVERFLOWS.
func subInt32sOverflows(a, b int32) bool {
	if a >= 0 {
		return int32Min+int32(a) >= b
	} else {
		return int32Max+int32(a)+1 < b
	}
}
