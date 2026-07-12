// Command unexport performs the one-time mechanical unexport sweep described
// in docs/public-api-architecture.md §7 / Phase 2: every exported identifier
// of the ported C layer that is not part of the intended public surface is
// renamed to its unexported form (case/format change only, so the identifier
// remains recognizably the C name).
//
// Dry run (prints the rename map and any collisions):
//
//	go run ./tools/unexport
//
// Apply:
//
//	go run ./tools/unexport -apply
//
// Renames are word-boundary textual replacements across every root *.go file
// (including tests and cgo parity files), with two guards:
//   - lines containing "Ported from H3 C:" or "H3 C API:" are never touched
//     (they reference C names and drive tools/apiinventory);
//   - occurrences preceded by '.' are never touched (selector expressions,
//     notably cgo's C.H3Index / C.H3Error);
//   - block comments in files that import "C" are never touched (cgo
//     preambles contain real C code).
//
// The tool refuses to apply if any target name collides with an existing
// top-level identifier or another target; known collisions are resolved in
// the special-case table below and documented in the architecture document.
//
// HISTORICAL: the sweep was executed during the public-API build-out and is
// not part of any Makefile target, workflow, or maintenance process. The
// tool is retained only as the reviewable record of that migration
// (referenced by docs/DEVIATIONS.md item 11); you should not need to run
// -apply again.
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
	"unicode"
)

// keep is the intended public surface (Phase 1/3 files plus public geometry
// types); everything else exported gets unexported.
var keep = map[string]bool{
	// types
	"Cell": true, "DirectedEdge": true, "Vertex": true, "Angle": true,
	"LatLng": true, "CellBoundary": true, "GeoLoop": true, "GeoPolygon": true,
	"CoordIJ": true, "ContainmentMode": true,
	// funcs
	"Deg": true, "Rad": true,
	// consts
	"RadPerDeg": true, "DegPerRad": true, "Pi": true, "TwoPi": true, "PiOver2": true,
	"MaxResolution": true, "NumBaseCells": true, "NumRes0Cells": true,
	"NumPentagons": true, "MaxCellBoundaryVerts": true,
	"VersionMajor": true, "VersionMinor": true, "VersionPatch": true,
}

// special resolves names whose mechanical rename would collide or mislead,
// plus C constants that become part of the public API under Go names.
var special = map[string]string{
	"H3Index": "h3Index",
	"H3Error": "h3Error",
	"BBox":    "bbox",
	"Vec2d":   "vec2d",
	"Vec3d":   "vec3d",
	// vertex_constants.go already has an unexported `directions` table.
	"DIRECTIONS": "algosDirections",
	// C declares `typedef struct {...} BaseCellData` AND a table variable
	// `baseCellData[]` (same for PentagonDirectionFaces); the case-lowered
	// type name would collide with the existing table, so the struct types
	// get an Entry suffix.
	"BaseCellData":           "baseCellDataEntry",
	"PentagonDirectionFaces": "pentagonDirectionFacesEntry",
	// plain `ij`/`ki`/`jk` would be shadowed by ubiquitous local variables.
	"IJ": "quadIJ", "KI": "quadKI", "JK": "quadJK",
	// localij__cellToLocalIjk.go has a local int32 `pentagonRotations` in the
	// same function that indexes this table; prefix with the C module name.
	"PENTAGON_ROTATIONS": "localijPentagonRotations",
	// ContainmentMode values are public API (PolygonToCellsExperimental).
	"CONTAINMENT_CENTER":           "ContainmentCenter",
	"CONTAINMENT_FULL":             "ContainmentFull",
	"CONTAINMENT_OVERLAPPING":      "ContainmentOverlapping",
	"CONTAINMENT_OVERLAPPING_BBOX": "ContainmentOverlappingBBox",
	"CONTAINMENT_INVALID":          "ContainmentInvalid",
}

func isScreaming(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
	}
	return true
}

