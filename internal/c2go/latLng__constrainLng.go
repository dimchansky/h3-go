package c2go

import "math"

// constrainLng makes sure longitudes are in the proper bounds.
// Ported from H3 C: latLng.c::constrainLng
func constrainLng(lng float64) float64 {
    for lng > math.Pi {
        lng = lng - (2 * math.Pi)
    }
    for lng < -math.Pi {
        lng = lng + (2 * math.Pi)
    }
    return lng
}

