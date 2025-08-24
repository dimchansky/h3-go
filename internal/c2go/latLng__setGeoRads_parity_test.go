//go:build cgo

package c2go

import "testing"

func Test_setGeoRads_ParityWithC(t *testing.T) {
	p := LatLng{Lat: 0.0, Lng: 0.0}
	goP := p
	cP := p

	_setGeoRads(&goP, 1.23, -0.77)
	_setGeoRadsC(&cP, 1.23, -0.77)

	if goP != cP {
		t.Fatalf("_setGeoRads mismatch: go=%+v c=%+v", goP, cP)
	}
}
