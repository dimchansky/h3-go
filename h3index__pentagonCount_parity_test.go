//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_pentagonCount_ParityWithC(t *testing.T) {
	if pentagonCount() != pentagonCountC() {
		t.Fatalf("pentagonCount mismatch: go=%d c=%d", pentagonCount(), pentagonCountC())
	}
}
