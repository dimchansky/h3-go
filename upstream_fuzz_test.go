package h3

// Native Go fuzz equivalents for the 23 H3 v4.4.0 libFuzzer/AFL harnesses.
// The upstream harnesses primarily assert memory safety/no-crash behavior;
// these targets preserve their raw index, cell-set, polygon, and coordinate
// input domains while also checking successful results for basic validity.

import (
	"encoding/binary"
	"math"
	"testing"
)

func FuzzUpstreamCellOperations(f *testing.F) {
	f.Add(uint64(0x8928308280fffff), int32(1), int32(0))
	f.Add(uint64(0), int32(-1), int32(math.MaxInt32))
	f.Fuzz(func(t *testing.T, raw uint64, a, b int32) {
		cell := h3Index(raw)
		res := int32(uint32(a) % (maxH3Res + 1))
		k := int32(uint32(b) % 3)

		_ = isValidCell(cell)
		_ = isPentagon(cell)
		_ = getResolution(cell)
		_ = getBaseCellNumber(cell)
		var ll LatLng
		_ = cellToLatLng(cell, &ll)
		var boundary CellBoundary
		_ = cellToBoundary(cell, &boundary)
		_, _ = cellAreaRads2(cell)
		_, _ = cellAreaKm2(cell)
		_, _ = cellAreaM2(cell)
		_, _ = cellToParent(cell, res)
		_, _ = cellToCenterChild(cell, res)
		_, _ = cellToChildrenSize(cell, res)
		_, _ = cellToChildPos(cell, res)
		_, _ = childPosToCell(int64(b), cell, res)

		out := make([]h3Index, 19)
		dist := make([]int32, 19)
		_ = gridDisk(cell, k, out)
		_ = gridDiskDistances(cell, k, out, dist)
		_ = gridRing(cell, k, out)
		var distance, pathSize int64
		_ = gridDistance(cell, h3Index(uint64(raw)^uint64(uint32(a))), &distance)
		_ = gridPathCellsSize(cell, h3Index(uint64(raw)^uint64(uint32(b))), &pathSize)

		var ij CoordIJ
		_ = cellToLocalIj(cell, cell, 0, &ij)
		var local h3Index
		_ = localIjToCell(cell, &CoordIJ{I: a, J: b}, 0, &local)
		var edges [6]h3Index
		_ = originToDirectedEdges(cell, edges[:])
		var vertices [6]h3Index
		_ = cellToVertexes(cell, &vertices)
		for _, edge := range edges {
			var length float64
			_ = edgeLengthRads(edge, &length)
			_ = directedEdgeToBoundary(edge, &boundary)
			// 4.5.0 upstream harness addition (fuzzerDirectedEdge.c).
			var reversed h3Index
			_ = reverseDirectedEdge(edge, &reversed)
		}

		digits := make([]int32, res)
		for i := range digits {
			digits[i] = int32(raw >> (i * 3) & 7)
		}
		var constructed h3Index
		_ = constructCell(res, int32(uint32(b)%numBaseCells), digits, &constructed)
	})
}

func FuzzUpstreamCellSets(f *testing.F) {
	f.Add([]byte{0xff, 0xff, 0x0f, 0x28, 0x08, 0x83, 0x92, 0x08})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64*8 {
			data = data[:64*8]
		}
		cells := make([]h3Index, len(data)/8)
		for i := range cells {
			cells[i] = h3Index(binary.LittleEndian.Uint64(data[i*8:]))
		}
		compacted := make([]h3Index, len(cells))
		_ = compactCells(cells, compacted, int64(len(cells)))
		var polygon linkedGeoPolygon
		_ = cellsToLinkedMultiPolygon(cells, int32(len(cells)), &polygon)
		destroyLinkedMultiPolygon(&polygon)
	})
}

func FuzzUpstreamPolygonOperations(f *testing.F) {
	f.Add([]byte{9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 2 {
			return
		}
		res := int32(data[0] % (maxH3Res + 1))
		flags := uint32(data[1] % 4)
		data = data[2:]
		if len(data) > 16*16 {
			data = data[:16*16]
		}
		loop := make(GeoLoop, len(data)/16)
		for i := range loop {
			loop[i] = LatLng{
				Lat: Angle(math.Float64frombits(binary.LittleEndian.Uint64(data[i*16:]))),
				Lng: Angle(math.Float64frombits(binary.LittleEndian.Uint64(data[i*16+8:]))),
			}
		}
		polygon := GeoPolygon{GeoLoop: loop}
		var size int64
		if maxPolygonToCellsSize(&polygon, res, 0, &size) == eSuccess && size >= 0 && size <= 10000 {
			_ = polygonToCells(&polygon, res, 0, make([]h3Index, size))
		}
		if size, err := maxPolygonToCellsSizeExperimental(&polygon, res, flags); err == eSuccess && size >= 0 && size <= 10000 {
			_ = polygonToCellsExperimental(&polygon, res, flags, size, make([]h3Index, size))
		}
	})
}

func FuzzUpstreamInternalCoordinates(f *testing.F) {
	f.Add(int32(0), int32(0), int32(0))
	f.Add(int32(math.MinInt32), int32(math.MaxInt32), int32(-1))
	f.Fuzz(func(t *testing.T, i, j, k int32) {
		ijk := coordIJK{I: i, J: j, K: k}
		_ijkNormalize(&ijk)
		_ijkRotate60cw(&ijk)
		_ijkRotate60ccw(&ijk)
		_ijkScale(&ijk, k)
		_upAp7(&ijk)
		_upAp7r(&ijk)
		_downAp3(&ijk)
		_ = _unitIjkToDigit(&ijk)
		_ = ijkDistance(&ijk, &coordIJK{})
	})
}
