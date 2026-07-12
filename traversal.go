package h3

import (
	"math"
	"slices"
)

// checkK validates a public grid-distance argument before it is narrowed to
// the C layer's int32.
func checkK(k int) error {
	if k < 0 || k > math.MaxInt32 {
		return ErrDomain
	}
	return nil
}

// compactNonNull moves all non-null indexes in win to its front, preserving
// their relative order, and returns their count. It is used to prune the
// H3_NULL holes that hash-set-placed C outputs may contain.
func compactNonNull(win []Cell) int {
	n := 0
	for _, c := range win {
		if c != h3Null {
			win[n] = c
			n++
		}
	}
	return n
}

// growZeroed extends dst by n zeroed elements and returns the extended slice
// plus the newly added window. The gridDisk family uses its output buffer as
// a linear-probing hash set, so the window must be zeroed even when reused.
func growZeroed(dst []Cell, n int) ([]Cell, []Cell) {
	start := len(dst)
	dst = slices.Grow(dst, n)[:start+n]
	win := dst[start:]
	clear(win)
	return dst, win
}

// MaxGridDiskSize returns the maximum number of cells in a grid disk of
// radius k, for pre-sizing AppendGridDisk destination buffers.
//
// H3 C API: maxGridDiskSize.
func MaxGridDiskSize(k int) (int64, error) {
	if err := checkK(k); err != nil {
		return 0, err
	}
	var sz int64
	if errC := maxGridDiskSize(int32(k), &sz); errC != eSuccess {
		return 0, toErr(errC)
	}
	return sz, nil
}

// MaxGridRingSize returns the maximum number of cells in a hollow grid ring
// of radius k, for pre-sizing AppendGridRing destination buffers.
//
// H3 C API: maxGridRingSize.
func MaxGridRingSize(k int) (int64, error) {
	if err := checkK(k); err != nil {
		return 0, err
	}
	var sz int64
	if errC := _maxGridRingSize(int32(k), &sz); errC != eSuccess {
		return 0, toErr(errC)
	}
	return sz, nil
}

// GridDisk returns all cells within grid distance k of c (including c), in
// no particular order.
//
// H3 C API: gridDisk.
func (c Cell) GridDisk(k int) ([]Cell, error) { return c.AppendGridDisk(nil, k) }

// AppendGridDisk appends all cells within grid distance k of c to dst and
// returns the extended slice, in no particular order. Pass dst[:0] (or nil)
// to reuse dst's capacity; when capacity suffices the only allocation is the
// algorithm's internal distance scratch (as in H3 C).
//
// H3 C API: gridDisk.
func (c Cell) AppendGridDisk(dst []Cell, k int) ([]Cell, error) {
	sz, err := MaxGridDiskSize(k)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz))
	if errC := gridDisk(c, int32(k), win); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst[:start+compactNonNull(win)], nil
}

// GridDiskUnsafe returns all cells within grid distance k of c in ring-walk
// order (origin first, then each ring counterclockwise). It fails with
// ErrPentagon if a pentagon or pentagon distortion is encountered; use
// GridDisk for the safe variant.
//
// H3 C API: gridDiskUnsafe.
func (c Cell) GridDiskUnsafe(k int) ([]Cell, error) { return c.AppendGridDiskUnsafe(nil, k) }

// AppendGridDiskUnsafe appends the ring-walk-ordered grid disk of radius k
// to dst; see GridDiskUnsafe.
//
// H3 C API: gridDiskUnsafe.
func (c Cell) AppendGridDiskUnsafe(dst []Cell, k int) ([]Cell, error) {
	sz, err := MaxGridDiskSize(k)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz))
	if errC := gridDiskUnsafe(c, int32(k), win); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst, nil
}

// GridDiskDistances returns all cells within grid distance k of c along with
// each cell's grid distance from c, in no particular order. The distances
// slice is int32 to match the H3 C representation without a conversion copy.
//
// H3 C API: gridDiskDistances.
func (c Cell) GridDiskDistances(k int) ([]Cell, []int32, error) {
	return c.AppendGridDiskDistances(nil, nil, k)
}

