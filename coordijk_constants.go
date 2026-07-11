package h3

// direction represents H3 directions/digits (from coordijk.h).
type direction int32

// direction enum (from coordijk.h).
const (
	centerDigit          direction = 0
	kAxesDigit           direction = 1
	jAxesDigit           direction = 2
	jkAxesDigit          direction = jAxesDigit | kAxesDigit // 3
	iAxesDigit           direction = 4
	ikAxesDigit          direction = iAxesDigit | kAxesDigit // 5
	ijAxesDigit          direction = iAxesDigit | jAxesDigit // 6
	invalidDigit         direction = 7
	numDigits            direction = invalidDigit
	pentagonSkippedDigit direction = kAxesDigit
)

// unitVecs mirrors coordijk.h::unitVecs - coordIJK unit vectors corresponding to the 7 H3 digits.
var unitVecs = [numDigits]coordIJK{
	{0, 0, 0}, // direction 0 (centerDigit)
	{0, 0, 1}, // direction 1 (kAxesDigit)
	{0, 1, 0}, // direction 2 (jAxesDigit)
	{0, 1, 1}, // direction 3 (jkAxesDigit)
	{1, 0, 0}, // direction 4 (iAxesDigit)
	{1, 0, 1}, // direction 5 (ikAxesDigit)
	{1, 1, 0}, // direction 6 (ijAxesDigit)
}
