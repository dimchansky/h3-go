package h3

// cellToChildren takes the given hexagon id and generates all of the children
// at the specified resolution storing them into the provided memory slice.
// It's assumed that cellToChildrenSize was used to determine the allocation.
// Ported from H3 C: h3Index.c::cellToChildren.
func cellToChildren(h H3Index, childRes int32, children []H3Index) H3Error {
	i := int64(0)
	var iter IterCellsChildren
	iterInitParent(h, childRes, &iter)
	for iter.H != H3_NULL {
		children[i] = iter.H
		i++
		iterStepChild(&iter)
	}
	return E_SUCCESS
}
