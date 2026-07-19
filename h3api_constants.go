package h3

// h3Error represents H3 operation result codes.
type h3Error uint32

// h3Error codes (from h3api.h).
const (
	eSuccess          h3Error = 0
	eFailed           h3Error = 1
	eDomain           h3Error = 2
	eLatlngDomain     h3Error = 3
	eResDomain        h3Error = 4
	eCellInvalid      h3Error = 5
	eDirEdgeInvalid   h3Error = 6
	eUndirEdgeInvalid h3Error = 7
	eVertexInvalid    h3Error = 8
	ePentagon         h3Error = 9
	eDuplicateInput   h3Error = 10
	eNotNeighbors     h3Error = 11
	eResMismatch      h3Error = 12
	eMemoryAlloc      h3Error = 13
	eMemoryBounds     h3Error = 14
	eOptionInvalid    h3Error = 15
	eIndexInvalid     h3Error = 16
	eBaseCellDomain   h3Error = 17
	eDigitDomain      h3Error = 18
	eDeletedDigit     h3Error = 19

	// h3ErrorEnd is a sentinel: one past the last valid code.
	// Ported from H3 C: h3api.h::H3_ERROR_END.
	h3ErrorEnd h3Error = 20
)

// H3 version constants.
// These constants mirror the H3_VERSION_* macros from h3api.h
// Ported from H3 C: h3api.h::H3_VERSION_MAJOR, H3_VERSION_MINOR, H3_VERSION_PATCH.
const (
	h3VersionMajor = 4
	h3VersionMinor = 5
	h3VersionPatch = 0
)
