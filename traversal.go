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
// radius k, for pre-sizing AppendGridDisk destination buffers. Like every
// k-taking function in this package, it accepts 0 <= k <= math.MaxInt32 and
// fails with ErrDomain otherwise (the pre-check guards the int32 narrowing;
// docs/DEVIATIONS.md).
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
// of radius k (1 for k=0), for pre-sizing AppendGridRing destination
// buffers. k outside 0..math.MaxInt32 fails with ErrDomain.
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
// no particular order. k outside 0..math.MaxInt32 fails with ErrDomain.
// Near pentagons the result may contain fewer than MaxGridDiskSize(k)
// cells: the H3_NULL holes C leaves in its output are pruned
// (docs/DEVIATIONS.md).
//
// H3 C API: gridDisk.
func (c Cell) GridDisk(k int) ([]Cell, error) { return c.AppendGridDisk(nil, k) }

// AppendGridDisk appends all cells within grid distance k of c to dst and
// returns the extended slice, in no particular order. Pass dst[:0] (or nil)
// to reuse dst's capacity; when capacity suffices the call may still
// allocate an internal distance scratch, but only when the fast unsafe
// algorithm falls back to the safe one near pentagons (as in H3 C). On
// error the returned slice has dst's original length and elements — no
// partial results are observable.
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

// GridDiskUnsafe returns all cells within grid distance k of c, in order of
// increasing distance from the origin (origin first, then each ring in
// turn; order within a ring is not part of the contract). It fails with
// ErrPentagon if a pentagon or pentagon distortion is encountered, because
// the fast algorithm's output would be undefined there; use GridDisk for
// the safe variant. The result is dense (no pruning).
//
// H3 C API: gridDiskUnsafe.
func (c Cell) GridDiskUnsafe(k int) ([]Cell, error) { return c.AppendGridDiskUnsafe(nil, k) }

// AppendGridDiskUnsafe appends the grid disk of radius k to dst in order of
// increasing distance; see GridDiskUnsafe. On error the returned slice has
// dst's original length and elements.
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
// each cell's grid distance from c (in grid moves, not a geographic
// measure), in no particular order. k outside 0..math.MaxInt32 fails with
// ErrDomain, and near pentagons the result may contain fewer than
// MaxGridDiskSize(k) cells (H3_NULL holes pruned, cells and distances in
// tandem). The distances slice is int32 to match the H3 C representation
// without a conversion copy.
//
// H3 C API: gridDiskDistances.
func (c Cell) GridDiskDistances(k int) ([]Cell, []int32, error) {
	return c.AppendGridDiskDistances(nil, nil, k)
}

// GridDiskDistancesGrouped returns cells within grid distance k of c grouped
// by distance: the result always has k+1 rings and result[d] contains
// exactly the cells at distance d. Empty rings are retained, H3_NULL holes
// near pentagons are omitted, and no order within a ring is guaranteed. All
// rings share one backing cell array and have their capacity limited to
// their length, so appending to one ring cannot overwrite another. k
// outside 0..math.MaxInt32 fails with ErrDomain (returning a nil result).
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
// both capacities suffice the call does not allocate. On error both
// returned slices have their original lengths and elements — the cell and
// distance outputs are rolled back in tandem, and no partial results are
// observable.
//
// H3 C API: gridDiskDistances.
func (c Cell) AppendGridDiskDistances(dst []Cell, dstDist []int32, k int) ([]Cell, []int32, error) {
	return c.appendGridDiskDistances(dst, dstDist, k, gridDiskDistances, true)
}

// GridDiskDistancesSafe is the variant of GridDiskDistances that never runs
// the optimistic unsafe fast path: GridDiskDistances tries the faster
// unsafe algorithm first and falls back to the safe one when a pentagon is
// encountered, while this function runs the safe algorithm directly. The
// results are equivalent.
//
// H3 C API: gridDiskDistancesSafe.
func (c Cell) GridDiskDistancesSafe(k int) ([]Cell, []int32, error) {
	return c.appendGridDiskDistances(nil, nil, k, gridDiskDistancesSafe, true)
}

// GridDiskDistancesUnsafe returns the grid disk and distances in order of
// increasing distance from the origin (order within a ring is not part of
// the contract). It fails with ErrPentagon if a pentagon or pentagon
// distortion is encountered; the result is dense (no pruning).
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

