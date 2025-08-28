package h3

// maxGridDiskSize calculates the maximum number of cells that can be in the result of gridDisk for the given k.
// Formula source and proof: https://oeis.org/A003215
// Ported from H3 C: algos.c::maxGridDiskSize
func maxGridDiskSize(k int32, out *int64) H3Error {
	if k < 0 {
		return E_DOMAIN
	}
	if k >= K_ALL_CELLS_AT_RES_15 {
		// If a k value of this value or above is provided, this function will
		// estimate more cells than exist in the H3 grid at the finest
		// resolution. This is a problem since the function does signed integer
		// arithmetic on `k`, which could overflow. To prevent that, instead
		// substitute the maximum number of cells in the grid, as it should not
		// be possible for the gridDisk functions to exceed that. Note this is
		// not resolution specific. So, when resolution < 15, this function may
		// still estimate a size larger than the number of cells in the grid.
		numCells, err := getNumCells(MAX_H3_RES)
		if err != E_SUCCESS {
			return err
		}
		*out = numCells
		return E_SUCCESS
	}
	*out = 3*int64(k)*(int64(k)+1) + 1
	return E_SUCCESS
}
