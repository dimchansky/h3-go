package h3

import (
	"math"
	"slices"
)

// MaxPolygonToCellsSize returns an upper bound on the number of cells
// PolygonToCells produces for the polygon at the given resolution, for
// pre-sizing AppendPolygonToCells destination buffers.
//
// H3 C API: maxPolygonToCellsSize.
func MaxPolygonToCellsSize(p GeoPolygon, res int) (int64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var sz int64
	if errC := maxPolygonToCellsSize(&p, int32(res), 0, &sz); errC != eSuccess {
		return 0, toErr(errC)
	}
	return sz, nil
}

// PolygonToCells returns all cells at the given resolution whose center point
// is contained in the polygon (holes excluded), in no particular order.
//
// H3 C API: polygonToCells.
func PolygonToCells(p GeoPolygon, res int) ([]Cell, error) {
	return AppendPolygonToCells(nil, p, res)
}

// AppendPolygonToCells appends the cells whose center point is contained in
// the polygon to dst and returns the extended slice. Pass dst[:0] (or nil)
// to reuse dst's capacity.
//
// H3 C API: polygonToCells.
func AppendPolygonToCells(dst []Cell, p GeoPolygon, res int) ([]Cell, error) {
	sz, err := MaxPolygonToCellsSize(p, res)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz)) // the algorithm probes the out buffer
	if errC := polygonToCells(&p, int32(res), 0, win); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst[:start+compactNonNull(win)], nil
}

// MaxPolygonToCellsSizeExperimental returns an upper bound on the number of
// cells PolygonToCellsExperimental produces for the polygon, resolution, and
// containment mode.
//
// Like its C counterpart, this API is experimental and may change in minor
// versions.
//
// H3 C API: maxPolygonToCellsSizeExperimental.
func MaxPolygonToCellsSizeExperimental(p GeoPolygon, res int, mode ContainmentMode) (int64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	sz, errC := maxPolygonToCellsSizeExperimental(&p, int32(res), uint32(mode))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return sz, nil
}

// PolygonToCellsExperimental returns all cells at the given resolution that
// match the polygon under the given containment mode, in no particular
// order.
//
// Like its C counterpart, this API is experimental and may change in minor
// versions.
//
// H3 C API: polygonToCellsExperimental.
func PolygonToCellsExperimental(p GeoPolygon, res int, mode ContainmentMode) ([]Cell, error) {
	return AppendPolygonToCellsExperimental(nil, p, res, mode)
}

// AppendPolygonToCellsExperimental appends the cells matching the polygon
// under the given containment mode to dst and returns the extended slice.
//
// H3 C API: polygonToCellsExperimental.
func AppendPolygonToCellsExperimental(dst []Cell, p GeoPolygon, res int, mode ContainmentMode) ([]Cell, error) {
	sz, err := MaxPolygonToCellsSizeExperimental(p, res, mode)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz))
	if errC := polygonToCellsExperimental(&p, int32(res), uint32(mode), sz, win); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst[:start+compactNonNull(win)], nil
}

// CellsToMultiPolygon returns the multipolygon (one GeoPolygon per contiguous
// region, holes included) describing the outline of the given set of
// same-resolution cells. Vertices are cell-boundary vertices; an empty input
// yields a nil result.
//
// H3 C API: cellsToLinkedMultiPolygon (the linked-list output is converted
// to slice-based GeoPolygon values; C's destroyLinkedMultiPolygon is
// unnecessary under garbage collection).
func CellsToMultiPolygon(cells []Cell) ([]GeoPolygon, error) {
	if len(cells) == 0 {
		return nil, nil
	}
	if len(cells) > math.MaxInt32 {
		return nil, ErrDomain
	}
	var linked linkedGeoPolygon
	if errC := cellsToLinkedMultiPolygon(cells, int32(len(cells)), &linked); errC != eSuccess {
		return nil, toErr(errC)
	}
	out := make([]GeoPolygon, 0, countLinkedPolygons(&linked))
	for poly := &linked; poly != nil; poly = poly.Next {
		var outer GeoLoop
		var holes []GeoLoop
		if n := countLinkedLoops(poly); n > 1 {
			holes = make([]GeoLoop, 0, n-1)
		}
		for loop := poly.First; loop != nil; loop = loop.Next {
			coords := make(GeoLoop, 0, countLinkedCoords(loop))
			for c := loop.First; c != nil; c = c.Next {
				coords = append(coords, c.Vertex)
			}
			// normalizeMultiPolygon guarantees the first loop is the outer one.
			if outer == nil {
				outer = coords
			} else {
				holes = append(holes, coords)
			}
		}
		out = append(out, GeoPolygon{GeoLoop: outer, Holes: holes})
	}
	return out, nil
}

// CompactCells returns the minimal set of cells of coarser resolutions that
// exactly covers the given set of cells. Input cells must have the same
// resolution and must not contain duplicates. Behavior for inputs violating
// these preconditions is not guaranteed; some such inputs may return
// ErrDuplicateInput or ErrResolutionMismatch. The order of the output is
// not guaranteed.
//
// H3 C API: compactCells.
func CompactCells(cells []Cell) ([]Cell, error) { return AppendCompactCells(nil, cells) }

// AppendCompactCells appends the compacted equivalent of cells to dst and
// returns the extended slice; see CompactCells for the input preconditions.
// It reuses the destination buffer but may allocate internal working
// storage for nontrivial inputs; the amount depends on the input and
// compaction depth (mirroring the C implementation's heap allocations).
// dst must not share memory with cells — results are unspecified when they
// overlap. On error the returned slice has dst's original length and
// elements.
//
// H3 C API: compactCells.
func AppendCompactCells(dst, cells []Cell) ([]Cell, error) {
	start := len(dst)
	dst, win := growZeroed(dst, len(cells))
	if errC := compactCells(cells, win, int64(len(cells))); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst[:start+compactNonNull(win)], nil
}

// UncompactCellsSize returns the number of cells UncompactCells produces
// for the given compacted set at the given resolution. Any input cell finer
// than res fails with ErrResolutionMismatch; res outside 0..MaxResolution
// fails with ErrResolutionDomain.
//
// H3 C API: uncompactCellsSize.
func UncompactCellsSize(cells []Cell, res int) (int64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	n, errC := uncompactCellsSize(cells, int64(len(cells)), int32(res))
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return n, nil
}

// UncompactCells expands a compacted set of cells to the given resolution.
// H3_NULL (zero) entries in the input are skipped, mirroring C. Any input
// cell finer than res fails with ErrResolutionMismatch; res outside
// 0..MaxResolution fails with ErrResolutionDomain.
//
// H3 C API: uncompactCells.
func UncompactCells(cells []Cell, res int) ([]Cell, error) {
	return AppendUncompactCells(nil, cells, res)
}

// AppendUncompactCells appends the expansion of the compacted set at the
// given resolution to dst and returns the extended slice; see
// UncompactCells for the error conditions. dst must not share memory with
// cells — the expansion is written while the input is read, so results are
// unspecified when they overlap. On error the returned slice has dst's
// original length and elements.
//
// H3 C API: uncompactCells.
func AppendUncompactCells(dst, cells []Cell, res int) ([]Cell, error) {
	n, err := UncompactCellsSize(cells, res)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst = slices.Grow(dst, int(n))[:start+int(n)]
	if errC := uncompactCells(cells, int64(len(cells)), dst[start:], n, int32(res)); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst, nil
}
