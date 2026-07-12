// Command testinventory enforces case-level completeness of the upstream
// H3 C test-suite port.
//
// The upstreamdiff tool answers "which upstream test FILES changed"; this
// tool goes one level deeper: it extracts every named upstream test case —
// TEST(name) blocks in src/apps/testapps/*.c, add_h3_cli_test entries in
// tests/cli/*.txt, fuzzer/benchmark/filter/helper executables, shared app
// support sources, and input fixtures — and
// cross-checks them against the committed registry
// docs/upstream-test-inventory.csv, which records a reviewed disposition
// for each case.
//
//	go run ./tools/testinventory -h3ver 4.4.0            # report
//	go run ./tools/testinventory -h3ver 4.4.0 -verify    # exit 1 on problems
//	go run ./tools/testinventory -h3ver 4.4.0 -init      # skeleton rows for unreviewed cases
//
// Registry columns: kind,upstream,case,status,go_test,notes,source_sha256
//
//   - kind:     test | cli | fuzzer | benchmark | filter | helper | support | fixture | build | generated
//   - upstream: path relative to the upstream tree root
//   - case:     TEST name / CLI test name / fuzzer file base name
//   - status:   full     — behaviorally equivalent Go test exists
//     partial  — ported with a documented, justified reduction
//     indirect — covered via cgo parity tests or another suite
//     na       — not applicable in Go (justify in notes)
//     deferred — reviewed, consciously not ported yet (justify)
//     missing  — no coverage; fails -verify (transient state only)
//   - go_test:  comma-separated Go test identifiers (TestX, FuzzX, BenchmarkX,
//     optionally with /subtest suffix). Required for full/partial/
//     indirect; each named function must exist in a *_test.go file.
//   - notes:    free text; required for partial/na/deferred.
//   - source_sha256: digest of the containing upstream file; changes require
//     renewed case-level review even when names remain unchanged.
//
// The tool never edits Go code; it assists human review during upstream
// syncs (docs/public-api-architecture.md §10, step 5).
package main

