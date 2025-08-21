package c2go

import "math"

// constrainLat makes sure latitudes are in the proper bounds.
// Ported from H3 C: latLng.c::constrainLat
func constrainLat(lat float64) float64 {
    for lat > math.Pi/2 {
        lat = lat - math.Pi
    }
    return lat
}

