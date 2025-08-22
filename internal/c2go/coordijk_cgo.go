//go:build cgo && c2go

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include "coordijk.h"

// Prototypes for the original C helpers in coordijk.c
void _ijkAdd(const CoordIJK* h1, const CoordIJK* h2, CoordIJK* sum);
void _ijkSub(const CoordIJK* h1, const CoordIJK* h2, CoordIJK* diff);
void _setIJK(CoordIJK* ijk, int i, int j, int k);
int _ijkMatches(const CoordIJK* c1, const CoordIJK* c2);
void _ijkScale(CoordIJK* c, int factor);
void _ijkNormalize(CoordIJK* c);
int ijkDistance(const CoordIJK* c1, const CoordIJK* c2);
void _ijkRotate60ccw(CoordIJK* ijk);
void _ijkRotate60cw(CoordIJK* ijk);
Direction _unitIjkToDigit(const CoordIJK* ijk);
void _neighbor(CoordIJK* ijk, Direction digit);
Direction _rotate60ccw(Direction digit);
Direction _rotate60cw(Direction digit);
void ijkToCube(CoordIJK* ijk);
void cubeToIjk(CoordIJK* ijk);
void _ijkToHex2d(const CoordIJK* h, Vec2d* v);
void _hex2dToCoordIJK(const Vec2d* v, CoordIJK* h);
void _upAp7(CoordIJK* ijk);
void _upAp7r(CoordIJK* ijk);
void _downAp7(CoordIJK* ijk);
void _downAp7r(CoordIJK* ijk);
void _downAp3(CoordIJK* ijk);
void _downAp3r(CoordIJK* ijk);
void ijkToIj(const CoordIJK* ijk, CoordIJ* ij);
bool _ijkNormalizeCouldOverflow(const CoordIJK* ijk);
H3Error ijToIjk(const CoordIJ* ij, CoordIJK* ijk);
H3Error _upAp7Checked(CoordIJK* ijk);
H3Error _upAp7rChecked(CoordIJK* ijk);
*/
import "C"

// _ijkAddC calls the original C implementation for adding IJK coordinates.
// Bridges to coordijk.c::_ijkAdd.
func _ijkAddC(h1, h2 *CoordIJK) CoordIJK {
	var ch1, ch2, sum C.CoordIJK
	ch1.i = C.int(h1.I)
	ch1.j = C.int(h1.J)
	ch1.k = C.int(h1.K)
	ch2.i = C.int(h2.I)
	ch2.j = C.int(h2.J)
	ch2.k = C.int(h2.K)
	C._ijkAdd(&ch1, &ch2, &sum)
	return CoordIJK{I: int(sum.i), J: int(sum.j), K: int(sum.k)}
}

// _ijkSubC calls the original C implementation for subtracting IJK coordinates.
// Bridges to coordijk.c::_ijkSub.
func _ijkSubC(h1, h2 *CoordIJK) CoordIJK {
	var ch1, ch2, diff C.CoordIJK
	ch1.i = C.int(h1.I)
	ch1.j = C.int(h1.J)
	ch1.k = C.int(h1.K)
	ch2.i = C.int(h2.I)
	ch2.j = C.int(h2.J)
	ch2.k = C.int(h2.K)
	C._ijkSub(&ch1, &ch2, &diff)
	return CoordIJK{I: int(diff.i), J: int(diff.j), K: int(diff.k)}
}

// _setIJKC calls the original C implementation for setting IJK coordinates.
// Bridges to coordijk.c::_setIJK.
func _setIJKC(i, j, k int) CoordIJK {
	var cijk C.CoordIJK
	C._setIJK(&cijk, C.int(i), C.int(j), C.int(k))
	return CoordIJK{I: int(cijk.i), J: int(cijk.j), K: int(cijk.k)}
}

// _ijkMatchesC calls the original C implementation for comparing IJK coordinates.
// Bridges to coordijk.c::_ijkMatches.
func _ijkMatchesC(c1, c2 *CoordIJK) bool {
	var cc1, cc2 C.CoordIJK
	cc1.i = C.int(c1.I)
	cc1.j = C.int(c1.J)
	cc1.k = C.int(c1.K)
	cc2.i = C.int(c2.I)
	cc2.j = C.int(c2.J)
	cc2.k = C.int(c2.K)
	return C._ijkMatches(&cc1, &cc2) != 0
}

