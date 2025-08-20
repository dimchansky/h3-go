// Package faceijk implements face-centered IJK coordinate transformations
// for the H3 icosahedral projection.
package faceijk

import (
	"math"
	
	"github.com/dimchansky/h3-go/internal/coordijk"
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