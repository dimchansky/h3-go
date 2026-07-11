//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_describeH3Error_ParityWithC(t *testing.T) {
	testCases := []struct {
		err  h3Error
		desc string
	}{
		// Valid error codes
		{eSuccess, "Success"},
		{eFailed, "The operation failed but a more specific error is not available"},
		{eDomain, "Argument was outside of acceptable range"},
		{eLatlngDomain, "Latitude or longitude arguments were outside of acceptable range"},
		{eResDomain, "Resolution argument was outside of acceptable range"},
		{eCellInvalid, "Cell argument was not valid"},
		{eDirEdgeInvalid, "Directed edge argument was not valid"},
		{eUndirEdgeInvalid, "Undirected edge argument was not valid"},
		{eVertexInvalid, "Vertex argument was not valid"},
		{ePentagon, "Pentagon distortion was encountered"},
		{eDuplicateInput, "Duplicate input"},
		{eNotNeighbors, "Cell arguments were not neighbors"},
		{eResMismatch, "Cell arguments had incompatible resolutions"},
		{eMemoryAlloc, "Memory allocation failed"},
		{eMemoryBounds, "Bounds of provided memory were insufficient"},
		{eOptionInvalid, "Mode or flags argument was not valid"},
	}

	for _, tc := range testCases {
		goResult := describeH3Error(tc.err)
		cResult := describeH3ErrorC(uint32(tc.err))

		if goResult != cResult {
			t.Fatalf("describeH3Error mismatch for err=%d: go='%s' c='%s'", tc.err, goResult, cResult)
		}

		// Also verify against expected string
		if goResult != tc.desc {
			t.Fatalf("describeH3Error wrong description for err=%d: got='%s' expected='%s'", tc.err, goResult, tc.desc)
		}
	}
}

func Test_h3index_describeH3Error_InvalidCodes_ParityWithC(t *testing.T) {
	invalidCodes := []h3Error{
		h3Error(16),  // Just above valid range
		h3Error(100), // Well above valid range
		h3Error(255), // Edge case
	}

	for _, code := range invalidCodes {
		goResult := describeH3Error(code)
		cResult := describeH3ErrorC(uint32(code))

		if goResult != cResult {
			t.Fatalf("describeH3Error mismatch for invalid code=%d: go='%s' c='%s'", code, goResult, cResult)
		}

		// Should be the invalid error message
		expectedMsg := "Invalid error code"
		if goResult != expectedMsg {
			t.Fatalf("describeH3Error wrong message for invalid code=%d: got='%s' expected='%s'", code, goResult, expectedMsg)
		}
	}
}

func Test_h3index_describeH3Error_BoundaryValues_ParityWithC(t *testing.T) {
	boundaryCases := []h3Error{
		0,  // Minimum valid
		15, // Maximum valid
	}

	for _, code := range boundaryCases {
		goResult := describeH3Error(code)
		cResult := describeH3ErrorC(uint32(code))

		if goResult != cResult {
			t.Fatalf("describeH3Error mismatch for boundary code=%d: go='%s' c='%s'", code, goResult, cResult)
		}

		// Should not be the invalid error message
		if goResult == "Invalid error code" {
			t.Fatalf("describeH3Error incorrectly returned invalid message for valid code=%d", code)
		}
	}
}