// _ijkScaleC calls the original C implementation for scaling IJK coordinates.
// Bridges to coordijk.c::_ijkScale.
func _ijkScaleC(c *CoordIJK, factor int) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(c.I)
	cc.j = C.int(c.J)
	cc.k = C.int(c.K)
	C._ijkScale(&cc, C.int(factor))
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _ijkNormalizeC calls the original C implementation for normalizing IJK coordinates.
// Bridges to coordijk.c::_ijkNormalize.
func _ijkNormalizeC(c *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(c.I)
	cc.j = C.int(c.J)
	cc.k = C.int(c.K)
	C._ijkNormalize(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// ijkDistanceC calls the original C implementation for computing IJK distance.
// Bridges to coordijk.c::ijkDistance.
func ijkDistanceC(c1, c2 *CoordIJK) int {
	var cc1, cc2 C.CoordIJK
	cc1.i = C.int(c1.I)
	cc1.j = C.int(c1.J)
	cc1.k = C.int(c1.K)
	cc2.i = C.int(c2.I)
	cc2.j = C.int(c2.J)
	cc2.k = C.int(c2.K)
	return int(C.ijkDistance(&cc1, &cc2))
}

// _ijkRotate60ccwC calls the original C implementation for rotating IJK coordinates 60° counter-clockwise.
// Bridges to coordijk.c::_ijkRotate60ccw.
func _ijkRotate60ccwC(c *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(c.I)
	cc.j = C.int(c.J)
	cc.k = C.int(c.K)
	C._ijkRotate60ccw(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _ijkRotate60cwC calls the original C implementation for rotating IJK coordinates 60° clockwise.
// Bridges to coordijk.c::_ijkRotate60cw.
func _ijkRotate60cwC(c *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(c.I)
	cc.j = C.int(c.J)
	cc.k = C.int(c.K)
	C._ijkRotate60cw(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _unitIjkToDigitC calls the original C implementation for converting unit IJK to digit.
// Bridges to coordijk.c::_unitIjkToDigit.
func _unitIjkToDigitC(ijk *CoordIJK) int {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	return int(C._unitIjkToDigit(&cc))
}

// _neighborC calls the original C implementation for applying a direction to IJK coordinates.
// Bridges to coordijk.c::_neighbor.
func _neighborC(ijk *CoordIJK, digit int) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._neighbor(&cc, C.Direction(digit))
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _rotate60ccwC calls the original C implementation for rotating a direction 60° counter-clockwise.
// Bridges to coordijk.c::_rotate60ccw.
func _rotate60ccwC(digit int) int {
	return int(C._rotate60ccw(C.Direction(digit)))
}

// _rotate60cwC calls the original C implementation for rotating a direction 60° clockwise.
// Bridges to coordijk.c::_rotate60cw.
func _rotate60cwC(digit int) int {
	return int(C._rotate60cw(C.Direction(digit)))
}

// ijkToCubeC calls the original C implementation for converting IJK to cube coordinates.
// Bridges to coordijk.c::ijkToCube.
func ijkToCubeC(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C.ijkToCube(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// cubeToIjkC calls the original C implementation for converting cube to IJK coordinates.
// Bridges to coordijk.c::cubeToIjk.
func cubeToIjkC(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C.cubeToIjk(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _ijkToHex2dC calls the original C implementation for converting IJK to 2D hex coordinates.
// Bridges to coordijk.c::_ijkToHex2d.
func _ijkToHex2dC(ijk *CoordIJK) Vec2d {
	var cc C.CoordIJK
	var cv C.Vec2d
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._ijkToHex2d(&cc, &cv)
	return Vec2d{X: float64(cv.x), Y: float64(cv.y)}
}

// _hex2dToCoordIJKC calls the original C implementation for converting 2D hex to IJK coordinates.
// Bridges to coordijk.c::_hex2dToCoordIJK.
func _hex2dToCoordIJKC(v *Vec2d) CoordIJK {
	var cv C.Vec2d
	var cc C.CoordIJK
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)
	C._hex2dToCoordIJK(&cv, &cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _upAp7C calls the original C implementation for aperture 7 up transformation.
// Bridges to coordijk.c::_upAp7.
func _upAp7C(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._upAp7(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _upAp7rC calls the original C implementation for aperture 7 up (clockwise) transformation.
// Bridges to coordijk.c::_upAp7r.
func _upAp7rC(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._upAp7r(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _downAp7C calls the original C implementation for aperture 7 down transformation.
// Bridges to coordijk.c::_downAp7.
func _downAp7C(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._downAp7(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _downAp7rC calls the original C implementation for aperture 7 down (clockwise) transformation.
// Bridges to coordijk.c::_downAp7r.
func _downAp7rC(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._downAp7r(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _downAp3C calls the original C implementation for aperture 3 down transformation.
// Bridges to coordijk.c::_downAp3.
func _downAp3C(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._downAp3(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// _downAp3rC calls the original C implementation for aperture 3 down (clockwise) transformation.
// Bridges to coordijk.c::_downAp3r.
func _downAp3rC(ijk *CoordIJK) CoordIJK {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	C._downAp3r(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}
}

// ijkToIjC calls the original C implementation to convert IJK to IJ coordinates.
// Bridges to coordijk.c::ijkToIj.
func ijkToIjC(ijk *CoordIJK) CoordIJ {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	var cij C.CoordIJ
	C.ijkToIj(&cc, &cij)
	return CoordIJ{I: int(cij.i), J: int(cij.j)}
}

// _ijkNormalizeCouldOverflowC calls the original C implementation to check for overflow.
// Bridges to coordijk.c::_ijkNormalizeCouldOverflow.
func _ijkNormalizeCouldOverflowC(ijk *CoordIJK) bool {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	return bool(C._ijkNormalizeCouldOverflow(&cc))
}

// ijToIjkC calls the original C implementation to convert IJ to IJK coordinates.
// Bridges to coordijk.c::ijToIjk.
func ijToIjkC(ij *CoordIJ) (CoordIJK, H3Error) {
	var cij C.CoordIJ
	cij.i = C.int(ij.I)
	cij.j = C.int(ij.J)
	var cc C.CoordIJK
	err := C.ijToIjk(&cij, &cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}, H3Error(err)
}

// _upAp7CheckedC calls the original C implementation for aperture 7 up transformation with overflow checking.
// Bridges to coordijk.c::_upAp7Checked.
func _upAp7CheckedC(ijk *CoordIJK) (CoordIJK, H3Error) {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	err := C._upAp7Checked(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}, H3Error(err)
}

// _upAp7rCheckedC calls the original C implementation for aperture 7 up (clockwise) transformation with overflow checking.
// Bridges to coordijk.c::_upAp7rChecked.
func _upAp7rCheckedC(ijk *CoordIJK) (CoordIJK, H3Error) {
	var cc C.CoordIJK
	cc.i = C.int(ijk.I)
	cc.j = C.int(ijk.J)
	cc.k = C.int(ijk.K)
	err := C._upAp7rChecked(&cc)
	return CoordIJK{I: int(cc.i), J: int(cc.j), K: int(cc.k)}, H3Error(err)
}
