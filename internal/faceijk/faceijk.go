// Package faceijk implements face-centered IJK coordinate transformations
// for the H3 icosahedral projection.
package faceijk

import (
	"math"
	
	"github.com/dimchansky/h3-go/internal/coordijk"
	"github.com/dimchansky/h3-go/internal/indexbits"
	"github.com/dimchansky/h3-go/internal/tables"
)

// NumIcosaFaces is the number of faces on an icosahedron.
const NumIcosaFaces = 20

// InvalidFace represents an invalid face index.
const InvalidFace = -1

// FaceIJK represents a face number and IJK coordinates on that
// face-centered coordinate system.
type FaceIJK struct {
	Face  int              // Face number (0-19)
	Coord coordijk.CoordIJK // IJK coordinates on that face
}

// FaceOrientIJK contains information to transform into an adjacent face IJK system.
type FaceOrientIJK struct {
	Face      int              // Face number
	Translate coordijk.CoordIJK // Resolution 0 translation relative to primary face
	CCWRot60  int              // Number of 60 degree CCW rotations relative to primary face
}

// Overage represents the type of overage when transforming between faces.
type Overage int

const (
	// NoOverage means coordinates are on the original face.
	NoOverage Overage = 0
	// FaceEdge means coordinates are on a face edge (substrate grids only).
	FaceEdge Overage = 1
	// NewFace means coordinates have overflowed to a new face interior.
	NewFace Overage = 2
)

// Face neighbor indexing
const (
	IJ = 1 // IJ quadrant face neighbor
	KI = 2 // KI quadrant face neighbor
	JK = 3 // JK quadrant face neighbor
)

// Constants from H3 source for gnomonic projection
const (
	// M_AP7_ROT_RADS is the rotation angle between Class II and Class III resolution axes
	M_AP7_ROT_RADS = 0.333473172251832115336090755351 // arcsin(sqrt(3/28))
	
	// RES0_U_GNOMONIC is the scaling factor for gnomonic projection at resolution 0
	RES0_U_GNOMONIC = 0.38196601125010515179541316563436
	
	// Constants for azimuth calculation
	M_SQRT7 = 2.6457513110645905509222807479026
)

// FaceCenterGeo contains the center point of each icosahedron face in lat/lng radians.
// From H3 C source faceijk.c
var FaceCenterGeo = [NumIcosaFaces][2]float64{
	{0.803582649718989942, 1.248397419617396099},    // face 0
	{1.307747883455638156, 2.536945009877921159},    // face 1
	{1.054751253523952054, -1.347517358900396623},   // face 2
	{0.600191595538186799, -0.450603909469755746},   // face 3
	{0.491715428198773866, 0.401988202911306943},    // face 4
	{0.172745327415618701, 1.678146885280433686},    // face 5
	{0.605929321571350690, 2.953923329812411617},    // face 6
	{0.427370518328979641, -1.888876200336285401},   // face 7
	{-0.079066118549212831, -0.733429513380867741},  // face 8
	{-0.230961644455383637, 0.506495587332349035},   // face 9
	{0.079066118549212831, 2.408163140208925497},    // face 10
	{0.230961644455383637, -2.635097066257444203},   // face 11
	{-0.172745327415618701, -1.463445768309359553},  // face 12
	{-0.605929321571350690, -0.187669323777381622},  // face 13
	{-0.427370518328979641, 1.252716453253507838},   // face 14
	{-0.600191595538186799, 2.690988744120037492},   // face 15
	{-0.491715428198773866, -2.739604450678486295},  // face 16
	{-0.803582649718989942, -1.893195233972397139},  // face 17
	{-1.307747883455638156, -0.604647643711872080},  // face 18
	{-1.054751253523952054, 1.794075294689396615},   // face 19
}

