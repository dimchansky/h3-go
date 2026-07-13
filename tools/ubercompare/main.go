// Command ubercompare maintains the uber/h3-go comparison matrix.
//
// The curated matrix lives in docs/comparison-uber-h3-go.csv: one row per
// public C function of the pinned H3 release, mapping it to this library's
// API and to the official cgo binding's API. This tool keeps that matrix
// consistent with the repository's authoritative inventories — offline, so
// it can gate CI without network or testref sources:
//
//   - every row must correspond to exactly one public C function from
//     docs/c-api-inventory.csv, and every public C function must have a
//     row (drift gate for H3 version bumps);
//   - every Go symbol named in the matrix's this_api column must exist in
//     docs/api-surface.txt (drift gate for API changes in this library);
//   - status and migration values must come from the fixed vocabularies;
//   - the generated tables in docs/comparison-uber-h3-go.md must match
//     the matrix (drift gate for hand-edited tables).
//
// The binding's side of the matrix cannot be verified offline from the
// root module (it would need the uber/h3-go dependency); that check lives
// in interop/uberbench (TestMappingSymbolsExist) and runs with the
// benchmark suite.
//
// Usage:
//
//	go run ./tools/ubercompare            # print the generated Markdown tables
//	go run ./tools/ubercompare -write     # rewrite the generated section of docs/comparison-uber-h3-go.md
//	go run ./tools/ubercompare -verify    # exit 1 on any inconsistency (CI gate)
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	matrixCSV    = "docs/comparison-uber-h3-go.csv"
	inventoryCSV = "docs/c-api-inventory.csv"
	surfaceFile  = "docs/api-surface.txt"
	targetDoc    = "docs/comparison-uber-h3-go.md"

	beginMarker = "<!-- BEGIN GENERATED: ubercompare (edit docs/comparison-uber-h3-go.csv and run `make gen-ubercompare`) -->"
	endMarker   = "<!-- END GENERATED: ubercompare -->"
)

var (
	uberStatuses = map[string]bool{"available": true, "different-shape": true, "absorbed": true, "missing": true}
	migrations   = map[string]bool{"mechanical": true, "adaptation": true, "n/a": true}

	// categoryOrder controls table grouping in the generated Markdown.
	categoryOrder = []string{"indexing", "inspection", "traversal", "hierarchy", "edges", "vertexes", "regions", "measurement", "misc"}
	categoryTitle = map[string]string{
		"indexing":    "Indexing",
		"inspection":  "Index inspection",
		"traversal":   "Grid traversal",
		"hierarchy":   "Hierarchy and compaction",
		"edges":       "Directed edges",
		"vertexes":    "Vertexes",
		"regions":     "Regions (polyfill and multi-polygon)",
		"measurement": "Measurement",
		"misc":        "Constants, conversions, and error description",
	}

	// Additive ergonomics that do not each correspond to their own public C
	// function row still belong to the comparison contract. Keep them locked
	// here so the migration/comparison docs cannot claim them after removal.
	ergonomicSurface = []string{
		"Index",
		"NumIcosahedronFaces",
		"Cell.ImmediateParent",
		"Cell.ImmediateChildren",
		"Cell.AppendImmediateChildren",
		"Cell.GridDiskDistancesGrouped",
		"DirectedEdge.IndexDigit",
		"Vertex.IndexDigit",
	}
)

type row struct {
	CFunc, Category, ThisAPI, UberAPI, UberStatus, Semantics, Allocation, Migration string
}

func main() {
	repo := flag.String("repo", ".", "repository root")
	write := flag.Bool("write", false, "rewrite the generated section of "+targetDoc)
	verify := flag.Bool("verify", false, "verify consistency and exit non-zero on drift")
	flag.Parse()

	if err := run(*repo, *write, *verify); err != nil {
		fmt.Fprintln(os.Stderr, "ubercompare:", err)
		os.Exit(1)
	}
}

