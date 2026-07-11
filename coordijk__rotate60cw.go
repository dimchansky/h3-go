package h3

// _rotate60cw rotates a direction digit 60 degrees clockwise.
// Mirrors H3's coordijk.c::_rotate60cw behavior.
// Ported from H3 C: coordijk.c::_rotate60cw.
func _rotate60cw(digit direction) direction {
	switch digit {
	case kAxesDigit:
		return jkAxesDigit
	case jkAxesDigit:
		return jAxesDigit
	case jAxesDigit:
		return ijAxesDigit
	case ijAxesDigit:
		return iAxesDigit
	case iAxesDigit:
		return ikAxesDigit
	case ikAxesDigit:
		return kAxesDigit
	default:
		return digit
	}
}