// FaceAxesAzRadsCII contains the azimuth in radians for the first axis of each face
// for Class II resolutions.
var FaceAxesAzRadsCII = [NumIcosaFaces][3]float64{
	{5.619958268523939882, 3.525563166130744542, 1.431168063737548730},  // face 0
	{5.760339081714187279, 3.665943979320991689, 1.571548876927796127},  // face 1
	{0.780213654393430055, 4.969003859179821079, 2.874608756786625655},  // face 2
	{0.430469363979999913, 4.619259568766391033, 2.524864466373195467},  // face 3
	{6.130269123335111400, 4.035874020941915804, 1.941478918548720291},  // face 4
	{2.692877706530642877, 0.598482604137447119, 4.787272808923838195},  // face 5
	{2.982963003477243874, 0.888567901084048369, 5.077358105870439581},  // face 6
	{3.532912002790141181, 1.438516900396945656, 5.627307105183336758},  // face 7
	{3.494305004758171946, 1.399909902364976048, 5.588700107151367234},  // face 8
	{3.003214169499538391, 0.908819067106342928, 5.097609271892733906},  // face 9
	{5.930472956509811562, 3.836077854116615875, 1.741682751723420374},  // face 10
	{0.138378484090254847, 4.327168688876645809, 2.232773586483450307},  // face 11
	{0.448714947059150361, 4.637505151845541521, 2.543110049452346120},  // face 12
	{0.158629650112549365, 4.347419854898940135, 2.253024752505744869},  // face 13
	{5.891865957979862386, 3.797470855586666806, 1.703075753193471330},  // face 14
	{5.401775122721453543, 3.307380020328257965, 1.212984917935062374},  // face 15
	{5.611690825861418396, 3.517295723468222704, 1.422900621075027214},  // face 16
	{5.321605529792489535, 3.227210427399293958, 1.132815325006098423},  // face 17
	{5.471966437919222939, 3.377571335526027311, 1.283176233132831783},  // face 18
	{4.831696968776721443, 2.737301866383525914, 0.642906764990330379},  // face 19
}

// FaceAxesAzRadsCIII contains the azimuth in radians for the first axis of each face
// for Class III resolutions.
var FaceAxesAzRadsCIII = [NumIcosaFaces][3]float64{
	{5.619958268523939882, 3.525563166130744542, 1.431168063737548730},  // face 0
	{5.760339081714187279, 3.665943979320991689, 1.571548876927796127},  // face 1
	{0.780213654393430055, 4.969003859179821079, 2.874608756786625655},  // face 2
	{0.897199975119535024, 5.086990179905926144, 2.992595077512730473},  // face 3
	{6.130269123335111400, 4.035874020941915804, 1.941478918548720291},  // face 4
	{2.692877706530642877, 0.598482604137447119, 4.787272808923838195},  // face 5
	{2.982963003477243874, 0.888567901084048369, 5.077358105870439581},  // face 6
	{3.532912002790141181, 1.438516900396945656, 5.627307105183336758},  // face 7
	{3.494305004758171946, 1.399909902364976048, 5.588700107151367234},  // face 8
	{3.003214169499538391, 0.908819067106342928, 5.097609271892733906},  // face 9
	{5.930472956509811562, 3.836077854116615875, 1.741682751723420374},  // face 10
	{0.138378484090254847, 4.327168688876645809, 2.232773586483450307},  // face 11
	{0.448714947059150361, 4.637505151845541521, 2.543110049452346120},  // face 12
	{0.158629650112549365, 4.347419854898940135, 2.253024752505744869},  // face 13
	{5.891865957979862386, 3.797470855586666806, 1.703075753193471330},  // face 14
	{5.401775122721453543, 3.307380020328257965, 1.212984917935062374},  // face 15
	{5.611690825861418396, 3.517295723468222704, 1.422900621075027214},  // face 16
	{5.321605529792489535, 3.227210427399293958, 1.132815325006098423},  // face 17
	{5.471966437919222939, 3.377571335526027311, 1.283176233132831783},  // face 18
	{4.831696968776721443, 2.737301866383525914, 0.642906764990330379},  // face 19
}

// MaxDimByCIIres contains maximum dimension for unit coordinates by Class II resolution.
// This is used for scaling coordinates during transformations.
var MaxDimByCIIres = []float64{
	2.0,         // res 0
	-1.0,        // res 1 (not used - Class III)
	14.0,        // res 2
	-1.0,        // res 3 (not used - Class III)
	98.0,        // res 4
	-1.0,        // res 5 (not used - Class III)
	686.0,       // res 6
	-1.0,        // res 7 (not used - Class III)
	4802.0,      // res 8
	-1.0,        // res 9 (not used - Class III)
	33614.0,     // res 10
	-1.0,        // res 11 (not used - Class III)
	235298.0,    // res 12
	-1.0,        // res 13 (not used - Class III)
	1647086.0,   // res 14
	-1.0,        // res 15 (not used - Class III)
}

// MaxDimByCIIIres contains maximum dimension for unit coordinates by Class III resolution.
var MaxDimByCIIIres = []float64{
	-1.0,        // res 0 (not used - Class II)
	6.0,         // res 1
	-1.0,        // res 2 (not used - Class II)
	42.0,        // res 3
	-1.0,        // res 4 (not used - Class II)
	294.0,       // res 5
	-1.0,        // res 6 (not used - Class II)
	2058.0,      // res 7
	-1.0,        // res 8 (not used - Class II)
	14406.0,     // res 9
	-1.0,        // res 10 (not used - Class II)
	100842.0,    // res 11
	-1.0,        // res 12 (not used - Class II)
	705894.0,    // res 13
	-1.0,        // res 14 (not used - Class II)
	4941258.0,   // res 15
}

