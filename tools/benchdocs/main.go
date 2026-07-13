// Command benchdocs generates and verifies the benchmark scorecard in
// README.md and the complete comparison in docs/benchmarks/results.md from
// committed benchmark CSV and metadata artifacts. It is intentionally
// offline and dependency-free.
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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	readmePath  = "README.md"
	resultsPath = "docs/benchmarks/results.md"

	readmeBegin = "<!-- BEGIN GENERATED: benchdocs README (run `make gen-benchdocs`) -->"
	readmeEnd   = "<!-- END GENERATED: benchdocs README -->"
)

type measurement struct {
	value        float64
	ci, delta, p string
	present      bool
}

type metric struct {
	pure, uber, cold, warm measurement
}

type memoryMeasurement struct {
	peakRSSKB, heapAllocKB float64
	checksum               string
	present                bool
}

type memoryResult struct {
	pure, uber memoryMeasurement
}

type result struct {
	dir, goos, goarch, cpu string
	meta                   map[string]string
	metrics                map[string]map[string]metric
	memory                 map[string]memoryResult
	order                  []string
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
	readmeFile := filepath.Join(repo, readmePath)
	data, err := os.ReadFile(readmeFile)
	if err != nil {
		return err
	}
	current, err := generatedSection(string(data), readmeBegin, readmeEnd)
	if err != nil {
		return fmt.Errorf("%s: %w", readmeFile, err)
	}
	if write {
		updated := strings.Replace(string(data), current, readme, 1)
		if err := os.WriteFile(readmeFile, []byte(updated), 0o644); err != nil {
			return err
		}
	} else if strings.TrimSpace(current) != strings.TrimSpace(readme) {
		return fmt.Errorf("%s benchmark scorecard is stale; run `make gen-benchdocs`", readmeFile)
	}

	fullResults := renderFullResults(results)
	fullResultsFile := filepath.Join(repo, resultsPath)
	if write {
		if err := os.WriteFile(fullResultsFile, []byte(fullResults), 0o644); err != nil {
			return err
		}
	} else {
		committed, err := os.ReadFile(fullResultsFile)
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(committed)) != strings.TrimSpace(fullResults) {
			return fmt.Errorf("%s is stale; run `make gen-benchdocs`", fullResultsFile)
		}
	}
	if write {
		fmt.Println("benchdocs: README scorecard and complete results updated")
	} else {
		fmt.Println("benchdocs: OK (scorecard and complete results match both artifact sets)")
	}
	return nil
}

