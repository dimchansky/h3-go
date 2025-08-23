package c2go

// _rotate60cw rotates a direction digit 60 degrees clockwise.
// Mirrors H3's coordijk.c::_rotate60cw behavior.
// Ported from H3 C: coordijk.c::_rotate60cw
func _rotate60cw(digit Direction) Direction {
	switch digit {
	case K_AXES_DIGIT:
		return JK_AXES_DIGIT
	case JK_AXES_DIGIT:
		return J_AXES_DIGIT
	case J_AXES_DIGIT:
		return IJ_AXES_DIGIT
	case IJ_AXES_DIGIT:
		return I_AXES_DIGIT
	case I_AXES_DIGIT:
		return IK_AXES_DIGIT
	case IK_AXES_DIGIT:
		return K_AXES_DIGIT
	default:
		return digit
	}
}
