package h3

import (
	"math"
	"sort"
)

// createGlobeMultiPolygon allocates a GeoMultiPolygon representing the
// entire globe. The globe is represented using 8 triangular polygons,
// with all edge arcs of exactly 90 degrees (i.e., pi/2 radians).
// Memory should be freed with `destroyGeoMultiPolygon` (in Go: the
// garbage collector).
// Ported from H3 C: cellsToMultiPoly.c::createGlobeMultiPolygon.
func createGlobeMultiPolygon(mpoly *geoMultiPolygon) h3Error {
	const numPolygons = 8
	const numVerts = 3
	verts := [8][3]LatLng{
		{{Lat: Rad(math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}},
		{{Lat: Rad(math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}, {Lat: Rad(0), Lng: Rad(math.Pi)}},
		{{Lat: Rad(math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi)}, {Lat: Rad(0), Lng: Rad(-math.Pi / 2)}},
		{{Lat: Rad(math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(-math.Pi / 2)}, {Lat: Rad(0), Lng: Rad(0)}},
		{{Lat: Rad(-math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(-math.Pi / 2)}},
		{{Lat: Rad(-math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(-math.Pi / 2)}, {Lat: Rad(0), Lng: Rad(-math.Pi)}},
		{{Lat: Rad(-math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(-math.Pi)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}},
		{{Lat: Rad(-math.Pi / 2), Lng: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}, {Lat: Rad(0), Lng: Rad(0)}},
	}

	spolys := make([]sortablePoly, numPolygons)

	for i := 0; i < numPolygons; i++ {
		poly := &spolys[i].poly
		poly.Holes = nil
		poly.GeoLoop = make(GeoLoop, numVerts)

		for j := 0; j < numVerts; j++ {
			poly.GeoLoop[j] = verts[i][j]
		}

		// Calculate outer area for sorting
		spolys[i].outerArea, _ = geoLoopAreaRads2(poly.GeoLoop)
	}

	sort.Slice(spolys, func(a, b int) bool {
		return cmp_SortablePoly(&spolys[a], &spolys[b]) < 0
	})

	mpoly.Polygons = make([]GeoPolygon, numPolygons)

	mpoly.NumPolygons = numPolygons
	for i := 0; i < numPolygons; i++ {
		mpoly.Polygons[i] = spolys[i].poly
	}

	return eSuccess
}
