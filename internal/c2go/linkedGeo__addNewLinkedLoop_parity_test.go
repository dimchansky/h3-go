//go:build cgo && c2go

package c2go

import (
	"testing"
)

func Test_addLinkedLoop_parity(t *testing.T) {
	t.Run("empty polygon", func(t *testing.T) {
		// Test Go implementation
		goPolygon := &LinkedGeoPolygon{First: nil, Last: nil, Next: nil}
		goLoop := &LinkedGeoLoop{First: nil, Last: nil, Next: nil}
		goResult := addLinkedLoop(goPolygon, goLoop)
		
		// Check Go behavior
		goReturnsLoop := (goResult == goLoop)
		goSetsFirst := (goPolygon.First == goLoop)
		goSetsLast := (goPolygon.Last == goLoop)
		
		// Test C implementation behavior
		cReturnsLoop, cSetsFirst, cSetsLast, _ := addLinkedLoopC(true, false)
		
		// Compare behaviors
		if goReturnsLoop != cReturnsLoop {
			t.Errorf("Return behavior mismatch: Go=%v, C=%v", goReturnsLoop, cReturnsLoop)
		}
		if goSetsFirst != cSetsFirst {
			t.Errorf("First pointer behavior mismatch: Go=%v, C=%v", goSetsFirst, cSetsFirst)
		}
		if goSetsLast != cSetsLast {
			t.Errorf("Last pointer behavior mismatch: Go=%v, C=%v", goSetsLast, cSetsLast)
		}
	})
	
	t.Run("polygon with existing loop", func(t *testing.T) {
		// Test Go implementation
		existingLoop := &LinkedGeoLoop{First: nil, Last: nil, Next: nil}
		goPolygon := &LinkedGeoPolygon{First: existingLoop, Last: existingLoop, Next: nil}
		newLoop := &LinkedGeoLoop{First: nil, Last: nil, Next: nil}
		goResult := addLinkedLoop(goPolygon, newLoop)
		
		// Check Go behavior
		goReturnsLoop := (goResult == newLoop)
		goKeepsFirst := (goPolygon.First == existingLoop)
		goSetsLast := (goPolygon.Last == newLoop)
		goLinksLoops := (existingLoop.Next == newLoop)
		
		// Test C implementation behavior
		cReturnsLoop, cKeepsFirst, cSetsLast, cLinksLoops := addLinkedLoopC(false, true)
		
		// Compare behaviors
		if goReturnsLoop != cReturnsLoop {
			t.Errorf("Return behavior mismatch: Go=%v, C=%v", goReturnsLoop, cReturnsLoop)
		}
		if goKeepsFirst != cKeepsFirst {
			t.Errorf("First pointer behavior mismatch: Go=%v, C=%v", goKeepsFirst, cKeepsFirst)
		}
		if goSetsLast != cSetsLast {
			t.Errorf("Last pointer behavior mismatch: Go=%v, C=%v", goSetsLast, cSetsLast)
		}
		if goLinksLoops != cLinksLoops {
			t.Errorf("Loop linking behavior mismatch: Go=%v, C=%v", goLinksLoops, cLinksLoops)
		}
	})
}

func Test_addNewLinkedLoop_parity(t *testing.T) {
	t.Run("empty polygon", func(t *testing.T) {
		// Test Go implementation
		goPolygon := &LinkedGeoPolygon{First: nil, Last: nil, Next: nil}
		goResult := addNewLinkedLoop(goPolygon)
		
		// Check Go behavior
		goCreatesLoop := (goResult != nil)
		goSetsFirst := (goPolygon.First == goResult)
		goSetsLast := (goPolygon.Last == goResult)
		
		// Test C implementation behavior
		cCreatesLoop, cSetsFirst, cSetsLast, _ := addNewLinkedLoopC(true, false)
		
		// Compare behaviors
		if goCreatesLoop != cCreatesLoop {
			t.Errorf("Loop creation behavior mismatch: Go=%v, C=%v", goCreatesLoop, cCreatesLoop)
		}
		if goSetsFirst != cSetsFirst {
			t.Errorf("First pointer behavior mismatch: Go=%v, C=%v", goSetsFirst, cSetsFirst)
		}
		if goSetsLast != cSetsLast {
			t.Errorf("Last pointer behavior mismatch: Go=%v, C=%v", goSetsLast, cSetsLast)
		}
	})
	
	t.Run("polygon with existing loop", func(t *testing.T) {
		// Test Go implementation
		existingLoop := &LinkedGeoLoop{First: nil, Last: nil, Next: nil}
		goPolygon := &LinkedGeoPolygon{First: existingLoop, Last: existingLoop, Next: nil}
		goResult := addNewLinkedLoop(goPolygon)
		
		// Check Go behavior
		goCreatesLoop := (goResult != nil)
		goKeepsFirst := (goPolygon.First == existingLoop)
		goSetsLast := (goPolygon.Last == goResult)
		goLinksLoops := (existingLoop.Next == goResult)
		
		// Test C implementation behavior
		cCreatesLoop, cKeepsFirst, cSetsLast, cLinksLoops := addNewLinkedLoopC(false, true)
		
		// Compare behaviors
		if goCreatesLoop != cCreatesLoop {
			t.Errorf("Loop creation behavior mismatch: Go=%v, C=%v", goCreatesLoop, cCreatesLoop)
		}
		if goKeepsFirst != cKeepsFirst {
			t.Errorf("First pointer behavior mismatch: Go=%v, C=%v", goKeepsFirst, cKeepsFirst)
		}
		if goSetsLast != cSetsLast {
			t.Errorf("Last pointer behavior mismatch: Go=%v, C=%v", goSetsLast, cSetsLast)
		}
		if goLinksLoops != cLinksLoops {
			t.Errorf("Loop linking behavior mismatch: Go=%v, C=%v", goLinksLoops, cLinksLoops)
		}
	})
}