import (
	"crypto/sha256"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type upstreamCase struct {
	kind     string // test, cli, fuzzer
	upstream string // e.g. src/apps/testapps/testGridDisk.c
	name     string // e.g. gridDisk_pentagon
	digest   string // SHA-256 of the upstream file containing this entry
}

type registryRow struct {
	upstreamCase
	status string
	goTest string
	notes  string
	line   int
}

func key(c upstreamCase) string { return c.kind + "|" + c.upstream + "|" + c.name }

// ---------------------------------------------------------------------------
// Upstream extraction
// ---------------------------------------------------------------------------

var (
	testCaseRe = regexp.MustCompile(`(?m)^\s*TEST\(([A-Za-z0-9_]+)\)`)
	cliTestRe  = regexp.MustCompile(`add_h3_cli_test\(\s*([A-Za-z0-9_]+)`)
)

func scanUpstream(root string) ([]upstreamCase, error) {
	var out []upstreamCase

	// 1. testapps: every TEST(name) in test*.c.
	testDir := filepath.Join(root, "src", "apps", "testapps")
	entries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w (run `make -C testref h3-source`)", testDir, err)
	}
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasPrefix(n, "test") || !strings.HasSuffix(n, ".c") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(testDir, n))
		if err != nil {
			return nil, err
		}
		rel := "src/apps/testapps/" + n
		ms := testCaseRe.FindAllStringSubmatch(string(data), -1)
		if len(ms) == 0 {
			// Suites whose main() is the test (e.g. input-file-driven tests)
			// still need a disposition; record a synthetic whole-file case.
			out = append(out, upstreamCase{kind: "test", upstream: rel, name: "(main)"})
			continue
		}
		// A few upstream files repeat a TEST name (copy-paste); keep every
		// occurrence distinct so each block gets its own disposition.
		seen := map[string]int{}
		for _, m := range ms {
			name := m[1]
			seen[name]++
			if k := seen[name]; k > 1 {
				name = fmt.Sprintf("%s#%d", name, k)
			}
			out = append(out, upstreamCase{kind: "test", upstream: rel, name: name})
		}
	}

	// 7. Authoritative build/registration definitions and the optional
	// generated country benchmark pipeline.
	out = append(out,
		upstreamCase{kind: "build", upstream: "CMakeLists.txt", name: "(test ecosystem definitions)"},
		upstreamCase{kind: "build", upstream: "CMakeTests.cmake", name: "(CTest registrations)"},
		upstreamCase{kind: "generated", upstream: "scripts/make_countries.js", name: "benchmarkCountries.c"},
	)
	// 4. Other executable test ecosystem components. Benchmarks are included
	// because several perform correctness setup/checks before timing. Filters
	// and mkRand helpers protect command/integration behavior even when their
	// disposition is "not applicable" to this library-only Go port.
	for _, spec := range []struct {
		kind, dir, prefix string
	}{
		{kind: "benchmark", dir: "src/apps/benchmarks"},
		{kind: "filter", dir: "src/apps/filters"},
		{kind: "helper", dir: "src/apps/testapps", prefix: "mkRand"},
	} {
		err = appendSourceFiles(root, spec.dir, spec.kind, spec.prefix, &out)
		if err != nil {
			return nil, err
		}
	}

	// 5. Shared test/app support sources. Headers matter here: they define the
	// assertion, argument parsing, benchmark, and AFL harness semantics used by
	// the executable suites.
	for _, dir := range []string{"src/apps/applib/include", "src/apps/applib/lib"} {
		entries, readErr := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
		if readErr != nil {
			return nil, readErr
		}
		for _, e := range entries {
			if e.IsDir() || (!strings.HasSuffix(e.Name(), ".c") && !strings.HasSuffix(e.Name(), ".h")) {
				continue
			}
			out = append(out, upstreamCase{kind: "support", upstream: dir + "/" + e.Name(), name: "(file)"})
		}
	}

	// 6. Checked-in inputs consumed by CLI/integration suites.
	fixtureDir := "tests/inputfiles"
	entries, err = os.ReadDir(filepath.Join(root, filepath.FromSlash(fixtureDir)))
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, upstreamCase{kind: "fixture", upstream: fixtureDir + "/" + e.Name(), name: "(file)"})
		}
	}

	// 2. CLI tests: add_h3_cli_test entries in tests/cli/*.txt.
	cliDir := filepath.Join(root, "tests", "cli")
	entries, err = os.ReadDir(cliDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".txt") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(cliDir, n))
		if err != nil {
			return nil, err
		}
		for _, m := range cliTestRe.FindAllStringSubmatch(string(data), -1) {
			out = append(out, upstreamCase{kind: "cli", upstream: "tests/cli/" + n, name: m[1]})
		}
	}

	// 3. Fuzzers: one harness per file.
	fuzzDir := filepath.Join(root, "src", "apps", "fuzzers")
	entries, err = os.ReadDir(fuzzDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".c") {
			continue
		}
		out = append(out, upstreamCase{
			kind:     "fuzzer",
			upstream: "src/apps/fuzzers/" + n,
			name:     strings.TrimSuffix(n, ".c"),
		})
	}

	for i := range out {
		data, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(out[i].upstream)))
		if readErr != nil {
			return nil, readErr
		}
		out[i].digest = fmt.Sprintf("%x", sha256.Sum256(data))
	}

	sort.Slice(out, func(i, j int) bool { return key(out[i]) < key(out[j]) })
	return out, nil
}

func appendSourceFiles(root, dir, kind, prefix string, out *[]upstreamCase) error {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".c") || !strings.HasPrefix(name, prefix) {
			continue
		}
		*out = append(*out, upstreamCase{kind: kind, upstream: dir + "/" + name, name: "(main)"})
	}
	return nil
}

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

var validStatus = map[string]bool{
	"full": true, "partial": true, "indirect": true,
	"na": true, "deferred": true, "missing": true,
}

func loadRegistry(path string) ([]registryRow, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil // bootstrap: no registry yet, everything is unreviewed
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	r := csv.NewReader(f)
	r.FieldsPerRecord = 7
	recs, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	var rows []registryRow
	for i, rec := range recs {
		if i == 0 && rec[0] == "kind" { // header
			continue
		}
		rows = append(rows, registryRow{
			upstreamCase: upstreamCase{kind: rec[0], upstream: rec[1], name: rec[2]},
			status:       rec[3], goTest: rec[4], notes: rec[5], line: i + 1,
		})
		rows[len(rows)-1].digest = rec[6]
	}
	return rows, nil
}

