//go:build cgo && c2go

package c2go

import (
	"testing"
)

func Test_addLinkedCoord_parity(t *testing.T) {
	t.Run("empty loop", func(t *testing.T) {
		// Test Go implementation
		goLoop := &LinkedGeoLoop{First: nil, Last: nil, Next: nil}
		vertex := &LatLng{Lat: 0.123, Lng: 0.456}
		goResult := addLinkedCoord(goLoop, vertex)

		// Check Go behavior
		goReturnsCoord := (goResult != nil)
		goSetsFirst := (goLoop.First == goResult)
		goSetsLast := (goLoop.Last == goResult)
		goVertexMatches := (goResult != nil && goResult.Vertex.Lat == vertex.Lat && goResult.Vertex.Lng == vertex.Lng)

		// Test C implementation behavior
		cReturnsCoord, cSetsFirst, cSetsLast, _, cVertexMatches := addLinkedCoordC(true, false, *vertex)

		// Compare behaviors
		if goReturnsCoord != cReturnsCoord {
			t.Errorf("Return behavior mismatch: Go=%v, C=%v", goReturnsCoord, cReturnsCoord)
		}
		if goSetsFirst != cSetsFirst {
			t.Errorf("First pointer behavior mismatch: Go=%v, C=%v", goSetsFirst, cSetsFirst)
		}
		if goSetsLast != cSetsLast {
			t.Errorf("Last pointer behavior mismatch: Go=%v, C=%v", goSetsLast, cSetsLast)
		}
		if goVertexMatches != cVertexMatches {
			t.Errorf("Vertex copy behavior mismatch: Go=%v, C=%v", goVertexMatches, cVertexMatches)
		}
	})

	t.Run("loop with existing coordinate", func(t *testing.T) {
		// Test Go implementation
		existingCoord := &LinkedLatLng{
			Vertex: LatLng{Lat: 0.5, Lng: 1.0},
			Next:   nil,
		}
		goLoop := &LinkedGeoLoop{First: existingCoord, Last: existingCoord, Next: nil}
		vertex := &LatLng{Lat: 0.789, Lng: 0.012}
		goResult := addLinkedCoord(goLoop, vertex)

		// Check Go behavior
		goReturnsCoord := (goResult != nil)
		goKeepsFirst := (goLoop.First == existingCoord)
		goSetsLast := (goLoop.Last == goResult)
		goLinksCoords := (existingCoord.Next == goResult)
		goVertexMatches := (goResult != nil && goResult.Vertex.Lat == vertex.Lat && goResult.Vertex.Lng == vertex.Lng)

		// Test C implementation behavior
		cReturnsCoord, cKeepsFirst, cSetsLast, cLinksCoords, cVertexMatches := addLinkedCoordC(false, true, *vertex)

		// Compare behaviors
		if goReturnsCoord != cReturnsCoord {
			t.Errorf("Return behavior mismatch: Go=%v, C=%v", goReturnsCoord, cReturnsCoord)
		}
		if goKeepsFirst != cKeepsFirst {
			t.Errorf("First pointer behavior mismatch: Go=%v, C=%v", goKeepsFirst, cKeepsFirst)
		}
		if goSetsLast != cSetsLast {
			t.Errorf("Last pointer behavior mismatch: Go=%v, C=%v", goSetsLast, cSetsLast)
		}
		if goLinksCoords != cLinksCoords {
			t.Errorf("Coordinate linking behavior mismatch: Go=%v, C=%v", goLinksCoords, cLinksCoords)
		}
		if goVertexMatches != cVertexMatches {
			t.Errorf("Vertex copy behavior mismatch: Go=%v, C=%v", goVertexMatches, cVertexMatches)
		}
	})

	t.Run("multiple coordinates", func(t *testing.T) {
		// Test Go implementation with multiple additions
		goLoop := &LinkedGeoLoop{First: nil, Last: nil, Next: nil}

		// Add first coordinate
		vertex1 := &LatLng{Lat: 0.1, Lng: 0.2}
		coord1 := addLinkedCoord(goLoop, vertex1)

		// Add second coordinate
		vertex2 := &LatLng{Lat: 0.3, Lng: 0.4}
		coord2 := addLinkedCoord(goLoop, vertex2)

		// Add third coordinate
		vertex3 := &LatLng{Lat: 0.5, Lng: 0.6}
		coord3 := addLinkedCoord(goLoop, vertex3)

		// Check Go behavior for final state
		goHasThreeCoords := (coord1 != nil && coord2 != nil && coord3 != nil)
		goFirstIsCoord1 := (goLoop.First == coord1)
		goLastIsCoord3 := (goLoop.Last == coord3)
		goProperLinking := (coord1.Next == coord2 && coord2.Next == coord3 && coord3.Next == nil)

		// We can't directly test multiple additions with C, but we can verify the pattern
		// by checking that each addition follows the expected behavior
		_, c1SetsFirst, c1SetsLast, _, _ := addLinkedCoordC(true, false, *vertex1)
		_, c2KeepsFirst, c2SetsLast, c2Links, _ := addLinkedCoordC(false, true, *vertex2)
		_, c3KeepsFirst, c3SetsLast, c3Links, _ := addLinkedCoordC(false, true, *vertex3)

		// Verify the pattern matches
		if !goHasThreeCoords {
			t.Error("Go failed to create three coordinates")
		}
		if !goFirstIsCoord1 || !c1SetsFirst {
			t.Errorf("First coordinate behavior mismatch: Go=%v, C=%v", goFirstIsCoord1, c1SetsFirst)
		}
		if !goLastIsCoord3 || !c3SetsLast {
			t.Errorf("Last coordinate behavior mismatch: Go=%v, C=%v", goLastIsCoord3, c3SetsLast)
		}
		if !goProperLinking || !c2Links || !c3Links {
			t.Errorf("Linking behavior mismatch: Go=%v, C2=%v, C3=%v", goProperLinking, c2Links, c3Links)
		}
		if !c2KeepsFirst || !c3KeepsFirst {
			t.Error("C should keep first pointer when adding to non-empty loop")
		}
		if !c1SetsLast || !c2SetsLast {
			t.Error("C should update last pointer on each addition")
		}
	})
}