// UnitScaleByCIIres contains unit scale for each Class II resolution.
var UnitScaleByCIIres = []float64{
	1.0,                        // res 0
	-1.0,                       // res 1 (not used)
	1.0 / 7.0,                  // res 2
	-1.0,                       // res 3 (not used)
	1.0 / 49.0,                 // res 4
	-1.0,                       // res 5 (not used)
	1.0 / 343.0,                // res 6
	-1.0,                       // res 7 (not used)
	1.0 / 2401.0,               // res 8
	-1.0,                       // res 9 (not used)
	1.0 / 16807.0,              // res 10
	-1.0,                       // res 11 (not used)
	1.0 / 117649.0,             // res 12
	-1.0,                       // res 13 (not used)
	1.0 / 823543.0,             // res 14
	-1.0,                       // res 15 (not used)
}

// UnitScaleByCIIIres contains unit scale for each Class III resolution.
var UnitScaleByCIIIres = []float64{
	-1.0,                       // res 0 (not used)
	1.0 / 3.0,                  // res 1
	-1.0,                       // res 2 (not used)
	1.0 / 21.0,                 // res 3
	-1.0,                       // res 4 (not used)
	1.0 / 147.0,                // res 5
	-1.0,                       // res 6 (not used)
	1.0 / 1029.0,               // res 7
	-1.0,                       // res 8 (not used)
	1.0 / 7203.0,               // res 9
	-1.0,                       // res 10 (not used)
	1.0 / 50421.0,              // res 11
	-1.0,                       // res 12 (not used)
	1.0 / 352947.0,             // res 13
	-1.0,                       // res 14 (not used)
	1.0 / 2470629.0,            // res 15
}

// IsResolutionClassII returns true if the resolution is Class II (even).
func IsResolutionClassII(res int) bool {
	return res%2 == 0
}

// IsResolutionClassIII returns true if the resolution is Class III (odd).
func IsResolutionClassIII(res int) bool {
	return res%2 == 1
}

// GeoToClosestFace finds the closest icosahedron face to a lat/lng point.
// Returns the face number and squared distance.
func GeoToClosestFace(lat, lng float64) (face int, sqDist float64) {
	// Convert to 3D cartesian coordinates on unit sphere
	x := math.Cos(lat) * math.Cos(lng)
	y := math.Cos(lat) * math.Sin(lng)
	z := math.Sin(lat)
	
	minDist := math.MaxFloat64
	closestFace := 0
	
	// Check distance to each face center
	for f := 0; f < NumIcosaFaces; f++ {
		faceLat := FaceCenterGeo[f][0]
		faceLng := FaceCenterGeo[f][1]
		
		// Convert face center to 3D
		fx := math.Cos(faceLat) * math.Cos(faceLng)
		fy := math.Cos(faceLat) * math.Sin(faceLng)
		fz := math.Sin(faceLat)
		
		// Squared distance in 3D
		dx := x - fx
		dy := y - fy
		dz := z - fz
		dist := dx*dx + dy*dy + dz*dz
		
		if dist < minDist {
			minDist = dist
			closestFace = f
		}
	}
	
	return closestFace, minDist
}

// PosAngleRads normalizes radians to the range [0, 2π).
func PosAngleRads(rads float64) float64 {
	if rads < 0.0 {
		rads += 2 * math.Pi
	}
	if rads >= 2*math.Pi {
		rads -= 2 * math.Pi
	}
	return rads
}

// GeoAzimuthRads calculates the azimuth from p1 to p2 in radians.
func GeoAzimuthRads(p1Lat, p1Lng, p2Lat, p2Lng float64) float64 {
	return math.Atan2(math.Cos(p2Lat)*math.Sin(p2Lng-p1Lng),
		math.Cos(p1Lat)*math.Sin(p2Lat)-math.Sin(p1Lat)*math.Cos(p2Lat)*math.Cos(p2Lng-p1Lng))
}

