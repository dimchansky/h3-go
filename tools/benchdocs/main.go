// Command benchdocs generates and verifies the selected benchmark excerpts
// in README.md and docs/benchmarks/README.md from committed benchmark CSV and
// metadata artifacts. It is intentionally offline and dependency-free.
//
// Usage:
//
//	go run ./tools/benchdocs -write   # refresh generated Markdown blocks
//	go run ./tools/benchdocs -verify  # fail when artifacts and docs drift
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	readmePath = "README.md"
	benchPath  = "docs/benchmarks/README.md"

	readmeBegin = "<!-- BEGIN GENERATED: benchdocs README (run `make gen-benchdocs`) -->"
	readmeEnd   = "<!-- END GENERATED: benchdocs README -->"
	benchBegin  = "<!-- BEGIN GENERATED: benchdocs details (run `make gen-benchdocs`) -->"
	benchEnd    = "<!-- END GENERATED: benchdocs details -->"
)

type metric struct {
	pure, uber, warm float64
	delta            string
	hasWarm          bool
}

type result struct {
	dir, goos, goarch, cpu string
	meta                   map[string]string
	metrics                map[string]map[string]metric
}

type selected struct {
	name, label string
	warm        bool
}

var selections = []selected{
	{"Resolution", "`Cell.Resolution`", false},
	{"CellToParent/res=9to7", "`Cell.Parent` (res 9→7)", false},
	{"LatLngToCell/res=9", "`LatLngToCell` (res 9)", false},
	{"CellToBoundary/res=9", "`Cell.Boundary` (res 9)", false},
	{"GridDisk/k=5", "`GridDisk` (k=5)", true},
	{"Compact/set=sf9", "`CompactCells` (1,253 cells)", true},
	{"PolygonToCells/poly=sf/res=9", "`PolygonToCells` (SF, res 9)", true},
	{"CellsToMultiPolygon/n=331", "`CellsToMultiPolygon` (331 cells)", false},
	{"ServiceWorkload/pts=256", "service workload (256 points)", true},
}

var gomaxSuffix = regexp.MustCompile(`-\d+$`)

func main() {
	repo := flag.String("repo", ".", "repository root")
	write := flag.Bool("write", false, "rewrite generated Markdown blocks")
	verify := flag.Bool("verify", false, "verify generated Markdown blocks")
	flag.Parse()
	if *write == *verify {
		fmt.Fprintln(os.Stderr, "benchdocs: exactly one of -write or -verify is required")
		os.Exit(2)
	}
	if err := run(*repo, *write); err != nil {
		fmt.Fprintln(os.Stderr, "benchdocs:", err)
		os.Exit(1)
	}
}

func run(repo string, write bool) error {
	var results []result
	for _, dir := range []string{"darwin-arm64", "linux-amd64"} {
		r, err := loadResult(filepath.Join(repo, "docs/benchmarks", dir), dir)
		if err != nil {
			return err
		}
		results = append(results, r)
	}

	readme := renderREADME(results)
	details := renderDetails(results)
	targets := []struct {
		path, begin, end, generated string
	}{
		{filepath.Join(repo, readmePath), readmeBegin, readmeEnd, readme},
		{filepath.Join(repo, benchPath), benchBegin, benchEnd, details},
	}
	for _, target := range targets {
		data, err := os.ReadFile(target.path)
		if err != nil {
			return err
		}
		current, err := generatedSection(string(data), target.begin, target.end)
		if err != nil {
			return fmt.Errorf("%s: %w", target.path, err)
		}
		if write {
			updated := strings.Replace(string(data), current, target.generated, 1)
			if err := os.WriteFile(target.path, []byte(updated), 0o644); err != nil {
				return err
			}
			continue
		}
		if strings.TrimSpace(current) != strings.TrimSpace(target.generated) {
			return fmt.Errorf("%s benchmark excerpt is stale; run `make gen-benchdocs`", target.path)
		}
	}
	if write {
		fmt.Println("benchdocs: generated benchmark excerpts updated")
	} else {
		fmt.Println("benchdocs: OK (README excerpts match both committed artifact sets)")
	}
	return nil
}

