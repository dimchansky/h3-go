package c2go

import (
    "math"
)

// _geoAzDistanceRads computes the point p2 at azimuth az and distance from p1.
// Ported from H3 C: latLng.c::_geoAzDistanceRads
func _geoAzDistanceRads(p1 LatLng, az, distance float64) LatLng {
    const epsilon = 1e-16
    if distance < epsilon {
        return p1
    }
    az = _posAngleRads(az)

    var p2 LatLng
    if az < epsilon || math.Abs(az-math.Pi) < epsilon {
        // Due north or south
        if az < epsilon {
            p2.lat = p1.lat + distance
        } else {
            p2.lat = p1.lat - distance
        }
        if math.Abs(p2.lat-math.Pi/2) < epsilon {
            p2.lat = math.Pi / 2
            p2.lng = 0
        } else if math.Abs(p2.lat+math.Pi/2) < epsilon {
            p2.lat = -math.Pi / 2
            p2.lng = 0
        } else {
            p2.lng = constrainLng(p1.lng)
        }
        return p2
    }

    sinlat := math.Sin(p1.lat)*math.Cos(distance) + math.Cos(p1.lat)*math.Sin(distance)*math.Cos(az)
    if sinlat > 1.0 {
        sinlat = 1.0
    }
    if sinlat < -1.0 {
        sinlat = -1.0
    }
    p2.lat = math.Asin(sinlat)
    if math.Abs(p2.lat-math.Pi/2) < epsilon {
        p2.lat = math.Pi / 2
        p2.lng = 0
        return p2
    } else if math.Abs(p2.lat+math.Pi/2) < epsilon {
        p2.lat = -math.Pi / 2
        p2.lng = 0
        return p2
    }

    invcosp2lat := 1.0 / math.Cos(p2.lat)
    sinlng := math.Sin(az) * math.Sin(distance) * invcosp2lat
    coslng := (math.Cos(distance) - math.Sin(p1.lat)*math.Sin(p2.lat)) / math.Cos(p1.lat) * invcosp2lat
    if sinlng > 1.0 {
        sinlng = 1.0
    }
    if sinlng < -1.0 {
        sinlng = -1.0
    }
    if coslng > 1.0 {
        coslng = 1.0
    }
    if coslng < -1.0 {
        coslng = -1.0
    }
    p2.lng = constrainLng(p1.lng + math.Atan2(sinlng, coslng))
    return p2
}

