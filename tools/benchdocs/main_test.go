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

func TestLoadResultRejectsDrift(t *testing.T) {
	tests := []struct {
		name        string
		metadataOld string
		metadataNew string
		csvOld      string
		csvNew      string
		want        string
	}{
		{
			name:   "selected benchmark disappeared",
			csvOld: "Resolution-10",
			csvNew: "ResolutionGone-10",
			want:   "selected benchmark \"Resolution\" metric sec/op disappeared",
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
	for _, name := range []string{"metadata.txt", "benchstat.csv"} {
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
