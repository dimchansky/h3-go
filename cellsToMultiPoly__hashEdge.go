package h3

// hashEdge hashes an H3Index to a bucket index for hash table lookups.
//
// Uses a mixing function based on SplitMix64 to ensure good
// distribution of hash values. x is the H3Index value to hash, n the
// number of hash table buckets; the result is a bucket index in range
// [0, n-1].
//
// Reference: Steele et al., "Fast splittable pseudorandom number
// generators" OOPSLA 2014. https://doi.org/10.1145/2660193.2660195
// Ported from H3 C: cellsToMultiPoly.c::hashEdge.
func hashEdge(x h3Index, n uint64) uint64 {
	v := uint64(x)
	v ^= v >> 30
	v *= 0xbf58476d1ce4e5b9
	v ^= v >> 27
	v *= 0x94d049bb133111eb
	v ^= v >> 31

	return v % n
}
