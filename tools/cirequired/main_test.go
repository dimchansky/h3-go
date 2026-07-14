package main

import (
	"fmt"
	"strings"
	"testing"
)

// needs builds the CIREQUIRED_NEEDS JSON for the given job results,
// mirroring the shape of ${{ toJSON(needs) }} (extra fields included to
// prove they are ignored).
func needs(results map[string]string) string {
	var b strings.Builder
	b.WriteString("{")
	first := true
	for _, job := range gatedJobs {
		result, ok := results[job]
		if !ok {
			continue
		}
		if !first {
			b.WriteString(",")
		}
		first = false
		fmt.Fprintf(&b, "%q: {\"result\": %q, \"outputs\": {}}", job, result)
	}
	b.WriteString("}")
	return b.String()
}

// allSuccess is the code-PR row: every gated job ran and succeeded.
func allSuccess() map[string]string {
	return map[string]string{
		"changes": "success", "docs": "success", "fast": "success",
		"race": "success", "api-gates": "success", "parity": "success",
	}
}

// docsOnly is the docs row: classifier and docs ran, all Go jobs skipped.
func docsOnly() map[string]string {
	return map[string]string{
		"changes": "success", "docs": "success", "fast": "skipped",
		"race": "skipped", "api-gates": "skipped", "parity": "skipped",
	}
}

func TestEvaluateTruthTable(t *testing.T) {
	codePush := allSuccess()
	codePush["race"] = "skipped" // race is PR-only

	tests := []struct {
		name       string
		results    map[string]string
		event      string
		code       string
		violations int
	}{
		// The four valid rows.
		{name: "code PR", results: allSuccess(), event: "pull_request", code: "true"},
		{name: "code push", results: codePush, event: "push", code: "true"},
		{name: "docs-only PR", results: docsOnly(), event: "pull_request", code: "false"},
		{name: "docs-only push", results: docsOnly(), event: "push", code: "false"},

		// Row-selection failures.
		{name: "unknown event", results: allSuccess(), event: "workflow_dispatch", code: "true", violations: 1},
		{name: "empty event", results: allSuccess(), event: "", code: "true", violations: 1},
		{name: "invalid classifier output", results: allSuccess(), event: "pull_request", code: "yes", violations: 1},
		{name: "missing classifier output", results: allSuccess(), event: "pull_request", code: "", violations: 1},

		// Result failures on a code PR (incl. the matrix job aggregating a failed leg).
		{name: "fast matrix failure", results: override(allSuccess(), "fast", "failure"), event: "pull_request", code: "true", violations: 1},
		{name: "parity canceled", results: override(allSuccess(), "parity", resultCancelled), event: "pull_request", code: "true", violations: 1},
		{name: "race failure", results: override(allSuccess(), "race", "failure"), event: "pull_request", code: "true", violations: 1},
		{name: "docs failure", results: override(allSuccess(), "docs", "failure"), event: "pull_request", code: "true", violations: 1},
		{name: "unknown result value", results: override(allSuccess(), "fast", "timed_out"), event: "pull_request", code: "true", violations: 1},

		// Unexpected skips: a broken job condition must fail the gate.
		{name: "fast unexpectedly skipped on code PR", results: override(allSuccess(), "fast", "skipped"), event: "pull_request", code: "true", violations: 1},
		{name: "changes skipped", results: override(allSuccess(), "changes", "skipped"), event: "pull_request", code: "true", violations: 1},
		{name: "docs skipped on docs-only PR", results: override(docsOnly(), "docs", "skipped"), event: "pull_request", code: "false", violations: 1},
		{name: "race unexpectedly ran on push", results: override(codePushCopy(), "race", "success"), event: "push", code: "true", violations: 1},

		// Unexpected runs: Go jobs executing on a docs-only change.
		{name: "fast unexpectedly ran on docs-only PR", results: override(docsOnly(), "fast", "success"), event: "pull_request", code: "false", violations: 1},

		// Everything wrong at once still reports every violation.
		{name: "all failures on code PR", results: map[string]string{
			"changes": "failure", "docs": "failure", "fast": "failure",
			"race": "failure", "api-gates": "failure", "parity": "failure",
		}, event: "pull_request", code: "true", violations: 6},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluate(needs(tc.results), tc.event, tc.code)
			if len(got) != tc.violations {
				t.Fatalf("violations = %d, want %d:\n%s", len(got), tc.violations, strings.Join(got, "\n"))
			}
		})
	}
}

func TestEvaluateInputValidation(t *testing.T) {
	tests := []struct {
		name  string
		json  string
		event string
		code  string
	}{
		{name: "empty needs", json: "", event: "pull_request", code: "true"},
		{name: "malformed JSON", json: "{", event: "pull_request", code: "true"},
		{name: "not an object", json: `["changes"]`, event: "pull_request", code: "true"},
		{name: "trailing content", json: needs(allSuccess()) + "{}", event: "pull_request", code: "true"},
		{name: "duplicate job", json: `{"changes":{"result":"success"},"changes":{"result":"failure"}}`, event: "pull_request", code: "true"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := evaluate(tc.json, tc.event, tc.code); len(got) == 0 {
				t.Fatal("want at least one violation, got none")
			}
		})
	}
}

func TestEvaluateJobSetDrift(t *testing.T) {
	t.Run("missing job", func(t *testing.T) {
		partial := allSuccess()
		delete(partial, "race")
		got := evaluate(needs(partial), "pull_request", "true")
		if len(got) != 1 || !strings.Contains(got[0], `"race"`) {
			t.Fatalf("want one missing-race violation, got: %v", got)
		}
	})
	t.Run("unknown job", func(t *testing.T) {
		extra := needs(allSuccess())
		extra = strings.TrimSuffix(extra, "}") + `,"vulncheck":{"result":"success"}}`
		got := evaluate(extra, "pull_request", "true")
		if len(got) != 1 || !strings.Contains(got[0], `"vulncheck"`) {
			t.Fatalf("want one unknown-job violation, got: %v", got)
		}
	})
}

// override returns results with one job's result replaced.
func override(results map[string]string, job, result string) map[string]string {
	results[job] = result
	return results
}

// codePushCopy is the valid code-push row (race skipped).
func codePushCopy() map[string]string {
	r := allSuccess()
	r["race"] = "skipped"
	return r
}
