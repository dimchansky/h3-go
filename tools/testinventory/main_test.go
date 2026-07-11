package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanUpstreamDiscoversWholeEcosystemAndFingerprintsChanges(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"src/apps/testapps/testOne.c":        "TEST(alpha) {}\n",
		"src/apps/testapps/mkRandOne.c":      "int main(void) {}\n",
		"src/apps/fuzzers/fuzzerOne.c":       "int LLVMFuzzerTestOneInput(void) {}\n",
		"src/apps/benchmarks/benchmarkOne.c": "int main(void) {}\n",
		"src/apps/filters/filterOne.c":       "int main(void) {}\n",
		"src/apps/applib/include/test.h":     "#define TEST(x)\n",
		"src/apps/applib/lib/test.c":         "int globalTestCount;\n",
		"tests/cli/one.txt":                  "add_h3_cli_test(testCliOne \"one\")\n",
		"tests/inputfiles/one.txt":           "fixture\n",
		"CMakeLists.txt":                     "add_h3_test(testOne)\n",
		"CMakeTests.cmake":                   "add_h3_test(testOne)\n",
		"scripts/make_countries.js":          "// generator\n",
	}
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	first, err := scanUpstream(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 12 {
		t.Fatalf("scanUpstream returned %d entries; want 12", len(first))
	}
	for _, entry := range first {
		if entry.digest == "" {
			t.Fatalf("%s has no source digest", key(entry))
		}
	}

	testPath := filepath.Join(root, "src/apps/testapps/testOne.c")
	if err := os.WriteFile(testPath, []byte("TEST(alpha) { changed(); }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	second, err := scanUpstream(root)
	if err != nil {
		t.Fatal(err)
	}
	if firstDigest := digestFor(first, "test", "src/apps/testapps/testOne.c", "alpha"); firstDigest == digestFor(second, "test", "src/apps/testapps/testOne.c", "alpha") {
		t.Fatal("source digest did not change after an in-place assertion edit")
	}
}

func digestFor(cases []upstreamCase, kind, upstream, name string) string {
	for _, c := range cases {
		if c.kind == kind && c.upstream == upstream && c.name == name {
			return c.digest
		}
	}
	return ""
}
