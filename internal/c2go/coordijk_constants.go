package c2go

// Direction represents H3 directions/digits (from coordijk.h)
type Direction int32

// Direction enum (from coordijk.h)
const (
	CENTER_DIGIT           Direction = 0
	K_AXES_DIGIT           Direction = 1
	J_AXES_DIGIT           Direction = 2
	JK_AXES_DIGIT          Direction = J_AXES_DIGIT | K_AXES_DIGIT // 3
	I_AXES_DIGIT           Direction = 4
	IK_AXES_DIGIT          Direction = I_AXES_DIGIT | K_AXES_DIGIT // 5
	IJ_AXES_DIGIT          Direction = I_AXES_DIGIT | J_AXES_DIGIT // 6
	INVALID_DIGIT          Direction = 7
	NUM_DIGITS             Direction = INVALID_DIGIT
	PENTAGON_SKIPPED_DIGIT Direction = K_AXES_DIGIT
)

// UNIT_VECS mirrors coordijk.h::UNIT_VECS - CoordIJK unit vectors corresponding to the 7 H3 digits.
var UNIT_VECS = [NUM_DIGITS]CoordIJK{
	{0, 0, 0}, // direction 0 (CENTER_DIGIT)
	{0, 0, 1}, // direction 1 (K_AXES_DIGIT)
	{0, 1, 0}, // direction 2 (J_AXES_DIGIT)
	{0, 1, 1}, // direction 3 (JK_AXES_DIGIT)
	{1, 0, 0}, // direction 4 (I_AXES_DIGIT)
	{1, 0, 1}, // direction 5 (IK_AXES_DIGIT)
	{1, 1, 0}, // direction 6 (IJ_AXES_DIGIT)
}
