//go:build cgo && c2go && !h3v450

package h3

import (
	"testing"
)

func Test_hashVertex_parity(t *testing.T) {
	testCases := []struct {
		name       string
		vertex     LatLng
		res        int32
		numBuckets int32
	}{
		{
			name:       "basic hash",
			vertex:     LatLng{Lat: 37.775, Lng: -122.418},
			res:        9,
			numBuckets: 100,
		},
		{
			name:       "zero coordinates",
			vertex:     LatLng{Lat: 0.0, Lng: 0.0},
			res:        5,
			numBuckets: 50,
		},
		{
			name:       "negative coordinates",
			vertex:     LatLng{Lat: -45.0, Lng: -90.0},
			res:        7,
			numBuckets: 200,
		},
		{
			name:       "high resolution",
			vertex:     LatLng{Lat: 40.7128, Lng: -74.0060},
			res:        15,
			numBuckets: 1000,
		},
		{
			name:       "low resolution",
			vertex:     LatLng{Lat: 51.5074, Lng: -0.1278},
			res:        0,
			numBuckets: 10,
		},
		{
			name:       "edge case - max resolution",
			vertex:     LatLng{Lat: 89.9, Lng: 179.9},
			res:        15,
			numBuckets: 500,
		},
		{
			name:       "small bucket count",
			vertex:     LatLng{Lat: 35.6762, Lng: 139.6503},
			res:        10,
			numBuckets: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			goResult := _hashVertex(&tc.vertex, tc.res, tc.numBuckets)
			cResult := _hashVertexC(&tc.vertex, tc.res, tc.numBuckets)

			if goResult != cResult {
				t.Errorf("_hashVertex parity mismatch for %s:\n  Go result: %d\n  C result: %d",
					tc.name, goResult, cResult)
			}
		})
	}
}
