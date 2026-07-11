package h3

// compactCells takes a set of hexagons all at the same resolution and
// compresses them by pruning full child branches to the parent level. This is
// also done for all parents recursively to get the minimum number of hex
// addresses that perfectly cover the defined space.
// The compactedSet slice is modified in place. Remaining slots are filled with h3Null.
// Ported from H3 C: h3Index.c::compactCells.
func compactCells(h3Set []h3Index, compactedSet []h3Index, numHexes int64) h3Error {
	if numHexes == 0 {
		return eSuccess
	}

	if int64(len(h3Set)) < numHexes || int64(len(compactedSet)) < numHexes {
		return eFailed
	}

	res := getResolution(h3Set[0])
	if res == 0 {
		// No compaction possible, just copy the set to output
		for i := int64(0); i < numHexes; i++ {
			compactedSet[i] = h3Set[i]
		}
		return eSuccess
	}

	// Create working arrays
	remainingHexes := make([]h3Index, numHexes)
	copy(remainingHexes, h3Set)
	hashSetArray := make([]h3Index, numHexes)

	compactedSetOffset := int64(0)
	numRemainingHexes := numHexes

	for numRemainingHexes > 0 {
		res = getResolution(remainingHexes[0])
		parentRes := res - 1

		// If parentRes is less than zero, we've compacted all the way up to the
		// base cells. Time to process the remaining cells.
		if parentRes >= 0 {
			// Put the parents of the hexagons into the temp array
			// via a hashing mechanism, and use the reserved bits
			// to track how many times a parent is duplicated
			for i := int64(0); i < numRemainingHexes; i++ {
				currIndex := remainingHexes[i]
				if currIndex != 0 {
					// If the reserved bits were set by the caller, the
					// algorithm below may encounter undefined behavior
					// because it expects to have set the reserved bits
					// itself.
					if getReservedBits(currIndex) != 0 {
						return eCellInvalid
					}

					parent, parentError := cellToParent(currIndex, parentRes)
					// Should never be reachable as a result of the compact
					// algorithm. Can happen if cellToParent errors e.g.
					// because of incompatible resolutions.
					if parentError != eSuccess {
						return parentError
					}

					// Modulus hash the parent into the temp array
					loc := int64(parent % h3Index(numRemainingHexes))
					loopCount := int64(0)
					for hashSetArray[loc] != 0 {
						if loopCount > numRemainingHexes {
							// This case should not be possible because at
							// most one index is placed into hashSetArray
							// per numRemainingHexes.
							return eFailed
						}
						tempIndex := hashSetArray[loc] & h3Index(h3ReservedMaskNegative)
						if tempIndex == parent {
							count := getReservedBits(hashSetArray[loc]) + 1
							limitCount := int32(7)
							if isPentagon(tempIndex & h3Index(h3ReservedMaskNegative)) {
								limitCount--
							}
							// One is added to count for this check to match
							// one being added to count later in this
							// function when checking for all children being
							// present.
							if count+1 > limitCount {
								// Only possible on duplicate input
								return eDuplicateInput
							}
							parent = setReservedBits(parent, count)
							hashSetArray[loc] = h3Null
						} else {
							loc = (loc + 1) % numRemainingHexes
						}
						loopCount++
					}
					hashSetArray[loc] = parent
				}
			}
		}

		// Determine which parent hexagons have a complete set
		// of children and put them in the compactableHexes array
		compactableCount := int64(0)
		maxCompactableCount := numRemainingHexes / 6 // Somehow all pentagons; conservative
		if maxCompactableCount == 0 {
			for i := int64(0); i < numRemainingHexes; i++ {
				compactedSet[compactedSetOffset+i] = remainingHexes[i]
			}
			compactedSetOffset += numRemainingHexes
			break
		}

		compactableHexes := make([]h3Index, maxCompactableCount)
		for i := int64(0); i < numRemainingHexes; i++ {
			if hashSetArray[i] == 0 {
				continue
			}
			count := getReservedBits(hashSetArray[i]) + 1
			// Include the deleted direction for pentagons as implicitly "there"
			if isPentagon(hashSetArray[i] & h3Index(h3ReservedMaskNegative)) {
				// We need this later on, no need to recalculate
				hashSetArray[i] = setReservedBits(hashSetArray[i], count)
				// Increment count after setting the reserved bits,
				// since count is already incremented above, so it
				// will be the expected value for a complete hexagon.
				count++
			}
			if count == 7 {
				// Bingo! Full set!
				compactableHexes[compactableCount] = hashSetArray[i] & h3Index(h3ReservedMaskNegative)
				compactableCount++
			}
		}

		// Uncompactable hexes are immediately copied into the
		// output compactedSetOffset
		uncompactableCount := int64(0)
		for i := int64(0); i < numRemainingHexes; i++ {
			currIndex := remainingHexes[i]
			if currIndex != h3Null {
				isUncompactable := true
				// Resolution 0 cells always uncompactable, and trying to take
				// the res -1 parent of a cell is invalid.
				if parentRes >= 0 {
					parent, parentError := cellToParent(currIndex, parentRes)
					if parentError != eSuccess {
						return parentError
					}
					// Modulus hash the parent into the temp array
					// to determine if this index was included in
					// the compactableHexes array
					loc := int64(parent % h3Index(numRemainingHexes))
					loopCount := int64(0)
					// Equivalent to C do-while loop
					for {
						if loopCount > numRemainingHexes {
							// This case should not be possible because at most
							// one index is placed into hashSetArray per input
							// hexagon.
							return eFailed
						}
						tempIndex := hashSetArray[loc] & h3Index(h3ReservedMaskNegative)
						if tempIndex == parent {
							count := getReservedBits(hashSetArray[loc]) + 1
							if count == 7 {
								isUncompactable = false
							}
							break
						} else {
							loc = (loc + 1) % numRemainingHexes
						}
						loopCount++
						// Exit condition: equivalent to C do-while condition
						if hashSetArray[loc] == parent {
							break
						}
					}
				}
				if isUncompactable {
					compactedSet[compactedSetOffset+uncompactableCount] = remainingHexes[i]
					uncompactableCount++
				}
			}
		}

		// Set up for the next loop
		// Clear hash set array (matches C: memset(hashSetArray, 0, numHexes * sizeof(h3Index)))
		for i := int64(0); i < numHexes; i++ {
			hashSetArray[i] = 0
		}
		compactedSetOffset += uncompactableCount
		copy(remainingHexes, compactableHexes[:compactableCount])
		numRemainingHexes = compactableCount
	}

	// Fill remaining slots with h3Null
	for i := compactedSetOffset; i < numHexes; i++ {
		compactedSet[i] = h3Null
	}

	return eSuccess
}
