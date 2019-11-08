package h3

//lint:file-ignore U1000 Ignore all unused code
const (
	// define's of constants and macros for bitwise manipulation of H3Index's.

	H3_NUM_BITS               = 64                               // The number of bits in an H3 index.
	H3_MAX_OFFSET             = 63                               // The bit offset of the max resolution digit in an H3 index.
	H3_MODE_OFFSET            = 59                               // The bit offset of the mode in an H3 index.
	H3_BC_OFFSET              = 45                               // The bit offset of the base cell in an H3 index.
	H3_RES_OFFSET             = 52                               // The bit offset of the resolution in an H3 index.
	H3_RESERVED_OFFSET        = 56                               // The bit offset of the reserved bits in an H3 index.
	H3_PER_DIGIT_OFFSET       = 3                                // The number of bits in a single H3 resolution digit.
	H3_MODE_MASK              = H3Index(15) << H3_MODE_OFFSET    // 1's in the 4 mode bits, 0's everywhere else.
	H3_MODE_MASK_NEGATIVE     = ^H3_MODE_MASK                    // 0's in the 4 mode bits, 1's everywhere else.
	H3_BC_MASK                = H3Index(127) << H3_BC_OFFSET     // 1's in the 7 base cell bits, 0's everywhere else.
	H3_BC_MASK_NEGATIVE       = ^H3_BC_MASK                      // 0's in the 7 base cell bits, 1's everywhere else.
	H3_RES_MASK               = H3Index(15) << H3_RES_OFFSET     // 1's in the 4 resolution bits, 0's everywhere else.
	H3_RES_MASK_NEGATIVE      = ^H3_RES_MASK                     // 0's in the 4 resolution bits, 1's everywhere else.
	H3_RESERVED_MASK          = H3Index(7) << H3_RESERVED_OFFSET // 1's in the 3 reserved bits, 0's everywhere else.
	H3_RESERVED_MASK_NEGATIVE = ^H3_RESERVED_MASK                // 0's in the 3 reserved bits, 1's everywhere else.
	H3_DIGIT_MASK             = H3Index(7)                       // 1's in the 3 bits of res 15 digit bits, 0's everywhere else.
	H3_INIT                   = H3Index(35184372088831)          // H3 index with mode 0, res 0, base cell 0, and 7 for all index digits.

	// H3_INVALID_INDEX index used to indicate an error from geoToH3 and related functions.
	H3_INVALID_INDEX = H3Index(0)
)

// H3_GET_MODE gets the integer mode of h3.
func H3_GET_MODE(h3 H3Index) H3Mode { return H3Mode((h3 & H3_MODE_MASK) >> H3_MODE_OFFSET) }

// H3_SET_MODE sets the integer mode of h3 to v.
func H3_SET_MODE(h3 *H3Index, v H3Mode) {
	*h3 = (*h3 & H3_MODE_MASK_NEGATIVE) | (H3Index(v) << H3_MODE_OFFSET)
}

// H3_GET_BASE_CELL gets the integer base cell of h3.
func H3_GET_BASE_CELL(h3 H3Index) int { return int((h3 & H3_BC_MASK) >> H3_BC_OFFSET) }

// H3_SET_BASE_CELL sets the integer base cell of h3 to bc.
func H3_SET_BASE_CELL(h3 *H3Index, bc int) {
	*h3 = ((*h3) & H3_BC_MASK_NEGATIVE) | (H3Index(bc) << H3_BC_OFFSET)
}

// H3_GET_RESOLUTION gets the integer resolution of h3.
func H3_GET_RESOLUTION(h3 H3Index) int { return int((h3 & H3_RES_MASK) >> H3_RES_OFFSET) }

// H3_SET_RESOLUTION sets the integer resolution of h3.
func H3_SET_RESOLUTION(h3 *H3Index, res int) {
	*h3 = (*h3 & H3_RES_MASK_NEGATIVE) | (H3Index(res) << H3_RES_OFFSET)
}

