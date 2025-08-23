package c2go

// addInt32sOverflows evaluates to true if a + b would overflow for int32.
// Mirrors H3's mathExtensions.h::ADD_INT32S_OVERFLOWS.
func addInt32sOverflows(a, b int32) bool {
	if a > 0 {
		return INT32_MAX-int32(a) < int32(b)
	} else {
		return INT32_MIN-int32(a) > int32(b)
	}
}

// subInt32sOverflows evaluates to true if a - b would overflow for int32.
// Mirrors H3's mathExtensions.h::SUB_INT32S_OVERFLOWS.
func subInt32sOverflows(a, b int32) bool {
	if a >= 0 {
		return INT32_MIN+int32(a) >= int32(b)
	} else {
		return INT32_MAX+int32(a)+1 < int32(b)
	}
}