// GeoToHex2d converts geographic coordinates to hex2d coordinates on a specific face.
func GeoToHex2d(lat, lng float64, res int) (face int, v coordijk.Vec2d) {
	// Determine the icosahedron face
	face, sqDist := GeoToClosestFace(lat, lng)
	
	// cos(r) = 1 - 2 * sin^2(r/2) = 1 - 2 * (sqd / 4) = 1 - sqd/2
	r := math.Acos(1 - sqDist*0.5)
	
	if r < 1e-12 { // EPSILON
		v = coordijk.Vec2d{X: 0.0, Y: 0.0}
		return face, v
	}
	
	// Find CCW theta from Class II i-axis
	faceLat := FaceCenterGeo[face][0]
	faceLng := FaceCenterGeo[face][1]
	
	azimuth := GeoAzimuthRads(faceLat, faceLng, lat, lng)
	theta := PosAngleRads(FaceAxesAzRadsCII[face][0] - PosAngleRads(azimuth))
	
	// Adjust theta for Class III (odd resolutions)
	if IsResolutionClassIII(res) {
		theta = PosAngleRads(theta - M_AP7_ROT_RADS)
	}
	
	// Perform gnomonic scaling of r
	r = math.Tan(r)
	
	// Scale for current resolution length u
	r *= 2.61803398874989588842 // INV_RES0_U_GNOMONIC
	for i := 0; i < res; i++ {
		r *= M_SQRT7
	}
	
	// Convert to local x,y
	v.X = r * math.Cos(theta)
	v.Y = r * math.Sin(theta)
	
	return face, v
}

// GeoToFaceIJK converts geographic coordinates to FaceIJK coordinates.
func GeoToFaceIJK(lat, lng float64, res int) FaceIJK {
	// First convert to hex2d
	face, v := GeoToHex2d(lat, lng, res)
	
	// Then convert to IJK
	coord := coordijk.Hex2dToCoordIJK(v)
	
	return FaceIJK{
		Face:  face,
		Coord: coord,
	}
}

// FaceIJKToH3 converts FaceIJK coordinates to an H3 index.
// This implements the H3 C _faceIjkToH3 function logic.
func FaceIJKToH3(fijk FaceIJK, res int) uint64 {
	// Initialize H3 index
	h := indexbits.H3_INIT
	h = indexbits.SetMode(h, 1) // H3_CELL_MODE
	h = indexbits.SetResolution(h, res)
	
	// Handle resolution 0 base cell case
	if res == 0 {
		// Check MAX_FACE_COORD bounds (should be 2 based on H3 C)
		if fijk.Coord.I > 2 || fijk.Coord.J > 2 || fijk.Coord.K > 2 {
			return 0 // H3_NULL - out of range
		}
		baseCell := faceIJKToBaseCell(fijk)
		if baseCell < 0 {
			return 0 // H3_NULL - invalid lookup
		}
		h = indexbits.SetBaseCell(h, baseCell)
		return h
	}
	
	// Make a copy - we'll scale this down to base cell level
	fijkBC := fijk
	ijk := &fijkBC.Coord
	
	// Build H3 index from finest resolution down to coarsest
	// This matches the H3 C algorithm exactly
	for r := res - 1; r >= 0; r-- {
		lastIJK := *ijk
		var lastCenter coordijk.CoordIJK
		
		if IsResolutionClassIII(r + 1) {
			// Class III: rotate ccw (UpAp7)
			ijk.UpAp7()
			lastCenter = *ijk
			lastCenter.DownAp7()
		} else {
			// Class II: rotate cw (UpAp7r)
			ijk.UpAp7r()
			lastCenter = *ijk  
			lastCenter.DownAp7r()
		}
		
		// Calculate difference and normalize
		diff := lastIJK.Sub(lastCenter)
		diff.Normalize()
		
		// Convert unit IJK to digit and set in H3 index
		digit := coordijk.UnitIJKToDigit(diff)
		if digit == coordijk.InvalidDigit {
			return 0 // H3_NULL - invalid digit
		}
		h = indexbits.SetDigit(h, r+1, int(digit))
	}
	
	// After scaling loop, coordinates should be in base cell range [0-2]
	// Look up base cell
	baseCell := faceIJKToBaseCell(fijkBC)
	if baseCell < 0 {
		return 0 // H3_NULL - invalid lookup
	}
	h = indexbits.SetBaseCell(h, baseCell)
	
	// Pentagon handling for special base cells
	if tables.IsPentagonBaseCell(baseCell) {
		// Check if leading non-zero digit matches the pentagon's orientation
		leadingNonZeroDigit := h3LeadingNonZeroDigit(h)
		
		// Pentagon base cells need special KAxesDigit rotation handling
		if leadingNonZeroDigit == coordijk.KAxesDigit {
			// Apply pentagon-specific 60-degree counter-clockwise rotation
			h = h3RotatePent60ccw(h)
		}
	}
	
	// Apply base cell rotation if needed
	baseCCWrot60 := faceIJKToBaseCellCCWrot60(fijkBC)
	if baseCCWrot60 > 0 {
		for i := 0; i < baseCCWrot60; i++ {
			h = h3Rotate60ccw(h)
		}
	}
	
	return h
}

