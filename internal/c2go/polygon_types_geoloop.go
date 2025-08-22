package c2go

// GeoLoop is a simple loop of LatLng coordinates (closed implicitly).
type GeoLoop []LatLng

// GeoPolygon is an outer loop with optional holes.
type GeoPolygon struct {
    Geoloop GeoLoop
    Holes   []GeoLoop
}