func loadResult(path, dir string) (result, error) {
	r := result{dir: dir, meta: map[string]string{}, metrics: map[string]map[string]metric{}}
	metaData, err := os.ReadFile(filepath.Join(path, "metadata.txt"))
	if err != nil {
		return r, err
	}
	for _, line := range strings.Split(string(metaData), "\n") {
		key, value, ok := strings.Cut(line, ": ")
		if ok {
			r.meta[key] = strings.TrimSpace(value)
		}
	}
	for _, key := range []string{"date_utc", "repo_commit", "go_version", "uber_h3_go", "pure_go_h3_target", "cpu", "cc", "bench_flags", "environment"} {
		if r.meta[key] == "" {
			return r, fmt.Errorf("%s/metadata.txt: missing %s", dir, key)
		}
	}
	// The refreshed local artifact uses the expanded metadata format. Keep the
	// older Linux artifact byte-for-byte stable until its next justified rerun.
	if dir == "darwin-arm64" {
		for _, key := range []string{"memory_bytes", "cgo_cppflags", "cgo_cxxflags", "cgo_ldflags", "gogccflags"} {
			if _, ok := r.meta[key]; !ok {
				return r, fmt.Errorf("%s/metadata.txt: missing %s", dir, key)
			}
		}
	}

	f, err := os.Open(filepath.Join(path, "benchstat.csv"))
	if err != nil {
		return r, err
	}
	defer func() { _ = f.Close() }()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return r, err
	}
	metricName := ""
	for _, record := range records {
		if len(record) == 0 {
			continue
		}
		line := strings.Join(record, ",")
		switch {
		case strings.HasPrefix(line, "goos: "):
			r.goos = strings.TrimSpace(strings.TrimPrefix(line, "goos: "))
		case strings.HasPrefix(line, "goarch: "):
			r.goarch = strings.TrimSpace(strings.TrimPrefix(line, "goarch: "))
		case strings.HasPrefix(line, "cpu: "):
			r.cpu = strings.TrimSpace(strings.TrimPrefix(line, "cpu: "))
		case len(record) > 1 && record[0] == "" && (record[1] == "sec/op" || record[1] == "B/op" || record[1] == "allocs/op"):
			metricName = record[1]
		case metricName != "" && record[0] != "" && record[0] != "geomean":
			if len(record) < 7 {
				return r, fmt.Errorf("%s/benchstat.csv: short %s row %q", dir, metricName, record[0])
			}
			name := gomaxSuffix.ReplaceAllString(record[0], "")
			pure, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				return r, fmt.Errorf("%s %s pure: %w", dir, name, err)
			}
			uber, err := strconv.ParseFloat(record[3], 64)
			if err != nil {
				return r, fmt.Errorf("%s %s uber: %w", dir, name, err)
			}
			m := metric{pure: pure, uber: uber, delta: record[5]}
			if len(record) > 11 && record[11] != "" {
				m.warm, err = strconv.ParseFloat(record[11], 64)
				if err != nil {
					return r, fmt.Errorf("%s %s warm: %w", dir, name, err)
				}
				m.hasWarm = true
			}
			if r.metrics[name] == nil {
				r.metrics[name] = map[string]metric{}
			}
			r.metrics[name][metricName] = m
		}
	}

	wantOS, wantArch, ok := strings.Cut(dir, "-")
	if !ok || r.goos != wantOS || r.goarch != wantArch {
		return r, fmt.Errorf("%s: benchstat platform is %s/%s", dir, r.goos, r.goarch)
	}
	if !strings.Contains(r.meta["go_version"], r.goos+"/"+r.goarch) {
		return r, fmt.Errorf("%s: metadata go_version disagrees with benchstat platform", dir)
	}
	if strings.TrimSpace(r.meta["cpu"]) != r.cpu {
		return r, fmt.Errorf("%s: metadata CPU %q disagrees with benchstat CPU %q", dir, r.meta["cpu"], r.cpu)
	}
	for _, selected := range selections {
		metrics := r.metrics[selected.name]
		for _, name := range []string{"sec/op", "B/op", "allocs/op"} {
			if _, ok := metrics[name]; !ok {
				return r, fmt.Errorf("%s: selected benchmark %q metric %s disappeared", dir, selected.name, name)
			}
		}
	}
	return r, nil
}

