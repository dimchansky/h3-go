package h3

// _rotate60ccw rotates a direction digit 60 degrees counter-clockwise.
// Mirrors H3's coordijk.h::_rotate60ccw behavior.
// Ported from H3 C: coordijk.h::_rotate60ccw.
func _rotate60ccw(digit direction) direction {
	switch digit {
	case kAxesDigit:
		return ikAxesDigit
	case ikAxesDigit:
		return iAxesDigit
	case iAxesDigit:
		return ijAxesDigit
	case ijAxesDigit:
		return jAxesDigit
	case jAxesDigit:
		return jkAxesDigit
	case jkAxesDigit:
		return kAxesDigit
	default:
		return digit
	}
}