// rename computes the unexported form of an identifier.
func rename(name string) string {
	if n, ok := special[name]; ok {
		return n
	}
	if strings.Contains(name, "_") || isScreaming(name) {
		segs := strings.Split(name, "_")
		var b strings.Builder
		b.WriteString(strings.ToLower(segs[0]))
		for _, s := range segs[1:] {
			if s == "" {
				continue
			}
			switch s {
			case "II", "III": // roman numerals stay uppercase
				b.WriteString(s)
			default:
				lower := strings.ToLower(s)
				r := []rune(lower)
				r[0] = unicode.ToUpper(r[0])
				b.WriteString(string(r))
			}
		}
		return b.String()
	}
	r := []rune(name)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

type decl struct {
	name, kind, file string
}

func main() {
	apply := flag.Bool("apply", false, "apply the renames (default: dry run)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "unexport is the HISTORICAL one-time unexport sweep of the ported C layer; kept as a migration record.")
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/unexport [-apply]   # dry run by default")
		flag.PrintDefaults()
	}
	flag.Parse()

	files, err := filepath.Glob("*.go")
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	var renames []decl
	allTopLevel := map[string]string{} // name -> file (every top-level ident, incl. tests)

	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}
		af, err := parser.ParseFile(fset, f, src, parser.SkipObjectResolution)
		if err != nil {
			panic(fmt.Sprintf("%s: %v", f, err))
		}
		isTest := strings.HasSuffix(f, "_test.go")
		for _, d := range af.Decls {
			for _, name := range declNames(d) {
				if prev, dup := allTopLevel[name]; dup && !isTest {
					_ = prev // duplicate decls across build tags are fine
				}
				allTopLevel[name] = f
				if isTest {
					continue // test helpers are not part of the surface
				}
				if !ast.IsExported(name) || keep[name] || strings.HasPrefix(name, "Err") {
					continue
				}
				renames = append(renames, decl{name: name, kind: declKind(d), file: f})
			}
		}
	}

	sort.Slice(renames, func(i, j int) bool { return renames[i].name < renames[j].name })

	// Build and validate the rename map.
	m := map[string]string{}
	targets := map[string]string{}
	bad := 0
	for _, d := range renames {
		to := rename(d.name)
		if to == d.name {
			fmt.Printf("ERROR: %-30s rename is identity\n", d.name)
			bad++
			continue
		}
		if prev, dup := targets[to]; dup && prev != d.name {
			fmt.Printf("ERROR: %-30s -> %-30s collides with rename of %s\n", d.name, to, prev)
			bad++
		}
		if f, exists := allTopLevel[to]; exists {
			fmt.Printf("ERROR: %-30s -> %-30s collides with existing decl in %s\n", d.name, to, f)
			bad++
		}
		targets[to] = d.name
		m[d.name] = to
	}
	fmt.Printf("%d exported identifiers to rename, %d collisions\n", len(m), bad)
	if !*apply {
		for _, d := range renames {
			fmt.Printf("%-8s %-35s -> %-35s (%s)\n", d.kind, d.name, m[d.name], d.file)
		}
		return
	}
	if bad > 0 {
		fmt.Println("refusing to apply with collisions")
		os.Exit(1)
	}

	// Longest-first application avoids any prefix ambiguity (word boundaries
	// already prevent it, but ordering makes it structurally impossible).
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	res := make([]*regexp.Regexp, len(names))
	for i, n := range names {
		res[i] = regexp.MustCompile(`(^|[^\w.])` + regexp.QuoteMeta(n) + `\b`)
	}

	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}
		text := string(src)
		isCgo := strings.Contains(text, `import "C"`)
		lines := strings.Split(text, "\n")
		inBlockComment := false
		changed := false
		for i, line := range lines {
			if isCgo {
				// Never rewrite cgo preamble C code inside block comments.
				trimmed := strings.TrimSpace(line)
				if inBlockComment {
					if strings.Contains(line, "*/") {
						inBlockComment = false
					}
					continue
				}
				if strings.HasPrefix(trimmed, "/*") && !strings.Contains(trimmed, "*/") {
					inBlockComment = true
					continue
				}
			}
			if strings.Contains(line, "Ported from H3 C:") || strings.Contains(line, "H3 C API:") {
				continue
			}
			orig := line
			for j, re := range res {
				line = re.ReplaceAllString(line, `${1}`+m[names[j]])
			}
			if line != orig {
				lines[i] = line
				changed = true
			}
		}
		if changed {
			if err := os.WriteFile(f, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("applied; run gofmt, build, and the full test suites")
}

func declNames(d ast.Decl) []string {
	switch d := d.(type) {
	case *ast.FuncDecl:
		if d.Recv != nil {
			return nil // methods keep their names; receiver types are renamed
		}
		return []string{d.Name.Name}
	case *ast.GenDecl:
		var out []string
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				out = append(out, s.Name.Name)
			case *ast.ValueSpec:
				for _, n := range s.Names {
					out = append(out, n.Name)
				}
			}
		}
		return out
	}
	return nil
}

func declKind(d ast.Decl) string {
	switch d := d.(type) {
	case *ast.FuncDecl:
		return "func"
	case *ast.GenDecl:
		return d.Tok.String()
	}
	return "?"
}
