//go:build !race

package h3

import "testing"

// assertMaxAllocsPerRun enforces production allocation budgets only in
// ordinary builds. Race instrumentation has a different escape/allocation
// profile and is covered by the companion race-tagged implementation.
func assertMaxAllocsPerRun(t *testing.T, name string, runs int, maxAllocs float64, f func()) {
	t.Helper()
	if got := testing.AllocsPerRun(runs, f); got > maxAllocs {
		t.Errorf("%s allocates %g/run, want <= %g", name, got, maxAllocs)
	}
}