// ---------------------------------------------------------------------------
// Go test reference resolution
// ---------------------------------------------------------------------------

// goTestDecls returns the set of Test/Fuzz/Benchmark/Example function names
// declared in *_test.go files in repository Go packages.
func goTestDecls(repoRoot string) (map[string]bool, error) {
	declRe := regexp.MustCompile(`(?m)^func ((?:Test|Fuzz|Benchmark|Example)\w*)\s*\(`)
	decls := map[string]bool{}
	err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != repoRoot && (entry.Name() == ".git" || entry.Name() == "testref") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range declRe.FindAllStringSubmatch(string(data), -1) {
			decls[m[1]] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return decls, nil
}

func resolveGoTest(ref string, decls map[string]bool) (missing []string) {
	for part := range strings.SplitSeq(ref, ",") {
		part = strings.TrimSpace(part)
		if part == "" || part == "-" {
			continue
		}
		fn := part
		if idx := strings.Index(fn, "/"); idx >= 0 {
			fn = fn[:idx]
		}
		if !decls[fn] {
			missing = append(missing, fn)
		}
	}
	return missing
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	h3ver := flag.String("h3ver", "4.4.0", "upstream H3 version under testref/")
	upstreamRoot := flag.String("upstream", "", "explicit upstream tree path (overrides -h3ver)")
	repo := flag.String("repo", ".", "repository root containing the Go port")
	registryPath := flag.String("registry", "docs/upstream-test-inventory.csv", "registry CSV path (relative to -repo)")
	verify := flag.Bool("verify", false, "exit 1 on unreviewed/stale/missing/invalid entries")
	initMode := flag.Bool("init", false, "print skeleton CSV rows for unreviewed upstream cases")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "testinventory checks that every upstream test-ecosystem entry has a reviewed disposition in docs/upstream-test-inventory.csv.")
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/testinventory [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()

	root := *upstreamRoot
	if root == "" {
		root = filepath.Join(*repo, "testref", "h3-"+*h3ver)
	}

	cases, err := scanUpstream(root)
	check(err)
	rows, err := loadRegistry(filepath.Join(*repo, *registryPath))
	check(err)
	decls, err := goTestDecls(*repo)
	check(err)

	regByKey := map[string]registryRow{}
	var problems []string
	for _, r := range rows {
		k := key(r.upstreamCase)
		if _, dup := regByKey[k]; dup {
			problems = append(problems, fmt.Sprintf("duplicate registry row (line %d): %s", r.line, k))
		}
		regByKey[k] = r
		if !validStatus[r.status] {
			problems = append(problems, fmt.Sprintf("invalid status %q (line %d): %s", r.status, r.line, k))
		}
		switch r.status {
		case "full", "partial", "indirect":
			if strings.TrimSpace(strings.Trim(r.goTest, "-")) == "" {
				problems = append(problems, fmt.Sprintf("status %s requires go_test (line %d): %s", r.status, r.line, k))
			}
		}
		switch r.status {
		case "partial", "na", "deferred":
			if strings.TrimSpace(r.notes) == "" {
				problems = append(problems, fmt.Sprintf("status %s requires notes (line %d): %s", r.status, r.line, k))
			}
		}
		if miss := resolveGoTest(r.goTest, decls); len(miss) > 0 {
			problems = append(problems, fmt.Sprintf("go_test not found (line %d): %s -> %s", r.line, k, strings.Join(miss, ", ")))
		}
	}

	upstreamKeys := map[string]bool{}
	var unreviewed []upstreamCase
	for _, c := range cases {
		upstreamKeys[key(c)] = true
		if registered, ok := regByKey[key(c)]; !ok {
			unreviewed = append(unreviewed, c)
		} else if registered.digest != c.digest {
			problems = append(problems, fmt.Sprintf("upstream file changed since review (line %d): %s", registered.line, key(c)))
		}
	}
	var stale []registryRow
	for _, r := range rows {
		if !upstreamKeys[key(r.upstreamCase)] {
			stale = append(stale, r)
		}
	}
	var open []registryRow
	for _, r := range rows {
		if r.status == "missing" {
			open = append(open, r)
		}
	}

	if *initMode {
		w := csv.NewWriter(os.Stdout)
		for _, c := range unreviewed {
			_ = w.Write([]string{c.kind, c.upstream, c.name, "missing", suggestGoTest(c, decls, *repo), "", c.digest})
		}
		w.Flush()
		return
	}

	// Report.
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.status]++
	}
	fmt.Printf("# Upstream test inventory: %s\n\n", filepath.Base(root))
	fmt.Printf("Upstream inventory entries: %d. Registry rows: %d.\n\n", len(cases), len(rows))
	fmt.Println("Entries by kind:")
	for _, kind := range []string{"test", "cli", "fuzzer", "benchmark", "filter", "helper", "support", "fixture", "build", "generated"} {
		fmt.Printf("  %-9s %d\n", kind, countKind(cases, kind))
	}
	fmt.Println()
	fmt.Printf("Dispositions: full %d, partial %d, indirect %d, na %d, deferred %d, missing %d.\n\n",
		counts["full"], counts["partial"], counts["indirect"], counts["na"], counts["deferred"], counts["missing"])

	fail := false
	section := func(title string, items []string) {
		if len(items) == 0 {
			return
		}
		fail = true
		fmt.Printf("## %s (%d)\n\n", title, len(items))
		for _, it := range items {
			fmt.Printf("- %s\n", it)
		}
		fmt.Println()
	}
	var unrevStr, staleStr, openStr []string
	for _, c := range unreviewed {
		unrevStr = append(unrevStr, fmt.Sprintf("`%s` :: `%s` — no registry row; review and add a disposition", c.upstream, c.name))
	}
	for _, r := range stale {
		staleStr = append(staleStr, fmt.Sprintf("`%s` :: `%s` (line %d) — not in upstream tree; remove or fix", r.upstream, r.name, r.line))
	}
	for _, r := range open {
		openStr = append(openStr, fmt.Sprintf("`%s` :: `%s` (line %d) — port or document", r.upstream, r.name, r.line))
	}
	section("Unreviewed upstream cases", unrevStr)
	section("Stale registry rows", staleStr)
	section("Registry integrity problems", problems)
	section("Open gaps (status=missing)", openStr)

	if !fail {
		fmt.Println("OK: every upstream test case has a reviewed disposition.")
	}
	if *verify && fail {
		os.Exit(1)
	}
}

