package h3

// geoMultiPolygon mirrors the GeoMultiPolygon struct from h3api.h: a
// simplified multiple-polygon type. It is internal in this port —
// the public CellsToMultiPolygon API returns []GeoPolygon
// (docs/DEVIATIONS.md §4); the ported 4.5.0 multipolygon pipeline and
// the area helpers consume this C-shaped form.
// Ported from H3 C: h3api.h.in::GeoMultiPolygon.
//
//nolint:unused // exercised by the cgo && c2go parity tests; consumed by the I-C multipolygon port (#34)
type geoMultiPolygon struct {
	NumPolygons int32
	Polygons    []GeoPolygon
}
