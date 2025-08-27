package angle

import (
	"fmt"
	"math"
	"testing"
)

// approx reports whether a and b are within tol.
func approx(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

// degreeEps is a safe epsilon when comparing degree outputs.
const degreeEps = 1e-9

// radEps is a safe epsilon when comparing radians.
const radEps = 1e-12

func TestDegRadRoundTrip(t *testing.T) {
	// Verify Deg -> Angle -> Deg round trip for representative values.
	cases := []float64{-180, -90, -45.5, 0, 23.456789, 45.5, 89.999, 90, 180}
	for _, d := range cases {
		t.Run(fmt.Sprintf("%gdeg", d), func(t *testing.T) {
			got := Deg(d).Deg()
			if !approx(got, d, degreeEps) {
				t.Fatalf("round-trip degrees mismatch: want=%v got=%v", d, got)
			}
		})
	}
}

func TestRadDegConversions(t *testing.T) {
	// Verify Rad and Deg constructors and accessors.
	cases := []float64{-2 * math.Pi, -math.Pi, -math.Pi / 3, 0, math.Pi / 7, math.Pi / 2, math.Pi, 2 * math.Pi}
	for _, r := range cases {
		t.Run(fmt.Sprintf("%grad", r), func(t *testing.T) {
			a := Rad(r)
			if !approx(a.Rad(), r, radEps) {
				t.Fatalf("Rad().Rad() mismatch: want=%v got=%v", r, a.Rad())
			}
			if !approx(Deg(a.Deg()).Rad(), r, radEps) {
				t.Fatalf("Deg(Deg()).Rad() mismatch: want=%v got=%v", r, Deg(a.Deg()).Rad())
			}
		})
	}
}

func TestWrapTwoPi(t *testing.T) {
	twoPi := 2 * math.Pi
	cases := []struct {
		in   float64
		want float64
	}{
		{0, 0},
		{twoPi, 0},
		{-0.1, twoPi - 0.1},
		{4*twoPi + 0.3, 0.3},
		{-4*twoPi - 0.3, twoPi - 0.3},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in=%g", tc.in), func(t *testing.T) {
			got := Rad(tc.in).WrapTwoPi().Rad()
			if !approx(got, tc.want, radEps) {
				t.Fatalf("WrapTwoPi(): want=%v got=%v", tc.want, got)
			}
		})
	}
}

func TestWrapPi(t *testing.T) {
	pi := math.Pi
	twoPi := 2 * math.Pi
	cases := []struct {
		in   float64
		want float64
	}{
		{0, 0},
		{pi, pi},             // π stays π (range is (-π, π])
		{-pi, pi},            // -π wraps to π
		{pi + pi/2, -pi / 2}, // 1.5π -> -0.5π
		{-pi - pi/2, pi / 2}, // -1.5π -> 0.5π
		{3 * pi, pi},         // 3π -> π (NOT -π in this convention)
		{5 * pi, pi},         // 5π -> π
		{10 * twoPi, 0},      // many turns -> 0
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("in=%g", tc.in), func(t *testing.T) {
			got := Rad(tc.in).WrapPi().Rad()
			if !approx(got, tc.want, radEps) {
				t.Fatalf("WrapPi(): want=%v got=%v", tc.want, got)
			}
			// Invariant: result is in (-π, π] and equal to input modulo 2π.
			if !(got > -pi && got <= pi) {
				t.Fatalf("WrapPi(): out of range (-π, π], got %v", got)
			}
			diff := math.Mod(got-tc.in, 2*math.Pi)
			if diff < -1e-12 || diff > 1e-12 {
				t.Fatalf("WrapPi(): not equivalent modulo 2π, diff=%v", diff)
			}
		})
	}
}

func TestArithmetic(t *testing.T) {
	a := Deg(30)  // π/6
	b := Deg(45)  // π/4
	s := a.Add(b) // 75°
	d := b.Sub(a) // 15°
	m := a.Mul(2) // 60°
	v := b.Div(2) // 22.5°
	n := a.Neg()  // -30°

	if !approx(s.Deg(), 75, degreeEps) {
		t.Fatalf("Add(): want 75 got %v", s.Deg())
	}
	if !approx(d.Deg(), 15, degreeEps) {
		t.Fatalf("Sub(): want 15 got %v", d.Deg())
	}
	if !approx(m.Deg(), 60, degreeEps) {
		t.Fatalf("Mul(): want 60 got %v", m.Deg())
	}
	if !approx(v.Deg(), 22.5, degreeEps) {
		t.Fatalf("Div(): want 22.5 got %v", v.Deg())
	}
	if !approx(n.Deg(), -30, degreeEps) {
		t.Fatalf("Neg(): want -30 got %v", n.Deg())
	}
}

func TestEqualApprox(t *testing.T) {
	a := Deg(10)
	b := Deg(10).Add(Rad(5e-13)) // tiny delta in radians

	// Default epsilon path (eps <= 0 -> 1e-12)
	if !a.EqualApprox(b, 0) {
		t.Fatalf("EqualApprox default eps should consider values equal")
	}

	// With tighter tolerance than delta.
	if a.EqualApprox(b, 1e-15) {
		t.Fatalf("EqualApprox with tighter eps should consider values not equal")
	}

	// With custom looser tolerance.
	if !a.EqualApprox(b, 1e-9) {
		t.Fatalf("EqualApprox with custom eps should be equal")
	}
}

func TestTrig(t *testing.T) {
	// Use 30° where sin, cos, tan are well-known exact values.
	a := Deg(30)
	s, c := a.Sin(), a.Cos()
	scS, scC := a.SinCos()

	if !approx(s, 0.5, 1e-12) {
		t.Fatalf("Sin(30°): want 0.5 got %v", s)
	}
	if !approx(c, math.Sqrt(3)/2, 1e-12) {
		t.Fatalf("Cos(30°): want √3/2 got %v", c)
	}
	if !approx(a.Tan(), 1/math.Sqrt(3), 1e-12) {
		t.Fatalf("Tan(30°): want 1/√3 got %v", a.Tan())
	}
	// SinCos should match separate Sin and Cos.
	if !approx(scS, s, 1e-12) || !approx(scC, c, 1e-12) {
		t.Fatalf("SinCos mismatch: want (%v,%v) got (%v,%v)", s, c, scS, scC)
	}
}

func TestString(t *testing.T) {
	// String prints degrees with 6 decimals and a trailing degree sign.
	a := Deg(12.3456789)
	want := "12.345679°" // rounded
	if got := a.String(); got != want {
		t.Fatalf("String(): want %q got %q", want, got)
	}
}

func TestConstants(t *testing.T) {
	if !approx(Pi.Rad(), math.Pi, radEps) {
		t.Fatalf("Pi constant mismatch: want %v got %v", math.Pi, Pi.Rad())
	}
	if !approx(TwoPi.Rad(), 2*math.Pi, radEps) {
		t.Fatalf("TwoPi constant mismatch: want %v got %v", 2*math.Pi, TwoPi.Rad())
	}
}