// H3_GET_INDEX_DIGIT gets the resolution res integer digit (0-7) of h3.
func H3_GET_INDEX_DIGIT(h3 H3Index, res int) Direction {
	return Direction((h3 >> ((MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET)) & H3_DIGIT_MASK)
}

// H3_SET_INDEX_DIGIT sets the resolution res digit of h3 to the integer digit (0-7)
func H3_SET_INDEX_DIGIT(h3 *H3Index, res int, digit Direction) {
	*h3 = (*h3 & ^(H3_DIGIT_MASK << ((MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET))) | (H3Index(digit) << ((MAX_H3_RES - (res)) * H3_PER_DIGIT_OFFSET))
}

// H3_GET_RESERVED_BITS gets a value in the reserved space. Should always be zero for valid indexes.
func H3_GET_RESERVED_BITS(h3 H3Index) int {
	return int((h3 & H3_RESERVED_MASK) >> H3_RESERVED_OFFSET)
}

// H3_SET_RESERVED_BITS sets a value in the reserved space. Setting to non-zero may produce invalid
// indexes.
func H3_SET_RESERVED_BITS(h3 *H3Index, v int) {
	*h3 = (*h3 & H3_RESERVED_MASK_NEGATIVE) | (H3Index(v) << H3_RESERVED_OFFSET)
}

// geoToH3 encodes a coordinate on the sphere to the H3 index of the containing cell at
// the specified resolution.
//
// Returns 0 on invalid input.
//
// `g`: The spherical coordinates to encode.
// `res`: The desired H3 resolution for the encoding.
// Returns the encoded H3 index (or 0 on failure).
func geoToH3(g *GeoCoord, res int) H3Index {
	if res < 0 || res > MAX_H3_RES {
		return H3_INVALID_INDEX
	}

	if !isFinite(g.lat) || !isFinite(g.lon) {
		return H3_INVALID_INDEX
	}

	var fijk FaceIJK
	_geoToFaceIjk(g, res, &fijk)
	return _faceIjkToH3(&fijk, res)
}

// ToGeo determines the spherical coordinates of the center point of an H3 index.
//
// `h`: The H3 index.
// Returns the spherical coordinates of the H3 cell center.
func ToGeo(h H3Index) GeoCoord {
	panic("not implemented")
}

// _faceIjkToH3 converts an FaceIJK address to the corresponding H3 index.
//
// `fijk`: The FaceIJK address.
// `res`: The cell resolution.
// Returns the encoded H3 Index (or 0 on failure).
func _faceIjkToH3(fijk *FaceIJK, res int) H3Index {
	/*

	   // initialize the index
	   H3Index h = H3_INIT;
	   H3_SET_MODE(h, H3_HEXAGON_MODE);
	   H3_SET_RESOLUTION(h, res);

	   // check for res 0/base cell
	   if (res == 0) {
	       if (fijk->coord.i > MAX_FACE_COORD || fijk->coord.j > MAX_FACE_COORD ||
	           fijk->coord.k > MAX_FACE_COORD) {
	           // out of range input
	           return H3_INVALID_INDEX;
	       }

	       H3_SET_BASE_CELL(h, _faceIjkToBaseCell(fijk));
	       return h;
	   }

	   // we need to find the correct base cell FaceIJK for this H3 index;
	   // start with the passed in face and resolution res ijk coordinates
	   // in that face's coordinate system
	   FaceIJK fijkBC = *fijk;

	   // build the H3Index from finest res up
	   // adjust r for the fact that the res 0 base cell offsets the indexing
	   // digits
	   CoordIJK* ijk = &fijkBC.coord;
	   for (int r = res - 1; r >= 0; r--) {
	       CoordIJK lastIJK = *ijk;
	       CoordIJK lastCenter;
	       if (isResClassIII(r + 1)) {
	           // rotate ccw
	           _upAp7(ijk);
	           lastCenter = *ijk;
	           _downAp7(&lastCenter);
	       } else {
	           // rotate cw
	           _upAp7r(ijk);
	           lastCenter = *ijk;
	           _downAp7r(&lastCenter);
	       }

	       CoordIJK diff;
	       _ijkSub(&lastIJK, &lastCenter, &diff);
	       _ijkNormalize(&diff);

	       H3_SET_INDEX_DIGIT(h, r + 1, _unitIjkToDigit(&diff));
	   }

	   // fijkBC should now hold the IJK of the base cell in the
	   // coordinate system of the current face

	   if (fijkBC.coord.i > MAX_FACE_COORD || fijkBC.coord.j > MAX_FACE_COORD ||
	       fijkBC.coord.k > MAX_FACE_COORD) {
	       // out of range input
	       return H3_INVALID_INDEX;
	   }

	   // lookup the correct base cell
	   int baseCell = _faceIjkToBaseCell(&fijkBC);
	   H3_SET_BASE_CELL(h, baseCell);

	   // rotate if necessary to get canonical base cell orientation
	   // for this base cell
	   int numRots = _faceIjkToBaseCellCCWrot60(&fijkBC);
	   if (_isBaseCellPentagon(baseCell)) {
	       // force rotation out of missing k-axes sub-sequence
	       if (_h3LeadingNonZeroDigit(h) == K_AXES_DIGIT) {
	           // check for a cw/ccw offset face; default is ccw
	           if (_baseCellIsCwOffset(baseCell, fijkBC.face)) {
	               h = _h3Rotate60cw(h);
	           } else {
	               h = _h3Rotate60ccw(h);
	           }
	       }

	       for (int i = 0; i < numRots; i++) h = _h3RotatePent60ccw(h);
	   } else {
	       for (int i = 0; i < numRots; i++) {
	           h = _h3Rotate60ccw(h);
	       }
	   }

	   return h;
	*/

	panic("not implemented")
}