// GridDisksUnsafe returns the concatenated grid disks of radius k around
// every origin. The result is grouped by origin in input order, each group
// of size MaxGridDiskSize(k) ordered by increasing ring distance; upstream
// guarantees no sorting within each ring group. It fails with ErrPentagon
// (returning a nil slice) if any disk encounters a pentagon or pentagon
// distortion.
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
// from c, in no particular order. k=0 returns just the origin cell. k
// outside 0..math.MaxInt32 fails with ErrDomain (in C, negative k is
// undefined behavior here — docs/DEVIATIONS.md). Near pentagons the ring
// may contain fewer than MaxGridRingSize(k) cells (H3_NULL holes pruned).
//
// H3 C API: gridRing.
func (c Cell) GridRing(k int) ([]Cell, error) { return c.AppendGridRing(nil, k) }

// AppendGridRing appends the hollow ring of radius k around c to dst; see
// GridRing. On error the returned slice has dst's original length and
// elements.
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

// GridRingUnsafe returns the hollow ring of cells at exactly grid distance
// k from c (k=0 returns just the origin cell); the order within the ring is
// not part of the contract. It fails with ErrPentagon if a pentagon or
// pentagon distortion is encountered — upstream notes these failure cases
// may be fixed in future versions. Use GridRing for the safe variant.
//
// H3 C API: gridRingUnsafe.
func (c Cell) GridRingUnsafe(k int) ([]Cell, error) { return c.AppendGridRingUnsafe(nil, k) }

// AppendGridRingUnsafe appends the hollow ring of radius k around c to dst;
// see GridRingUnsafe. On error the returned slice has dst's original length
// and elements.
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

// GridDistance returns the grid distance — the minimum number of grid
// moves between the two cells, not a Euclidean or great-circle measure.
// Cells of different resolutions fail with ErrResolutionMismatch. The
// computation itself can fail with ErrFailed, e.g. across pentagon
// distortion or for very distant cells.
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
// including both endpoints — always GridDistance(c, other)+1. It shares
// GridDistance's failure modes: ErrResolutionMismatch for cells of
// different resolutions, ErrFailed when no path can be computed.
//
// H3 C API: gridPathCellsSize.
func (c Cell) GridPathLen(other Cell) (int, error) {
	var sz int64
	if errC := gridPathCellsSize(c, other, &sz); errC != eSuccess {
		return 0, toErr(errC)
	}
	return int(sz), nil
}

// GridPath returns the line of cells from c to other, ordered from c to
// other inclusive. Upstream guarantees the path's length is
// GridDistance(c, other)+1 and that every cell is a neighbor of the
// preceding one. Paths are drawn in grid space and may not correspond to
// Cartesian lines or great arcs; the specific path is not guaranteed unique
// or stable across H3 versions. Cells of different resolutions fail with
// ErrResolutionMismatch; ErrFailed or ErrPentagon when no path can be
// computed (interpolation is attempted from both endpoints' local
// coordinate charts, so many pentagon-adjacent paths succeed, but some
// pairs remain uncomputable).
//
// H3 C API: gridPathCells.
func (c Cell) GridPath(other Cell) ([]Cell, error) { return c.AppendGridPath(nil, other) }

// AppendGridPath appends the grid path from c to other (inclusive) to dst;
// see GridPath. On error the returned slice has dst's original length and
// elements.
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

// CellToLocalIJ returns the local IJ coordinates of cell anchored at
// origin. Coordinates are only comparable when produced with the same
// origin, and the local coordinate space is not guaranteed to be compatible
// across H3 versions — do not persist IJ coordinates across library
// upgrades. origin and cell of different resolutions fail with
// ErrResolutionMismatch; the conversion itself can fail with ErrFailed for
// cells too far apart or across pentagon distortion.
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
// origin. It fails with ErrFailed when the coordinates are too far from the
// origin, cross pentagon distortion, or do not correspond to a cell. Like
// CellToLocalIJ, the local coordinate space is not guaranteed to be
// compatible across H3 versions, so it inverts CellToLocalIJ only within a
// single H3 version and for the same origin.
//
// H3 C API: localIjToCell.
func LocalIJToCell(origin Cell, ij CoordIJ) (Cell, error) {
	var out h3Index
	if errC := localIjToCell(origin, &ij, 0, &out); errC != eSuccess {
		return 0, toErr(errC)
	}
	return out, nil
}
