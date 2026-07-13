//go:build race

package h3

import "testing"

// assertMaxAllocsPerRun preserves behavioral and race coverage for the
// measured closure without applying a production allocation budget to a
// race-instrumented binary.
func assertMaxAllocsPerRun(t *testing.T, _ string, _ int, _ float64, f func()) {
	t.Helper()
	f()
}
