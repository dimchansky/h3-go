package h3

// kadd adds x to the adder using Kahan compensated summation.
//
// Do not "simplify" the four statements algebraically: the algorithm
// depends on the exact rounding of each step, which Go's
// non-reassociating float64 arithmetic preserves as long as they stay
// as written (docs/sync/4.4.0-to-4.5.0.md §5.2).
// Ported from H3 C: adder.h::kadd.
func kadd(a *adder, x float64) {
	y := x - a.c
	t := a.sum + y
	a.c = (t - a.sum) - y
	a.sum = t
}