func run(repo string, write, verify bool) error {
	rows, err := readMatrix(repo + "/" + matrixCSV)
	if err != nil {
		return err
	}
	public, err := readPublicCFuncs(repo + "/" + inventoryCSV)
	if err != nil {
		return err
	}
	surface, err := readSurface(repo + "/" + surfaceFile)
	if err != nil {
		return err
	}

	var problems []string

	// 1. Row set must equal the public C function set.
	seen := map[string]bool{}
	for _, r := range rows {
		if seen[r.CFunc] {
			problems = append(problems, fmt.Sprintf("duplicate matrix row: %s", r.CFunc))
		}
		seen[r.CFunc] = true
		if !public[r.CFunc] {
			problems = append(problems, fmt.Sprintf("matrix row %q is not a public C function in %s", r.CFunc, inventoryCSV))
		}
	}
	for _, sym := range ergonomicSurface {
		if !symbolInSurface(surface, sym) {
			problems = append(problems, fmt.Sprintf("ergonomic comparison symbol %q not found in %s", sym, surfaceFile))
		}
	}
	for f := range public {
		if !seen[f] {
			problems = append(problems, fmt.Sprintf("public C function %q has no matrix row", f))
		}
	}

	// 2. Vocabularies and shape rules.
	for _, r := range rows {
		if !uberStatuses[r.UberStatus] {
			problems = append(problems, fmt.Sprintf("%s: unknown uber_status %q", r.CFunc, r.UberStatus))
		}
		if !migrations[r.Migration] {
			problems = append(problems, fmt.Sprintf("%s: unknown migration %q", r.CFunc, r.Migration))
		}
		hasUberAPI := !strings.HasPrefix(r.UberAPI, "—")
		switch r.UberStatus {
		case "available", "different-shape":
			if !hasUberAPI {
				problems = append(problems, fmt.Sprintf("%s: status %s but no uber_api", r.CFunc, r.UberStatus))
			}
		case "absorbed", "missing":
			if hasUberAPI {
				problems = append(problems, fmt.Sprintf("%s: status %s but uber_api %q", r.CFunc, r.UberStatus, r.UberAPI))
			}
		}
	}

	// 3. Every this_api symbol must exist in the locked API surface.
	for _, r := range rows {
		for _, sym := range apiSymbols(r.ThisAPI) {
			if !symbolInSurface(surface, sym) {
				problems = append(problems, fmt.Sprintf("%s: this_api symbol %q not found in %s", r.CFunc, sym, surfaceFile))
			}
		}
	}

	generated := render(rows)

	if verify {
		doc, err := os.ReadFile(repo + "/" + targetDoc)
		if err != nil {
			return err
		}
		current, err := generatedSection(string(doc))
		if err != nil {
			problems = append(problems, err.Error())
		} else if strings.TrimSpace(current) != strings.TrimSpace(generated) {
			problems = append(problems, fmt.Sprintf("%s generated section is stale; run `make gen-ubercompare`", targetDoc))
		}
		if len(problems) > 0 {
			for _, p := range problems {
				fmt.Fprintln(os.Stderr, "ubercompare:", p)
			}
			return fmt.Errorf("%d problem(s)", len(problems))
		}
		fmt.Printf("ubercompare: OK (%d public C functions, matrix and doc in sync)\n", len(rows))
		return nil
	}

	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, "ubercompare:", p)
		}
		return fmt.Errorf("%d problem(s)", len(problems))
	}

	if write {
		doc, err := os.ReadFile(repo + "/" + targetDoc)
		if err != nil {
			return err
		}
		updated, err := replaceGeneratedSection(string(doc), generated)
		if err != nil {
			return err
		}
		if err := os.WriteFile(repo+"/"+targetDoc, []byte(updated), 0o644); err != nil {
			return err
		}
		fmt.Printf("ubercompare: %s updated (%d rows)\n", targetDoc, len(rows))
		return nil
	}

	fmt.Print(generated)
	return nil
}

func readMatrix(path string) ([]row, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	rec, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(rec) == 0 {
		return nil, fmt.Errorf("%s: empty", path)
	}
	want := []string{"c_function", "category", "this_api", "uber_api", "uber_status", "semantics", "allocation", "migration"}
	if !slices.Equal(rec[0], want) {
		return nil, fmt.Errorf("%s: header %v, want %v", path, rec[0], want)
	}
	rows := make([]row, 0, len(rec)-1)
	for _, r := range rec[1:] {
		rows = append(rows, row{r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7]})
	}
	return rows, nil
}

func readPublicCFuncs(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	rec, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	out := map[string]bool{}
	for i, r := range rec {
		if i == 0 || len(r) < 6 {
			continue
		}
		if r[5] == "true" { // is_c_public
			out[r[0]] = true
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s: no public C functions found", path)
	}
	return out, nil
}

func readSurface(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := map[string]bool{}
	for line := range strings.SplitSeq(string(data), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out[line] = true
		}
	}
	return out, nil
}

