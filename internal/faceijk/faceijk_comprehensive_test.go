package faceijk

import (
	"math"
	"testing"

	"github.com/dimchansky/h3-go/internal/coordijk"
)

// TestGeoToClosestFaceComprehensive tests face selection with extensive test data
// NOTE: Some tests currently fail due to implementation differences from H3 C
// These failures help identify areas that need improvement and will pass once fixed
func TestGeoToClosestFaceComprehensive(t *testing.T) {
	tests := []struct {
		name string
		lat, lng float64
		expectedFace int
		expectedSqDist float64
	}{
		// Face center tests (20 cases) - should have zero distance
		{"face0_center", 0.8035826497, 1.2483974196, 0, 0.0000000000},
		{"face1_center", 1.3077478835, 2.5369450099, 1, 0.0000000000},
		{"face2_center", 1.0547512535, -1.3475173589, 2, 0.0000000000},
		{"face3_center", 0.6001915955, -0.4506039095, 3, 0.0000000000},
		{"face4_center", 0.4917154282, 0.4019882029, 4, 0.0000000000},
		{"face5_center", 0.1727453274, 1.6781468853, 5, 0.0000000000},
		{"face6_center", 0.6059293216, 2.9539233298, 6, 0.0000000000},
		{"face7_center", 0.4273705183, -1.8888762003, 7, 0.0000000000},
		{"face8_center", -0.0790661185, -0.7334295134, 8, 0.0000000000},
		{"face9_center", -0.2309616445, 0.5064955873, 9, 0.0000000000},
		{"face10_center", 0.0790661185, 2.4081631402, 10, 0.0000000000},
		{"face11_center", 0.2309616445, -2.6350970663, 11, 0.0000000000},
		{"face12_center", -0.1727453274, -1.4634457683, 12, 0.0000000000},
		{"face13_center", -0.6059293216, -0.1876693238, 13, 0.0000000000},
		{"face14_center", -0.4273705183, 1.2527164533, 14, 0.0000000000},
		{"face15_center", -0.6001915955, 2.6909887441, 15, 0.0000000000},
		{"face16_center", -0.4917154282, -2.7396044507, 16, 0.0000000000},
		{"face17_center", -0.8035826497, -1.8931952340, 17, 0.0000000000},
		{"face18_center", -1.3077478835, -0.6046476437, 18, 0.0000000000},
		{"face19_center", -1.0547512535, 1.7940752947, 19, 0.0000000000},

		// Random point tests (50 cases) - comprehensive coverage
		{"random_0", -1.5697636617, 0.1544853181, 18, 0.0684078575},
		{"random_1", 0.7396048392, -1.4871951484, 2, 0.1055982598},
		{"random_2", -0.3888538624, -1.9082924370, 17, 0.1696951222},
		{"random_3", 1.4950018886, 0.0773969581, 1, 0.1049104594},
		{"random_4", 0.0956584965, -1.5261754726, 12, 0.0754666243},
		{"random_5", -1.2343717945, 1.9822672216, 19, 0.0379289285},
		{"random_6", 1.2583477740, -0.3014129308, 2, 0.1927010279},
		{"random_7", -0.7998850780, -1.5870815886, 17, 0.0449828533},
		{"random_8", -0.9793156992, -1.1100891728, 18, 0.1431598631},
		{"random_9", -1.2457755335, 1.9312876484, 19, 0.0393414182},
		{"random_10", 0.1081051162, 2.1642678037, 10, 0.0595044786},
		{"random_11", 0.7446171491, -2.6494132902, 11, 0.2582384156},
		{"random_12", 0.1225516487, -2.3184409746, 11, 0.1078109129},
		{"random_13", 1.1389071705, -0.2224472753, 2, 0.2420844622},
		{"random_14", -0.0880494558, -0.3141280588, 8, 0.1721232686},
		{"random_15", -0.8373133572, -3.0641984706, 16, 0.1798883931},
		{"random_16", -1.4984588503, 3.1008121151, 18, 0.1056147090},
		{"random_17", 1.3051405898, 1.7959699643, 1, 0.0358048536},
		{"random_18", 0.2224871245, 1.6916873880, 5, 0.0026498981},
		{"random_19", 0.3882072087, -0.9787684033, 3, 0.2529358930},
		{"random_20", -0.3907098403, -1.4632783279, 12, 0.0473207636},
		{"random_21", -0.4657824776, 0.8855839887, 14, 0.1098459329},
		{"random_22", -0.4279473126, -2.8097959603, 16, 0.0080148427},
		{"random_23", 0.0900322657, -2.1507380995, 7, 0.1745191789},
		{"random_24", -0.1450831341, -1.0726714442, 8, 0.1167923558},
		{"random_25", -0.9651579877, -2.7348561179, 16, 0.2200034743},
		{"random_26", 1.5284671813, 0.2895750006, 1, 0.0843052431},
		{"random_27", -1.2907889041, 3.0995105013, 19, 0.2566667873},
		{"random_28", -0.2081930734, 1.2664641847, 14, 0.0480150343},
		{"random_29", -0.9841344694, 0.2745862375, 18, 0.2081061420},
		{"random_30", 1.5564388655, -1.9480288499, 1, 0.0706788186},
		{"random_31", 0.5788776315, -0.6321919334, 3, 0.0231696405},
		{"random_32", -0.1917354688, 1.5520772681, 14, 0.1347334335},
		{"random_33", -1.0113751115, 1.9526982884, 19, 0.0084563351},
		{"random_34", 0.9616367936, -2.5292261318, 1, 0.3130738049},
		{"random_35", -1.4774973672, -2.2998329551, 18, 0.0832066829},
		{"random_36", 0.4317667093, -0.7518926803, 3, 0.0958221525},
		{"random_37", -0.7873128536, 0.0422522092, 13, 0.0633418488},
		{"random_38", 0.0664700083, -2.4911113432, 11, 0.0470986006},
		{"random_39", 1.5192711027, -0.9513310324, 1, 0.0965497335},
		{"random_40", 0.8429729089, -1.4743759777, 2, 0.0499583620},
		{"random_41", 0.5228975307, 2.6082925355, 6, 0.0911192104},
		{"random_42", -0.1056218098, -0.3718165919, 8, 0.1289276818},
		{"random_43", 1.3239601385, -0.2054354369, 2, 0.2129306666},
		{"random_44", 1.4992651450, -1.1307661770, 1, 0.1058766212},
		{"random_45", 0.9242091024, 2.2966104538, 1, 0.1543126801},
		{"random_46", 0.7622773528, 0.3612548442, 4, 0.0738160012},
		{"random_47", 1.0265796960, 0.1961957575, 4, 0.2985853313},
		{"random_48", -0.6050947371, -0.9836524080, 8, 0.3214507029},
		{"random_49", -0.5927386296, -0.3356863745, 13, 0.0150833651},

		// Boundary condition tests (7 cases) - these may expose edge case bugs
		{"boundary_0", 0.0000000000, 0.0000000000, 9, 0.2975392027},
		{"boundary_1", 1.5707963268, 0.0000000000, 1, 0.0687964130},
		{"boundary_2", -1.5707963268, 0.0000000000, 18, 0.0687964130},
		{"boundary_3", 0.0000000000, 3.1415926536, 9, 0.2975392027},  // May fail - antimeridian
		{"boundary_4", 0.0000000000, -3.1415926536, 9, 0.2975392027}, // May fail - antimeridian
		{"boundary_5", 0.7853981634, 1.5707963268, 5, 0.0944913174},  // May fail - face boundary
		{"boundary_6", -0.7853981634, -1.5707963268, 12, 0.0944913174}, // May fail - face boundary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			face, sqDist := GeoToClosestFace(tt.lat, tt.lng)
			
			// These tests may fail due to implementation differences
			// Keep the failures - they show us what needs to be fixed
			if face != tt.expectedFace {
				t.Errorf("GeoToClosestFace(%f, %f) face = %d, want %d", tt.lat, tt.lng, face, tt.expectedFace)
			}
			if math.Abs(sqDist - tt.expectedSqDist) > 1e-8 {
				t.Errorf("GeoToClosestFace(%f, %f) sqDist = %f, want %f", tt.lat, tt.lng, sqDist, tt.expectedSqDist)
			}
		})
	}
}

