package c2go

// LinkedLatLng mirrors the C struct used in h3api.h - a coordinate node in a linked geo structure
type LinkedLatLng struct {
	Vertex LatLng
	Next   *LinkedLatLng
}

// LinkedGeoLoop mirrors the C struct used in h3api.h - a loop node in a linked geo structure
type LinkedGeoLoop struct {
	First *LinkedLatLng
	Last  *LinkedLatLng
	Next  *LinkedGeoLoop
}

// LinkedGeoPolygon mirrors the C struct used in h3api.h - a polygon node in a linked geo structure
type LinkedGeoPolygon struct {
	First *LinkedGeoLoop
	Last  *LinkedGeoLoop
	Next  *LinkedGeoPolygon
}
