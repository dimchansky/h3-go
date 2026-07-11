// Tests ported from testVec3dInternal.c
package h3

import (
	"math"
	"testing"
)

func Test_pointSquareDist(t *testing.T) {
	t.Parallel()

	v1 := &vec3d{0, 0, 0}
	v2 := &vec3d{1, 0, 0}
	v3 := &vec3d{0, 1, 1}
	v4 := &vec3d{1, 1, 1}
	v5 := &vec3d{1, 1, 2}

	const dblEpsilon = 2.220446049250313e-16

	// Test distance to self is 0
	if math.Abs(_pointSquareDist(v1, v1)) >= dblEpsilon {
		t.Error("distance to self should be 0")
	}

	// Test distance to <1,0,0> is 1
	if math.Abs(_pointSquareDist(v1, v2)-1) >= dblEpsilon {
		t.Error("distance to <1,0,0> should be 1")
	}

	// Test distance to <0,1,1> is 2
	if math.Abs(_pointSquareDist(v1, v3)-2) >= dblEpsilon {
		t.Error("distance to <0,1,1> should be 2")
	}

	// Test distance to <1,1,1> is 3
	if math.Abs(_pointSquareDist(v1, v4)-3) >= dblEpsilon {
		t.Error("distance to <1,1,1> should be 3")
	}

	// Test distance to <1,1,2> is 6
	if math.Abs(_pointSquareDist(v1, v5)-6) >= dblEpsilon {
		t.Error("distance to <1,1,2> should be 6")
	}
}
