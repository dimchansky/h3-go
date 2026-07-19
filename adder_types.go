package h3

// adder mirrors the Adder struct from adder.h: Kahan compensated
// summation state (running total plus compensation term). Zero value is
// ready to use.
// Ported from H3 C: adder.h::Adder.
type adder struct {
	sum float64
	c   float64
}