func suggestGoTest(c upstreamCase, decls map[string]bool, repo string) string {
	want := normalizedTestName(strings.TrimSuffix(c.name, "#2"))
	if want == "" || want == "main" || want == "file" {
		return ""
	}
	var matches []string
	for decl := range decls {
		candidate := strings.TrimPrefix(decl, "Test")
		if normalizedTestName(candidate) == want {
			matches = append(matches, decl)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}

	// Disambiguate common case names (invalidInputs, roundtrip, etc.) using
	// the provenance header in the translated Go file, then allow an
	// idiomatic suite prefix on the Go test name.
	base := filepath.Base(c.upstream)
	entries, err := os.ReadDir(repo)
	if err != nil {
		return ""
	}
	declRe := regexp.MustCompile(`(?m)^func ((?:Test|Fuzz|Benchmark|Example)\w*)\s*\(`)
	var associated []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		data, readErr := os.ReadFile(filepath.Join(repo, e.Name()))
		if readErr != nil || !strings.Contains(string(data), base) {
			continue
		}
		for _, m := range declRe.FindAllStringSubmatch(string(data), -1) {
			candidate := normalizedTestName(strings.TrimPrefix(m[1], "Test"))
			if candidate == want || strings.HasSuffix(candidate, want) {
				associated = append(associated, m[1])
			}
		}
	}
	if len(associated) == 1 {
		return associated[0]
	}
	return ""
}

func normalizedTestName(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func countKind(cases []upstreamCase, kind string) int {
	n := 0
	for _, c := range cases {
		if c.kind == kind {
			n++
		}
	}
	return n
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