// GridDiskDistancesGrouped returns cells within grid distance k of c grouped
// by distance: result[d] contains exactly the cells at distance d. Empty
// rings are retained, H3_NULL holes near pentagons are omitted, and no order
// within a ring is guaranteed. All rings share one backing cell array and
// have their capacity limited to their length, so appending to one ring
// cannot overwrite another.
//
// The flat GridDiskDistances and AppendGridDiskDistances forms remain the
// allocation-efficient choices when grouped slices are unnecessary.
//
// H3 C API: gridDiskDistances.
func (c Cell) GridDiskDistancesGrouped(k int) ([][]Cell, error) {
	cells, distances, err := c.GridDiskDistances(k)
	if err != nil {
		return nil, err
	}

	const stackRings = 64
	var countBuf, nextBuf [stackRings]int
	var counts, next []int
	if k < stackRings {
		counts = countBuf[:k+1]
		next = nextBuf[:k+1]
	} else {
		counts = make([]int, k+1)
		next = make([]int, k+1)
	}
	for _, distance := range distances {
		counts[int(distance)]++
	}
	for d := 1; d <= k; d++ {
		next[d] = next[d-1] + counts[d-1]
	}

	// In-place counting sort of the parallel cell/distance slices. This keeps
	// the flat result's cell allocation as the backing store for every ring.
	start := 0
	for d := 0; d <= k; d++ {
		end := start + counts[d]
		for next[d] < end {
			i := next[d]
			actual := int(distances[i])
			if actual == d {
				next[d]++
				continue
			}
			j := next[actual]
			cells[i], cells[j] = cells[j], cells[i]
			distances[i], distances[j] = distances[j], distances[i]
			next[actual]++
		}
		start = end
	}

	rings := make([][]Cell, k+1)
	start = 0
	for d := range rings {
		end := start + counts[d]
		rings[d] = cells[start:end:end]
		start = end
	}
	return rings, nil
}

// AppendGridDiskDistances appends the cells within grid distance k of c to
// dst and their distances to dstDist, returning both extended slices. When
// both capacities suffice the call does not allocate.
//
// H3 C API: gridDiskDistances.
func (c Cell) AppendGridDiskDistances(dst []Cell, dstDist []int32, k int) ([]Cell, []int32, error) {
	return c.appendGridDiskDistances(dst, dstDist, k, gridDiskDistances, true)
}

// GridDiskDistancesSafe is the always-correct but slower variant of
// GridDiskDistances (no unsafe-algorithm fast path).
//
// H3 C API: gridDiskDistancesSafe.
func (c Cell) GridDiskDistancesSafe(k int) ([]Cell, []int32, error) {
	return c.appendGridDiskDistances(nil, nil, k, gridDiskDistancesSafe, true)
}

// GridDiskDistancesUnsafe returns the grid disk and distances in ring-walk
// order. It fails with ErrPentagon if a pentagon or pentagon distortion is
// encountered.
//
// H3 C API: gridDiskDistancesUnsafe.
func (c Cell) GridDiskDistancesUnsafe(k int) ([]Cell, []int32, error) {
	return c.appendGridDiskDistances(nil, nil, k, gridDiskDistancesUnsafe, false)
}

func (c Cell) appendGridDiskDistances(dst []Cell, dstDist []int32, k int,
	impl func(h3Index, int32, []h3Index, []int32) h3Error, prune bool,
) ([]Cell, []int32, error) {
	sz, err := MaxGridDiskSize(k)
	if err != nil {
		return dst, dstDist, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz))
	distStart := len(dstDist)
	dstDist = slices.Grow(dstDist, int(sz))[:distStart+int(sz)]
	distWin := dstDist[distStart:]
	clear(distWin)
	if errC := impl(c, int32(k), win, distWin); errC != eSuccess {
		return dst[:start], dstDist[:distStart], toErr(errC)
	}
	if prune {
		n := 0
		for i, cell := range win {
			if cell != h3Null {
				win[n], distWin[n] = cell, distWin[i]
				n++
			}
		}
		dst, dstDist = dst[:start+n], dstDist[:distStart+n]
	}
	return dst, dstDist, nil
}

// GridDisksUnsafe returns the concatenated ring-walk-ordered grid disks of
// radius k around every origin. It fails with ErrPentagon if any disk
// encounters a pentagon; the result is grouped by origin, each group of size
// MaxGridDiskSize(k).
//
// H3 C API: gridDisksUnsafe.
func GridDisksUnsafe(origins []Cell, k int) ([]Cell, error) {
	sz, err := MaxGridDiskSize(k)
	if err != nil {
		return nil, err
	}
	out := make([]Cell, int64(len(origins))*sz)
	if errC := gridDisksUnsafe(origins, int32(k), out); errC != eSuccess {
		return nil, toErr(errC)
	}
	return out, nil
}