func renderREADME(results []result) string {
	var b strings.Builder
	b.WriteString(readmeBegin + "\n")
	for _, r := range results {
		fmt.Fprintf(&b, "### %s — %s\n\n", platformTitle(r), r.dir)
		fmt.Fprintf(&b, "%s; %s; %s; %s; repository `%s`. ", r.meta["go_version"], r.meta["cc"], r.meta["uber_h3_go"], r.meta["pure_go_h3_target"], shortCommit(r.meta["repo_commit"]))
		if r.dir == "linux-amd64" {
			b.WriteString("This is a shared GitHub Actions runner, so small deltas are noisier. ")
		}
		fmt.Fprintf(&b, "[Metadata](docs/benchmarks/%s/metadata.txt) · [full benchstat table](docs/benchmarks/%s/benchstat.txt) · [raw output](docs/benchmarks/%s/bench-raw.txt).\n\n", r.dir, r.dir, r.dir)
		b.WriteString("| Operation | this library | warm `Append*` | uber/h3-go v4.4.1 |\n")
		b.WriteString("|---|---:|---:|---:|\n")
		for _, selected := range selections {
			metrics := r.metrics[selected.name]
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", selected.label, formatCell(metrics, "pure"), formatWarm(metrics, selected.warm), formatCell(metrics, "uber"))
		}
		b.WriteString("\n")
	}
	b.WriteString("**Measured fact:** several results reverse between environments: `LatLngToCell`, `Cell.Boundary`, and `PolygonToCells` favor the binding on the M1 Max but pure Go on the Linux runner; `CompactCells` and `CellsToMultiPolygon` favor the binding in both. **Plausible explanation, not isolated by these measurements:** cgo-call cost, compiler code generation, and CPU microarchitecture all contribute. The artifacts do not identify a single cause, and absolute timings must never be compared across the two machines.\n")
	b.WriteString(readmeEnd)
	return b.String()
}

func renderDetails(results []result) string {
	var b strings.Builder
	b.WriteString(benchBegin + "\n")
	for _, r := range results {
		fmt.Fprintf(&b, "### %s (`%s`)\n\n", platformTitle(r), r.dir)
		fmt.Fprintf(&b, "Run `%s` at repository `%s`; %s; %s; `%s`.\n\n", r.meta["date_utc"], shortCommit(r.meta["repo_commit"]), r.meta["go_version"], r.meta["cc"], r.meta["bench_flags"])
		b.WriteString("| Selected operation | pure sec/op | uber sec/op | uber vs pure |\n")
		b.WriteString("|---|---:|---:|---:|\n")
		for _, selected := range selections {
			m := r.metrics[selected.name]["sec/op"]
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", selected.label, formatTime(m.pure), formatTime(m.uber), m.delta)
		}
		b.WriteString("\n")
	}
	b.WriteString("These tables are generated from the committed `benchstat.csv` files. Positive “uber vs pure” values mean the binding took longer; negative values mean it was faster. See the full tables for confidence intervals and p-values.\n")
	b.WriteString(benchEnd)
	return b.String()
}

func formatCell(metrics map[string]metric, impl string) string {
	time, bytes, allocs := metrics["sec/op"], metrics["B/op"], metrics["allocs/op"]
	if impl == "uber" {
		return fmt.Sprintf("%s · %s · %s", formatTime(time.uber), formatBytes(bytes.uber), formatAllocs(allocs.uber))
	}
	return fmt.Sprintf("%s · %s · %s", formatTime(time.pure), formatBytes(bytes.pure), formatAllocs(allocs.pure))
}

func formatWarm(metrics map[string]metric, selected bool) string {
	if !selected || !metrics["sec/op"].hasWarm {
		return "—"
	}
	return fmt.Sprintf("%s · %s · %s", formatTime(metrics["sec/op"].warm), formatBytes(metrics["B/op"].warm), formatAllocs(metrics["allocs/op"].warm))
}

func formatTime(seconds float64) string {
	switch {
	case seconds < 1e-6:
		return fmt.Sprintf("%.3g ns", seconds*1e9)
	case seconds < 1e-3:
		return fmt.Sprintf("%.4g µs", seconds*1e6)
	default:
		return fmt.Sprintf("%.4g ms", seconds*1e3)
	}
}

func formatBytes(bytes float64) string {
	if bytes >= 1024 {
		return fmt.Sprintf("%.4g KiB", bytes/1024)
	}
	return fmt.Sprintf("%.4g B", bytes)
}

func formatAllocs(allocs float64) string {
	if math.Trunc(allocs) == allocs {
		return fmt.Sprintf("%.0f allocs", allocs)
	}
	return fmt.Sprintf("%.3g allocs", allocs)
}

func platformTitle(r result) string {
	if r.dir == "darwin-arm64" {
		return "Apple M1 Max"
	}
	return "GitHub Actions AMD EPYC 7763"
}

func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

func generatedSection(doc, begin, end string) (string, error) {
	start := strings.Index(doc, begin)
	if start < 0 {
		return "", fmt.Errorf("missing marker %q", begin)
	}
	finishRel := strings.Index(doc[start:], end)
	if finishRel < 0 {
		return "", fmt.Errorf("missing marker %q", end)
	}
	finish := start + finishRel + len(end)
	return doc[start:finish], nil
}