// apiSymbols extracts checkable Go symbols from a this_api cell: entries
// are "; "-separated, absorbed entries start with "—", and parenthetical
// annotations are ignored.
func apiSymbols(cell string) []string {
	var out []string
	for entry := range strings.SplitSeq(cell, ";") {
		entry = strings.TrimSpace(entry)
		if entry == "" || strings.HasPrefix(entry, "—") {
			continue
		}
		if i := strings.Index(entry, " ("); i >= 0 {
			entry = entry[:i]
		}
		out = append(out, entry)
	}
	return out
}

// symbolInSurface matches "Name" or "Type.Method" against the api-surface
// line formats ("func Name", "const Name", "var Name", "type Name",
// "method (Type) Method", "method (*Type) Method").
func symbolInSurface(surface map[string]bool, sym string) bool {
	if typ, method, ok := strings.Cut(sym, "."); ok {
		return surface["method ("+typ+") "+method] || surface["method (*"+typ+") "+method]
	}
	return surface["func "+sym] || surface["const "+sym] || surface["var "+sym] || surface["type "+sym]
}

func render(rows []row) string {
	byCat := map[string][]row{}
	statusCounts := map[string]int{}
	for _, r := range rows {
		byCat[r.Category] = append(byCat[r.Category], r)
		statusCounts[r.UberStatus]++
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", beginMarker)
	fmt.Fprintf(&b, "Legend — **uber/h3-go status**: `available` (same operation, directly\n")
	fmt.Fprintf(&b, "comparable shape), `different-shape` (same operation, different Go\n")
	fmt.Fprintf(&b, "types/structure), `absorbed` (not exposed because Go semantics make it\n")
	fmt.Fprintf(&b, "unnecessary — sizing helpers, memory destructors, unit converters),\n")
	fmt.Fprintf(&b, "`missing` (no equivalent). **Migration**: `mechanical` (rename/re-type at\n")
	fmt.Fprintf(&b, "the call site), `adaptation` (surrounding code must change shape),\n")
	fmt.Fprintf(&b, "`n/a` (nothing to migrate). A long dash (—) marks an intentionally\n")
	fmt.Fprintf(&b, "absent API.\n")
	fmt.Fprintf(&b, "\n**Status totals for uber/h3-go v4.4.1:** %d `available`, %d\n", statusCounts["available"], statusCounts["different-shape"])
	fmt.Fprintf(&b, "`different-shape`, %d `absorbed`, %d `missing` — %d public H3 C\n", statusCounts["absorbed"], statusCounts["missing"], len(rows))
	fmt.Fprintf(&b, "functions accounted for.\n")
	for _, cat := range categoryOrder {
		rs := byCat[cat]
		if len(rs) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n### %s\n\n", categoryTitle[cat])
		fmt.Fprintf(&b, "| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |\n")
		fmt.Fprintf(&b, "|---|---|---|---|---|---|---|\n")
		for _, r := range rs {
			fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %s | %s | %s |\n",
				r.CFunc, codeify(r.ThisAPI), codeify(r.UberAPI), r.UberStatus, r.Semantics, r.Allocation, r.Migration)
		}
	}
	fmt.Fprintf(&b, "\n%s\n", endMarker)
	return b.String()
}

// codeify wraps each API entry in backticks, leaving absorbed/missing
// dashes and parenthetical annotations readable.
func codeify(cell string) string {
	if cell == "" {
		return ""
	}
	parts := strings.Split(cell, ";")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "—") {
			parts[i] = p
			continue
		}
		if j := strings.Index(p, " ("); j >= 0 {
			parts[i] = "`" + p[:j] + "`" + " " + p[j+1:]
		} else {
			parts[i] = "`" + p + "`"
		}
	}
	return strings.Join(parts, "; ")
}

func generatedSection(doc string) (string, error) {
	i := strings.Index(doc, beginMarker)
	j := strings.Index(doc, endMarker)
	if i < 0 || j < 0 || j < i {
		return "", fmt.Errorf("%s: generated-section markers not found", targetDoc)
	}
	return doc[i : j+len(endMarker)], nil
}

func replaceGeneratedSection(doc, generated string) (string, error) {
	i := strings.Index(doc, beginMarker)
	j := strings.Index(doc, endMarker)
	if i < 0 || j < 0 || j < i {
		return "", fmt.Errorf("%s: generated-section markers not found (add %q and %q)", targetDoc, beginMarker, endMarker)
	}
	return doc[:i] + strings.TrimSpace(generated) + doc[j+len(endMarker):], nil
}
