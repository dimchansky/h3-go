package h3

import "math"

// InvalidIndex used to indicate an error from `FromGeo` and related functions.
const InvalidIndex = Index(0)

// FromGeo encodes a coordinate on the sphere to the H3 index of the containing cell at
// the specified resolution.
//
// Returns 0 on invalid input.
//
// `g` is the spherical coordinates to encode.
// `res` is the desired H3 resolution for the encoding.
// Returns the encoded H3 index (or 0 on failure).
func FromGeo(g *GeoCoord, res int) Index {
	if res < 0 || res > MaxH3Res {
		return InvalidIndex
	}

	if !isFinite(g.Lat) || !isFinite(g.Lon) {
		return InvalidIndex
	}

	var fijk FaceIJK
	geoToFaceIjk(g, res, &fijk)
	return faceIjkToH3(&fijk, res)
}

// ToGeo determines the spherical coordinates of the center point of an H3 index.
//
// `h` is the H3 index.
// Returns the spherical coordinates of the H3 cell center.
func ToGeo(h Index) GeoCoord {
	panic("not implemented")
}

// faceIjkToH3 converts an FaceIJK address to the corresponding H3 index.
//
// `fijk` is the FaceIJK address.
// `res` is the cell resolution.
// Returns the encoded H3 Index (or 0 on failure).
func faceIjkToH3(fijk *FaceIJK, res int) Index {
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

// isFinite reports whether f is neither NaN nor an infinity.
func isFinite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}
