// Command layoutinventory classifies every Go and C file in the repository
// root — the flat h3 package that deliberately mixes the public API layer
// with the mechanically ported C layer (DR-001/DR-008) — into its
// architectural layer, so that the layering stays discoverable without a
// package split.
//
// Usage (from the repository root):
//
//	go run ./tools/layoutinventory > docs/file-layer-inventory.csv
//	go run ./tools/layoutinventory -verify   # fail on unclassifiable files
//
// Output: CSV to stdout with columns
// file,layer,package,build_tags,attributions,h3_c_api_refs,exported_decls
// and a per-layer summary to stderr.
//
// Layers (first matching rule wins; rules documented in
// docs/repository-layout-review.md):
//
//	parity-c-shim        h3lib_*.c — C shims compiling pristine upstream sources
//	parity-test          *_parity_test.go — Go-vs-C comparisons (cgo && c2go)
//	parity-cgo-bridge    *_cgo.go — cgo declarations for the parity harness
//	public-api-test      tests of the exported API (black-box example included)
//	ported-upstream-test test<Upstream>_test.go / upstream_*_test.go — ports of
//	                     upstream C test programs, fixtures, and fuzz harnesses
//	ported-internal-test <cmodule>-prefixed white-box tests of ported code
//	ported-public-types  ported-convention files that also declare exported types
//	ported-impl          <cfile>_<name>.go — the mechanically ported C layer
//	public-api           topical files (cell.go, traversal.go, …) — the public API
//
// Classification uses file content (package clause, build tags, attribution
// comments, exported declarations) plus the filename shape; it never needs
// the upstream tree, so it runs without testref/.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// cModules are the C translation-unit prefixes used by ported files
// (kept in sync with the lint-exclusion path regex in .golangci.yml).
// Both historical casings of h3Index and vertexGraph appear in filenames.
var cModules = map[string]bool{
	"adder": true, "area": true, "cellsToMultiPoly": true,
	"algos": true, "baseCells": true, "bbox": true, "coordijk": true,
	"directedEdge": true, "faceijk": true, "h3api": true, "h3Index": true,
	"h3index": true, "iterators": true, "latLng": true, "linkedGeo": true,
	"localij": true, "mathExtensions": true, "polyfill": true,
	"polygon": true, "utility": true, "vec2d": true, "vec3d": true,
	"vertex": true, "vertexGraph": true, "vertexgraph": true,
}

var (
	upstreamTestRe = regexp.MustCompile(`^(test[A-Z]|upstream_)`)
	packageRe      = regexp.MustCompile(`(?m)^package\s+(\w+)`)
	buildTagRe     = regexp.MustCompile(`(?m)^//go:build\s+(.+)$`)
)

type row struct {
	file          string
	layer         string
	pkg           string
	buildTags     string
	attributions  int
	apiRefs       int
	exportedDecls string // integer for non-test Go files, "-" otherwise
}

func main() {
	repoRoot := flag.String("repo", ".", "repository root")
	verify := flag.Bool("verify", false, "exit non-zero if any file cannot be classified")
	flag.Parse()

	entries, err := os.ReadDir(*repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "layoutinventory: %v\n", err)
		os.Exit(1)
	}

	var rows []row
	var unclassified []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			continue
		}
		isGo := strings.HasSuffix(name, ".go")
		isC := strings.HasSuffix(name, ".c") || strings.HasSuffix(name, ".h")
		if !isGo && !isC {
			continue
		}
		data, err := os.ReadFile(filepath.Join(*repoRoot, name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "layoutinventory: %v\n", err)
			os.Exit(1)
		}
		r := classify(name, string(data), *repoRoot)
		if r.layer == "unclassified" {
			unclassified = append(unclassified, name)
		}
		rows = append(rows, r)
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].file < rows[j].file })

	fmt.Println("file,layer,package,build_tags,attributions,h3_c_api_refs,exported_decls")
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.layer]++
		fmt.Printf("%s,%s,%s,%q,%d,%d,%s\n",
			r.file, r.layer, r.pkg, r.buildTags, r.attributions, r.apiRefs, r.exportedDecls)
	}

	var layers []string
	for l := range counts {
		layers = append(layers, l)
	}
	sort.Strings(layers)
	fmt.Fprintf(os.Stderr, "layoutinventory: %d files\n", len(rows))
	for _, l := range layers {
		fmt.Fprintf(os.Stderr, "  %-22s %d\n", l, counts[l])
	}

	if len(unclassified) > 0 {
		fmt.Fprintf(os.Stderr, "layoutinventory: %d file(s) do not match any layer rule:\n", len(unclassified))
		for _, f := range unclassified {
			fmt.Fprintf(os.Stderr, "  %s\n", f)
		}
		if *verify {
			os.Exit(1)
		}
	}
}

func classify(name, content, repoRoot string) row {
	r := row{
		file:          name,
		pkg:           "-",
		buildTags:     "-",
		attributions:  strings.Count(content, "Ported from H3 C:"),
		apiRefs:       strings.Count(content, "H3 C API:"),
		exportedDecls: "-",
	}
	if m := buildTagRe.FindStringSubmatch(content); m != nil {
		r.buildTags = strings.TrimSpace(m[1])
	}

	if !strings.HasSuffix(name, ".go") {
		if strings.HasPrefix(name, "h3lib_") {
			r.layer = "parity-c-shim"
		} else {
			r.layer = "unclassified"
		}
		return r
	}

	if m := packageRe.FindStringSubmatch(content); m != nil {
		r.pkg = m[1]
	}

	base := strings.TrimSuffix(name, ".go")
	isTest := strings.HasSuffix(base, "_test")
	prefix, _, hasUnderscore := strings.Cut(name, "_")

	switch {
	case strings.HasSuffix(base, "_parity_test"):
		r.layer = "parity-test"
	case isTest && upstreamTestRe.MatchString(name):
		r.layer = "ported-upstream-test"
	case isTest && hasUnderscore && cModules[prefix]:
		r.layer = "ported-internal-test"
	case isTest:
		r.layer = "public-api-test"
	case strings.HasSuffix(base, "_cgo"):
		r.layer = "parity-cgo-bridge"
	case hasUnderscore && cModules[prefix]:
		n := countExportedDecls(filepath.Join(repoRoot, name))
		r.exportedDecls = fmt.Sprint(n)
		if n > 0 {
			r.layer = "ported-public-types"
		} else {
			r.layer = "ported-impl"
		}
	case !hasUnderscore:
		r.exportedDecls = fmt.Sprint(countExportedDecls(filepath.Join(repoRoot, name)))
		r.layer = "public-api"
	default:
		r.layer = "unclassified"
	}
	return r
}

// countExportedDecls counts exported top-level identifiers (types, funcs,
// consts, vars) declared in the file.
func countExportedDecls(path string) int {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return 0
	}
	n := 0
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			if decl.Recv == nil && decl.Name.IsExported() {
				n++
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name.IsExported() {
						n++
					}
				case *ast.ValueSpec:
					for _, id := range s.Names {
						if id.IsExported() {
							n++
						}
					}
				}
			}
		}
	}
	return n
}
