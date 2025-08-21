package faceijk

import (
	"math"
	"testing"
)

func TestFaceConstants(t *testing.T) {
	if len(FaceCenterGeo) != NumIcosaFaces {
		t.Errorf("FaceCenterGeo length = %d, want %d", len(FaceCenterGeo), NumIcosaFaces)
	}
	if len(FaceAxesAzRadsCII) != NumIcosaFaces {
		t.Errorf("FaceAxesAzRadsCII length = %d, want %d", len(FaceAxesAzRadsCII), NumIcosaFaces)
	}
	if len(FaceAxesAzRadsCIII) != NumIcosaFaces {
		t.Errorf("FaceAxesAzRadsCIII length = %d, want %d", len(FaceAxesAzRadsCIII), NumIcosaFaces)
	}
	for i, center := range FaceCenterGeo {
		lat, lng := center[0], center[1]
		if lat < -math.Pi/2 || lat > math.Pi/2 {
			t.Errorf("Face %d latitude %f out of range", i, lat)
		}
		if lng < -math.Pi || lng > math.Pi {
			t.Errorf("Face %d longitude %f out of range", i, lng)
		}
	}
	expectedMaxRes := 16
	if len(MaxDimByCIIres) != expectedMaxRes {
		t.Errorf("MaxDimByCIIres length = %d, want %d", len(MaxDimByCIIres), expectedMaxRes)
	}
	if len(MaxDimByCIIIres) != expectedMaxRes {
		t.Errorf("MaxDimByCIIIres length = %d, want %d", len(MaxDimByCIIIres), expectedMaxRes)
	}
}