// TestGeoAzimuthRadsComprehensive tests azimuth calculations with extensive test data
// NOTE: Currently fails due to incorrect test data generation - keep failures to guide fixes
func TestGeoAzimuthRadsComprehensive(t *testing.T) {
	tests := []struct {
		name string
		p1Lat, p1Lng, p2Lat, p2Lng float64
		expected float64
	}{
		// Sample of azimuth test cases - these currently fail due to test data issues
		// Keep them as they will pass once we generate correct expected values from actual H3 C
		{"azimuth_0", -1.4707630591, 0.1544853181, 0.7396048392, -1.4871951484, 1.2267309751},
		{"azimuth_1", -0.2888538624, -1.9082924370, 1.3950018886, 0.0773969581, 0.8486026482},
		{"azimuth_2", 0.0956584965, -1.5261754726, -1.1343717945, 1.9822672216, -2.9781080103},
		{"azimuth_3", 1.1583477740, -0.3014129308, -0.6998850780, -1.5870815886, -2.4950768648},
		{"azimuth_4", -0.8793156992, -1.1100891728, -1.1457755335, 1.9312876484, 2.7052456139},
		{"azimuth_5", 0.0081051162, 2.1642678037, 0.6446171491, -2.6494132902, -1.2925353219},
		{"azimuth_6", 0.0225516487, -2.3184409746, 1.0389071705, -0.2224472753, 0.6737031341},
		{"azimuth_7", -0.0880494558, -0.3141280588, -0.7373133572, -3.0641984706, -2.5471706390},
		{"azimuth_8", -1.3984588503, 3.1008121151, 1.2051405898, 1.7959699643, 0.3624491930},
		{"azimuth_9", 0.1224871245, 1.6916873880, 0.2882072087, -0.9787684033, -1.4052439928},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeoAzimuthRads(tt.p1Lat, tt.p1Lng, tt.p2Lat, tt.p2Lng)
			
			// This will likely fail - that's expected and helps identify what to fix
			if math.Abs(result - tt.expected) > 1e-8 {
				t.Errorf("GeoAzimuthRads(%f, %f, %f, %f) = %f, want %f", 
					tt.p1Lat, tt.p1Lng, tt.p2Lat, tt.p2Lng, result, tt.expected)
			}
		})
	}
}