// faceIJKToBaseCell looks up base cell from FaceIJK coordinates.
func faceIJKToBaseCell(fijk FaceIJK) int {
	if fijk.Face < 0 || fijk.Face >= NumIcosaFaces ||
		fijk.Coord.I < 0 || fijk.Coord.I > 2 ||
		fijk.Coord.J < 0 || fijk.Coord.J > 2 ||
		fijk.Coord.K < 0 || fijk.Coord.K > 2 {
		return -1 // Invalid coordinates
	}
	return tables.FaceIJKBaseCells[fijk.Face][fijk.Coord.I][fijk.Coord.J][fijk.Coord.K].BaseCell
}

// faceIJKToBaseCellCCWrot60 looks up rotation from FaceIJK coordinates.
func faceIJKToBaseCellCCWrot60(fijk FaceIJK) int {
	if fijk.Face < 0 || fijk.Face >= NumIcosaFaces ||
		fijk.Coord.I < 0 || fijk.Coord.I > 2 ||
		fijk.Coord.J < 0 || fijk.Coord.J > 2 ||
		fijk.Coord.K < 0 || fijk.Coord.K > 2 {
		return -1 // Invalid coordinates
	}
	return tables.FaceIJKBaseCells[fijk.Face][fijk.Coord.I][fijk.Coord.J][fijk.Coord.K].CCWRot60
}

// h3LeadingNonZeroDigit returns the leading non-zero digit in an H3 index
func h3LeadingNonZeroDigit(h uint64) coordijk.Direction {
	resolution := indexbits.GetResolution(h)
	for r := 1; r <= resolution; r++ {
		digit := indexbits.GetDigit(h, r)
		if digit != 0 {
			return coordijk.Direction(digit)
		}
	}
	return coordijk.CenterDigit // All digits are zero
}

// h3Rotate60cw rotates an H3 index 60 degrees clockwise
func h3Rotate60cw(h uint64) uint64 {
	resolution := indexbits.GetResolution(h)
	rotated := h
	
	for r := 1; r <= resolution; r++ {
		digit := indexbits.GetDigit(h, r)
		// Rotate digit 60 degrees clockwise: (d + 1) % 6, but skip 0
		if digit != 0 {
			newDigit := (digit % 6) + 1
			rotated = indexbits.SetDigit(rotated, r, newDigit)
		}
	}
	return rotated
}

// h3Rotate60ccw rotates an H3 index 60 degrees counter-clockwise
func h3Rotate60ccw(h uint64) uint64 {
	resolution := indexbits.GetResolution(h)
	rotated := h
	
	for r := 1; r <= resolution; r++ {
		digit := indexbits.GetDigit(h, r)
		// Rotate digit 60 degrees counter-clockwise: (d - 1), wrapping from 1 to 6
		if digit != 0 {
			var newDigit int
			if digit == 1 {
				newDigit = 6
			} else {
				newDigit = digit - 1
			}
			rotated = indexbits.SetDigit(rotated, r, newDigit)
		}
	}
	return rotated
}

// h3RotatePent60ccw applies pentagon-specific 60-degree counter-clockwise rotation
// Pentagon cells have special handling for the KAxesDigit orientation
func h3RotatePent60ccw(h uint64) uint64 {
	baseCell := indexbits.GetBaseCell(h)
	
	// Verify this is actually a pentagon base cell
	if !tables.IsPentagonBaseCell(baseCell) {
		return h // No rotation for non-pentagon cells
	}
	
	// Get pentagon-specific rotation offset from tables
	pentOffset := tables.BaseCells[baseCell].CWOffsetPent
	
	// Apply pentagon rotation logic - this is a simplified implementation
	// TODO: Full pentagon rotation requires more complex digit manipulation
	// For now, apply standard counter-clockwise rotation with pentagon constraints
	rotated := h3Rotate60ccw(h)
	
	// Pentagon cells may need additional digit adjustments
	// This depends on the specific pentagon's CWOffsetPent values
	if pentOffset[0] != -1 || pentOffset[1] != -1 {
		// Pentagon has specific rotation offsets - apply them
		// This is a placeholder for more complex pentagon rotation logic
		// that would involve digit pattern adjustments specific to each pentagon
	}
	
	return rotated
}