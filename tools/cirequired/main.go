// Command cirequired is the aggregate merge gate behind the "CI / required"
// status check (the `required` job in .github/workflows/ci.yml). It reads
// the final result of every gated CI job and passes only when the
// combination matches the required-check truth table exactly.
//
// Inputs are environment variables (the workflow passes them via `env:`,
// never interpolated into a shell command line):
//
//	CIREQUIRED_NEEDS  ${{ toJSON(needs) }} — one entry per gated job
//	                  (changes, docs, fast, race, api-gates, parity).
//	CIREQUIRED_EVENT  ${{ github.event_name }} — pull_request or push.
//	CIREQUIRED_CODE   ${{ needs.changes.outputs.code }} — the docs-only
//	                  classifier verdict ("true" means a code change).
//
// Truth table (one row per event/classifier combination):
//
//	code=true  pull_request  changes, docs, fast, race, api-gates, parity = success
//	code=true  push          race = skipped (PR-only merge gate); the rest = success
//	code=false pull_request  changes, docs = success; fast, race, api-gates, parity = skipped
//	code=false push          changes, docs = success; fast, race, api-gates, parity = skipped
//
// Every other combination fails: unknown events, invalid classifier
// values, malformed or duplicate-key JSON, unknown or missing jobs,
// results outside the four GitHub Actions values, and any failed,
// canceled, unexpectedly skipped, or unexpectedly run job. A generic
// "success or skipped" acceptance is deliberately absent: a broken job
// condition must fail this gate, not slip a change past it.
//
// The tool is read-only and dependency-free; it exits 1 with a per-job
// report on any violation.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// gatedJobs is the exact set of ci.yml jobs the gate evaluates; the
// workflow's `needs:` list must stay in lockstep with it.
var gatedJobs = []string{"changes", "docs", "fast", "race", "api-gates", "parity"}

// resultCancelled is the literal job result GitHub Actions produces for a
// canceled job (the API uses the British spelling).
//
//nolint:misspell // must match the exact GitHub Actions result string
const resultCancelled = "cancelled"

// validResults are the only job results GitHub Actions produces.
var validResults = map[string]bool{
	"success":       true,
	"failure":       true,
	resultCancelled: true,
	"skipped":       true,
}

func main() {
	violations := evaluate(
		os.Getenv("CIREQUIRED_NEEDS"),
		os.Getenv("CIREQUIRED_EVENT"),
		os.Getenv("CIREQUIRED_CODE"),
	)
	if len(violations) > 0 {
		for _, v := range violations {
			fmt.Fprintln(os.Stderr, v)
		}
		fmt.Fprintf(os.Stderr, "cirequired: %d violation(s) of the required-check truth table\n", len(violations))
		os.Exit(1)
	}
	fmt.Println("cirequired: all gated jobs match the required-check truth table")
}

// evaluate checks the needs results against the truth-table row selected by
// event and code, and returns one violation string per problem (empty means
// the gate passes).
func evaluate(needsJSON, event, code string) []string {
	var violations []string

	expected, err := expectedResults(event, code)
	if err != nil {
		violations = append(violations, err.Error())
	}

	results, err := parseNeeds(needsJSON)
	if err != nil {
		return append(violations, err.Error())
	}

	for _, job := range gatedJobs {
		got, ok := results[job]
		if !ok {
			violations = append(violations, fmt.Sprintf("job %q: missing from the needs context", job))
			continue
		}
		if !validResults[got] {
			violations = append(violations, fmt.Sprintf("job %q: unknown result %q", job, got))
			continue
		}
		if expected == nil {
			continue // the row itself is invalid; that violation is already recorded
		}
		if want := expected[job]; got != want {
			violations = append(violations,
				fmt.Sprintf("job %q: result %q, want %q (event=%s, code=%s)", job, got, want, event, code))
		}
	}

	var unknown []string
	for job := range results {
		if !slices.Contains(gatedJobs, job) {
			unknown = append(unknown, fmt.Sprintf("job %q: not a gated job; update tools/cirequired and ci.yml together", job))
		}
	}
	slices.Sort(unknown)
	return append(violations, unknown...)
}

// expectedResults returns the truth-table row for the given event and
// classifier output, or an error when the combination is not in the table.
func expectedResults(event, code string) (map[string]string, error) {
	if event != "pull_request" && event != "push" {
		return nil, fmt.Errorf("event %q: not in the required-check truth table", event)
	}
	if code != "true" && code != "false" {
		return nil, fmt.Errorf("classifier output %q: want \"true\" or \"false\"", code)
	}
	expected := map[string]string{
		"changes":   "success",
		"docs":      "success",
		"fast":      "skipped",
		"race":      "skipped",
		"api-gates": "skipped",
		"parity":    "skipped",
	}
	if code == "true" {
		expected["fast"] = "success"
		expected["api-gates"] = "success"
		expected["parity"] = "success"
		if event == "pull_request" {
			expected["race"] = "success" // race is a PR-only merge gate; skipped on push is expected
		}
	}
	return expected, nil
}

// parseNeeds decodes the `needs` context JSON into job → result, rejecting
// anything but a single flat object with unique keys (encoding/json would
// silently keep the last duplicate, so keys are walked token by token).
func parseNeeds(data string) (map[string]string, error) {
	dec := json.NewDecoder(strings.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("needs JSON: %w", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, errors.New("needs JSON: top-level value must be an object")
	}
	results := make(map[string]string)
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("needs JSON: %w", err)
		}
		key, ok := keyTok.(string)
		if !ok {
			return nil, errors.New("needs JSON: object key is not a string")
		}
		if _, dup := results[key]; dup {
			return nil, fmt.Errorf("needs JSON: duplicate job %q", key)
		}
		var job struct {
			Result string `json:"result"`
		}
		if err := dec.Decode(&job); err != nil {
			return nil, fmt.Errorf("needs JSON: job %q: %w", key, err)
		}
		results[key] = job.Result
	}
	if _, err := dec.Token(); err != nil { // consume the closing brace
		return nil, fmt.Errorf("needs JSON: %w", err)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("needs JSON: trailing content after the object")
	}
	return results, nil
}