func loadResult(path, dir string) (result, error) {
	r := result{
		dir:     dir,
		meta:    map[string]string{},
		metrics: map[string]map[string]metric{},
		memory:  map[string]memoryResult{},
	}
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
	for _, key := range []string{"date_utc", "repo_commit", "go_version", "uber_h3_go", "pure_go_h3_target", "cpu", "cc", "bench_flags", "memprobe_iters", "environment"} {
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
			pure, err := parseMeasurement(record, 1, -1)
			if err != nil {
				return r, fmt.Errorf("%s %s pure: %w", dir, name, err)
			}
			uber, err := parseMeasurement(record, 3, 5)
			if err != nil {
				return r, fmt.Errorf("%s %s uber: %w", dir, name, err)
			}
			cold, err := parseMeasurement(record, 7, 9)
			if err != nil {
				return r, fmt.Errorf("%s %s pure-cold: %w", dir, name, err)
			}
			warm, err := parseMeasurement(record, 11, 13)
			if err != nil {
				return r, fmt.Errorf("%s %s pure-warm: %w", dir, name, err)
			}
			m := metric{pure: pure, uber: uber, cold: cold, warm: warm}
			if r.metrics[name] == nil {
				r.metrics[name] = map[string]metric{}
			}
			r.metrics[name][metricName] = m
			if metricName == "sec/op" {
				r.order = append(r.order, name)
			}
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
	if len(r.order) != len(scenarios) {
		return r, fmt.Errorf("%s: got %d benchmark scenarios, catalog has %d", dir, len(r.order), len(scenarios))
	}
	for _, scenario := range scenarios {
		metrics := r.metrics[scenario.name]
		for _, name := range []string{"sec/op", "B/op", "allocs/op"} {
			if _, ok := metrics[name]; !ok {
				return r, fmt.Errorf("%s: benchmark scenario %q metric %s disappeared", dir, scenario.name, name)
			}
		}
	}
	for _, name := range r.order {
		if _, ok := scenarioByName[name]; !ok {
			return r, fmt.Errorf("%s: undocumented benchmark scenario %q", dir, name)
		}
	}
	if err := loadMemoryResults(path, &r); err != nil {
		return r, err
	}
	return r, nil
}

func loadMemoryResults(path string, r *result) error {
	f, err := os.Open(filepath.Join(path, "memory.tsv"))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	reader := csv.NewReader(f)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("%s/memory.tsv: empty", r.dir)
	}
	header := make(map[string]int, len(records[0]))
	for i, field := range records[0] {
		header[field] = i
	}
	for _, field := range []string{"impl", "workload", "peak_rss_kb", "heap_alloc_kb", "checksum"} {
		if _, ok := header[field]; !ok {
			return fmt.Errorf("%s/memory.tsv: missing %s", r.dir, field)
		}
	}
	for _, record := range records[1:] {
		peakRSSKB, err := parseTSVFloat(record, header["peak_rss_kb"])
		if err != nil {
			return fmt.Errorf("%s/memory.tsv peak_rss_kb: %w", r.dir, err)
		}
		heapAllocKB, err := parseTSVFloat(record, header["heap_alloc_kb"])
		if err != nil {
			return fmt.Errorf("%s/memory.tsv heap_alloc_kb: %w", r.dir, err)
		}
		workload := record[header["workload"]]
		m := r.memory[workload]
		value := memoryMeasurement{
			peakRSSKB:   peakRSSKB,
			heapAllocKB: heapAllocKB,
			checksum:    record[header["checksum"]],
			present:     true,
		}
		switch record[header["impl"]] {
		case "pure":
			m.pure = value
		case "uber":
			m.uber = value
		default:
			return fmt.Errorf("%s/memory.tsv: unknown implementation %q", r.dir, record[header["impl"]])
		}
		r.memory[workload] = m
	}
	for _, scenario := range memoryScenarios {
		m := r.memory[scenario.name]
		if !m.pure.present || !m.uber.present {
			return fmt.Errorf("%s: process-memory scenario %q disappeared", r.dir, scenario.name)
		}
		if m.pure.checksum != m.uber.checksum {
			return fmt.Errorf("%s: process-memory scenario %q checksums differ", r.dir, scenario.name)
		}
	}
	if len(r.memory) != len(memoryScenarios) {
		return fmt.Errorf("%s: got %d process-memory scenarios, catalog has %d", r.dir, len(r.memory), len(memoryScenarios))
	}
	return nil
}

func parseTSVFloat(record []string, index int) (float64, error) {
	if index >= len(record) {
		return 0, fmt.Errorf("short row")
	}
	return strconv.ParseFloat(record[index], 64)
}

func parseMeasurement(record []string, valueIndex, deltaIndex int) (measurement, error) {
	if valueIndex >= len(record) || record[valueIndex] == "" {
		return measurement{}, nil
	}
	value, err := strconv.ParseFloat(record[valueIndex], 64)
	if err != nil {
		return measurement{}, err
	}
	m := measurement{value: value, present: true}
	if valueIndex+1 < len(record) {
		m.ci = record[valueIndex+1]
	}
	if deltaIndex >= 0 && deltaIndex < len(record) {
		m.delta = record[deltaIndex]
		if deltaIndex+1 < len(record) {
			m.p = record[deltaIndex+1]
		}
	}
	return m, nil
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