// GridRing returns the "hollow" ring of cells at exactly grid distance k
// from c, in no particular order.
//
// H3 C API: gridRing.
func (c Cell) GridRing(k int) ([]Cell, error) { return c.AppendGridRing(nil, k) }

// AppendGridRing appends the hollow ring of radius k around c to dst; see
// GridRing.
//
// H3 C API: gridRing.
func (c Cell) AppendGridRing(dst []Cell, k int) ([]Cell, error) {
	sz, err := MaxGridRingSize(k)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz))
	if errC := gridRing(c, int32(k), win); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst[:start+compactNonNull(win)], nil
}

// GridRingUnsafe returns the hollow ring of cells at exactly grid distance k
// from c in counterclockwise ring-walk order. It fails with ErrPentagon if a
// pentagon or pentagon distortion is encountered; use GridRing for the safe
// variant.
//
// H3 C API: gridRingUnsafe.
func (c Cell) GridRingUnsafe(k int) ([]Cell, error) { return c.AppendGridRingUnsafe(nil, k) }

// AppendGridRingUnsafe appends the counterclockwise-ordered hollow ring of
// radius k around c to dst; see GridRingUnsafe.
//
// H3 C API: gridRingUnsafe.
func (c Cell) AppendGridRingUnsafe(dst []Cell, k int) ([]Cell, error) {
	sz, err := MaxGridRingSize(k)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst, win := growZeroed(dst, int(sz))
	if errC := gridRingUnsafe(c, int32(k), win); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst, nil
}

// GridDistance returns the grid distance (minimum number of grid moves)
// between the two cells. It fails with ErrFailed when the distance cannot be
// computed (e.g. across pentagon distortion or very distant cells).
//
// H3 C API: gridDistance.
func (c Cell) GridDistance(other Cell) (int, error) {
	var d int64
	if errC := gridDistance(c, other, &d); errC != eSuccess {
		return 0, toErr(errC)
	}
	return int(d), nil
}

// GridPathLen returns the number of cells in the grid path from c to other,
// including both endpoints.
//
// H3 C API: gridPathCellsSize.
func (c Cell) GridPathLen(other Cell) (int, error) {
	var sz int64
	if errC := gridPathCellsSize(c, other, &sz); errC != eSuccess {
		return 0, toErr(errC)
	}
	return int(sz), nil
}

// GridPath returns the contiguous line of cells from c to other (inclusive).
// The path is not guaranteed unique or stable across library versions and
// may fail with ErrFailed across pentagon distortion.
//
// H3 C API: gridPathCells.
func (c Cell) GridPath(other Cell) ([]Cell, error) { return c.AppendGridPath(nil, other) }

// AppendGridPath appends the grid path from c to other (inclusive) to dst;
// see GridPath.
//
// H3 C API: gridPathCells.
func (c Cell) AppendGridPath(dst []Cell, other Cell) ([]Cell, error) {
	n, err := c.GridPathLen(other)
	if err != nil {
		return dst, err
	}
	start := len(dst)
	dst = slices.Grow(dst, n)[:start+n]
	if errC := gridPathCells(dst[start:], c, other); errC != eSuccess {
		return dst[:start], toErr(errC)
	}
	return dst, nil
}

// CellToLocalIJ returns the local IJ coordinates of cell anchored at origin.
// Coordinates are only comparable when produced with the same origin, and
// the conversion can fail with ErrFailed for cells too far apart or across
// pentagon distortion.
//
// H3 C API: cellToLocalIj.
func CellToLocalIJ(origin, cell Cell) (CoordIJ, error) {
	var ij CoordIJ
	if errC := cellToLocalIj(origin, cell, 0, &ij); errC != eSuccess {
		return CoordIJ{}, toErr(errC)
	}
	return ij, nil
}

// LocalIJToCell returns the cell at the local IJ coordinates anchored at
// origin; it is the inverse of CellToLocalIJ.
//
// H3 C API: localIjToCell.
func LocalIJToCell(origin Cell, ij CoordIJ) (Cell, error) {
	var out h3Index
	if errC := localIjToCell(origin, &ij, 0, &out); errC != eSuccess {
		return 0, toErr(errC)
	}
	return out, nil
}
