//go:build c2go

package c2go

import (
	"math"
	"testing"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func Test_vec3d__v3dNormalize_Parity(t *testing.T) {
	tests := []Vec3d{
		{X: 1, Y: 0, Z: 0},
		{X: 0, Y: -3, Z: 4},
		{X: 1.2345, Y: -6.789, Z: 0.001},
		{X: 0, Y: 0, Z: 0}, // zero magnitude edge case
	}
	const tol = 1e-15
	for _, v := range tests {
		gotGo := _v3dNormalize(&v)
		gotC := v3dNormalizeC(v)
		if !almostEqual(gotGo.X, gotC.X, tol) || !almostEqual(gotGo.Y, gotC.Y, tol) || !almostEqual(gotGo.Z, gotC.Z, tol) {
			t.Fatalf("_v3dNormalize mismatch for v=%+v: go=%+v c=%+v", v, gotGo, gotC)
		}
	}
}
