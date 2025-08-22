//go:build c2go

package c2go

import "testing"

func Test_cellBoundaryCrossesGeoLoop_ParityWithC(t *testing.T) {
    loop := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
    loopBBox := bboxFromGeoLoop(loop)
    // Boundary crossing the square
    boundary := CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: -1, Lng: 1}, {Lat: 3, Lng: 1}}}
    boundaryBBox := bboxFromGeoLoop(boundary.Verts)
    goVal := cellBoundaryCrossesGeoLoop(loop, loopBBox, boundary, boundaryBBox)
    cVal := cellBoundaryCrossesGeoLoopC(loop, loopBBox, boundary, boundaryBBox)
    if goVal != cVal {
        t.Fatalf("cellBoundaryCrossesGeoLoop mismatch (crossing)")
    }
    // Boundary outside
    boundary = CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: -1, Lng: -1}, {Lat: -2, Lng: -2}}}
    boundaryBBox = bboxFromGeoLoop(boundary.Verts)
    goVal = cellBoundaryCrossesGeoLoop(loop, loopBBox, boundary, boundaryBBox)
    cVal = cellBoundaryCrossesGeoLoopC(loop, loopBBox, boundary, boundaryBBox)
    if goVal != cVal {
        t.Fatalf("cellBoundaryCrossesGeoLoop mismatch (outside)")
    }
}

