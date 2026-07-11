// Tests ported from H3 v4.4.0: src/apps/testapps/testGridDistanceInternal.c.
package h3

import (
	"testing"
)

func Test_ijkDistance(t *testing.T) {
	t.Parallel()

	z := coordIJK{I: 0, J: 0, K: 0}
	i := coordIJK{I: 1, J: 0, K: 0}
	ik := coordIJK{I: 1, J: 0, K: 1}
	ij := coordIJK{I: 1, J: 1, K: 0}
	j2 := coordIJK{I: 0, J: 2, K: 0}

	// Identity distance tests
	if dist := ijkDistance(&z, &z); dist != 0 {
		t.Errorf("Expected identity distance 0,0,0 to be 0, got %d", dist)
	}
	if dist := ijkDistance(&i, &i); dist != 0 {
		t.Errorf("Expected identity distance 1,0,0 to be 0, got %d", dist)
	}
	if dist := ijkDistance(&ik, &ik); dist != 0 {
		t.Errorf("Expected identity distance 1,0,1 to be 0, got %d", dist)
	}
	if dist := ijkDistance(&ij, &ij); dist != 0 {
		t.Errorf("Expected identity distance 1,1,0 to be 0, got %d", dist)
	}
	if dist := ijkDistance(&j2, &j2); dist != 0 {
		t.Errorf("Expected identity distance 0,2,0 to be 0, got %d", dist)
	}

	// Distance tests between different coordinates
	if dist := ijkDistance(&z, &i); dist != 1 {
		t.Errorf("Expected distance from 0,0,0 to 1,0,0 to be 1, got %d", dist)
	}
	if dist := ijkDistance(&z, &j2); dist != 2 {
		t.Errorf("Expected distance from 0,0,0 to 0,2,0 to be 2, got %d", dist)
	}
	if dist := ijkDistance(&z, &ik); dist != 1 {
		t.Errorf("Expected distance from 0,0,0 to 1,0,1 to be 1, got %d", dist)
	}
	if dist := ijkDistance(&i, &ik); dist != 1 {
		t.Errorf("Expected distance from 1,0,0 to 1,0,1 to be 1, got %d", dist)
	}
	if dist := ijkDistance(&ik, &j2); dist != 3 {
		t.Errorf("Expected distance from 1,0,1 to 0,2,0 to be 3, got %d", dist)
	}
	if dist := ijkDistance(&ij, &ik); dist != 2 {
		t.Errorf("Expected distance from 1,1,0 to 1,0,1 to be 2, got %d", dist)
	}
}
