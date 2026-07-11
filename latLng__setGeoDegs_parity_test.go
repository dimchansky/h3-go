//go:build cgo && c2go

package h3

import "testing"

func Test_setGeoDegs_ParityWithC(t *testing.T) {
	p := LatLng{Lat: 0.0, Lng: 0.0}
	goP := p
	cP := p

	setGeoDegs(&goP, 45.0, -120.0)
	setGeoDegsC(&cP, 45.0, -120.0)

	if goP != cP {
		t.Fatalf("setGeoDegs mismatch: go=%+v c=%+v", goP, cP)
	}
}
