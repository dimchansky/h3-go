package h3

// H3Error represents H3 operation result codes.
type H3Error uint32

// H3Error codes (from h3api.h)
const (
	E_SUCCESS            H3Error = 0
	E_FAILED             H3Error = 1
	E_DOMAIN             H3Error = 2
	E_LATLNG_DOMAIN      H3Error = 3
	E_RES_DOMAIN         H3Error = 4
	E_CELL_INVALID       H3Error = 5
	E_DIR_EDGE_INVALID   H3Error = 6
	E_UNDIR_EDGE_INVALID H3Error = 7
	E_VERTEX_INVALID     H3Error = 8
	E_PENTAGON           H3Error = 9
	E_DUPLICATE_INPUT    H3Error = 10
	E_NOT_NEIGHBORS      H3Error = 11
	E_RES_MISMATCH       H3Error = 12
	E_MEMORY_ALLOC       H3Error = 13
	E_MEMORY_BOUNDS      H3Error = 14
	E_OPTION_INVALID     H3Error = 15
)