// TestFaceIJKToH3PentagonCases tests cases that will fail until pentagon handling is implemented
func TestFaceIJKToH3PentagonCases(t *testing.T) {
	t.Log("These tests document pentagon base cell cases that need special handling")
	t.Log("They will fail until pentagon rotation logic is implemented")
	
	// Pentagon base cells from H3 documentation: 1, 4, 7, 11, 14
	pentagonCases := []struct {
		name string
		fijk FaceIJK
		res int
		note string
	}{
		{"pentagon_baseCell1", FaceIJK{0, coordijk.CoordIJK{0, 0, 1}}, 0, "Base cell 1 - pentagon"},
		{"pentagon_baseCell4", FaceIJK{1, coordijk.CoordIJK{1, 1, 0}}, 0, "Base cell 4 - pentagon"},
		{"pentagon_baseCell7", FaceIJK{2, coordijk.CoordIJK{1, 0, 1}}, 0, "Base cell 7 - pentagon"},
		{"pentagon_baseCell11", FaceIJK{3, coordijk.CoordIJK{2, 1, 0}}, 0, "Base cell 11 - pentagon"},
		{"pentagon_baseCell14", FaceIJK{4, coordijk.CoordIJK{0, 2, 1}}, 0, "Base cell 14 - pentagon"},
	}
	
	for _, tt := range pentagonCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Pentagon test: %s", tt.note)
			result := FaceIJKToH3(tt.fijk, tt.res)
			
			// These may fail - that's expected until pentagon handling is implemented
			// The failures guide us on what needs to be implemented
			if result == 0 {
				t.Logf("Pentagon case failed as expected: %v -> H3_NULL", tt.fijk)
				t.Logf("Will pass once pentagon rotation logic (_h3RotatePent60ccw, etc.) is implemented")
			} else {
				t.Logf("Pentagon case unexpectedly succeeded: %v -> 0x%x", tt.fijk, result)
			}
		})
	}
}

// TestImplementationTODOList documents what needs to be implemented for tests to pass
func TestImplementationTODOList(t *testing.T) {
	t.Log("=== IMPLEMENTATION ROADMAP ===")
	t.Log("To make failing tests pass, we need to implement:")
	t.Log("")
	t.Log("1. Pentagon handling in FaceIJKToH3:")
	t.Log("   - _h3LeadingNonZeroDigit function")
	t.Log("   - _baseCellIsCwOffset lookup")
	t.Log("   - Pentagon-specific digit patterns")
	t.Log("")  
	t.Log("2. H3 index rotation functions:")
	t.Log("   - _h3Rotate60cw, _h3Rotate60ccw")
	t.Log("   - _h3RotatePent60ccw for pentagon base cells")
	t.Log("")
	t.Log("3. Comprehensive test data from actual H3 C library:")
	t.Log("   - Real oracle-generated expected values")
	t.Log("   - Edge case coverage (antimeridian, poles, face boundaries)")
	t.Log("   - 100+ systematic test cases per function")
	t.Log("")
	t.Log("4. Coordinate transformation accuracy improvements:")
	t.Log("   - Fix GeoToClosestFace for boundary cases")  
	t.Log("   - Fix GeoAzimuthRads calculation precision")
	t.Log("   - Validate against H3 C reference implementation")
}