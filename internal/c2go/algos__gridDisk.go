package c2go

// gridDisk produces cells within grid distance k of the origin cell.
// This is a convenience wrapper around gridDiskDistances that ignores distance values.
// Ported from H3 C: algos.c::gridDisk
func gridDisk(origin H3Index, k int32, out []H3Index) H3Error {
	return gridDiskDistances(origin, k, out, nil)
}
