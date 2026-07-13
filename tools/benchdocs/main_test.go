package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadResult(t *testing.T) {
	if _, err := loadResult("../../docs/benchmarks/darwin-arm64", "darwin-arm64"); err != nil {
		t.Fatalf("committed darwin artifact: %v", err)
	}
}

func TestCompleteResultsCoverCatalog(t *testing.T) {
	var results []result
	for _, dir := range []string{"darwin-arm64", "linux-amd64"} {
		r, err := loadResult(filepath.Join("../../docs/benchmarks", dir), dir)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, r)
	}

	doc := renderFullResults(results)
	for _, scenario := range scenarios {
		if count := strings.Count(doc, "| "+scenario.label+" |"); count < 4 {
			t.Errorf("scenario %q appears in only %d generated table rows, want at least 4", scenario.name, count)
		}
	}
	for _, scenario := range memoryScenarios {
		if count := strings.Count(doc, "| "+scenario.label+" |"); count != 2 {
			t.Errorf("process-memory scenario %q appears in %d generated table rows, want 2", scenario.name, count)
		}
	}
	for _, want := range []string{
		"For each point: `LatLngToCell` (res 9), `GridDisk` (k 1), then `Parent` (res 7)",
		"🟢 +",
		"🔴 −",
		"<summary>Run provenance and raw artifacts</summary>",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("complete results missing %q", want)
		}
	}
}

func TestLoadResultRejectsDrift(t *testing.T) {
	tests := []struct {
		name        string
		metadataOld string
		metadataNew string
		csvOld      string
		csvNew      string
		memoryOld   string
		memoryNew   string
		want        string
	}{
		{
			name:   "catalog benchmark disappeared",
			csvOld: "Resolution-10",
			csvNew: "ResolutionGone-10",
			want:   "benchmark scenario \"Resolution\" metric sec/op disappeared",
		},
		{
			name:        "CPU metadata mismatch",
			metadataOld: "cpu: Apple M1 Max",
			metadataNew: "cpu: Other CPU",
			want:        "metadata CPU",
		},
		{
			name:        "platform metadata mismatch",
			metadataOld: "darwin/arm64",
			metadataNew: "linux/amd64",
			want:        "metadata go_version disagrees",
		},
		{
			name:        "expanded metadata missing",
			metadataOld: "memory_bytes:",
			metadataNew: "memory_bytes_removed:",
			want:        "missing memory_bytes",
		},
		{
			name:      "process-memory scenario disappeared",
			memoryOld: "polyfill-large",
			memoryNew: "polyfill-large-renamed",
			want:      "process-memory scenario \"polyfill-large\" disappeared",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := copyArtifact(t)
			if tc.metadataOld != "" {
				replaceFile(t, filepath.Join(dir, "metadata.txt"), tc.metadataOld, tc.metadataNew)
			}
			if tc.csvOld != "" {
				replaceFile(t, filepath.Join(dir, "benchstat.csv"), tc.csvOld, tc.csvNew)
			}
			if tc.memoryOld != "" {
				replaceFile(t, filepath.Join(dir, "memory.tsv"), tc.memoryOld, tc.memoryNew)
			}
			_, err := loadResult(dir, "darwin-arm64")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("loadResult error = %v, want substring %q", err, tc.want)
			}
		})
	}
}

func copyArtifact(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"metadata.txt", "benchstat.csv", "memory.tsv"} {
		data, err := os.ReadFile(filepath.Join("../../docs/benchmarks/darwin-arm64", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func replaceFile(t *testing.T, path, old, replacement string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.ReplaceAll(string(data), old, replacement)
	if updated == string(data) {
		t.Fatalf("%q not found in %s", old, path)
	}
	if err := os.WriteFile(path, []byte(updated), 0o600); err != nil {
		t.Fatal(err)
	}
}
