package h3

// nullIter creates a fully nulled-out child iterator for when an iterator is exhausted.
// This helps minimize the chance that a user will depend on the iterator internal state
// after it's exhausted, like the child resolution, for example.
// Ported from H3 C: iterators.c::_null_iter.
func nullIter() iterCellsChildren {
	return iterCellsChildren{
		H:         0, // h3Null
		ParentRes: -1,
		SkipDigit: -1,
	}
}
