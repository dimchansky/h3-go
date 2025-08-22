//go:build c2go

package c2go

import "testing"

func Test_setGeoDegs_ParityWithC(t *testing.T) {
	p := LatLng{Lat: 0.0, Lng: 0.0}
	goP := p
	setGeoDegs(&goP, 45.0, -120.0)

	cP := setGeoDegsC(p, 45.0, -120.0)

	if goP != cP {
		t.Fatalf("setGeoDegs mismatch: go=%+v c=%+v", goP, cP)
	}
}
