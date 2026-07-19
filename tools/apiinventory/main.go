// Command apiinventory maps the H3 C public API (h3api.h.in) to the Go port
// in the repository root, using the "// Ported from H3 C: <file>::<name>"
// attribution comments that every ported declaration carries.
//
// Usage (from the repository root, after `make -C testref h3-source` has
// populated testref/):
//
//	go run ./tools/apiinventory > docs/c-api-inventory.csv
//	go run ./tools/apiinventory -h3ver 4.4.0        # future upstream versions
//
// Output: CSV to stdout with columns
// c_function,c_signature,go_file,go_func,go_signature,is_c_public,notes
// and a human-readable summary (counts, missing functions) to stderr.
//
// The tool is the mechanical backbone of the upstream-synchronization
// workflow described in docs/public-api-architecture.md: after vendoring a
// new upstream version under testref/, rerun it to detect added, removed,
// or still-unported public C functions.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type cDecl struct {
	name       string
	signature  string
	deprecated bool
}

type goDecl struct {
	attrFile string // e.g. "algos.c" or "h3Index.h"
	attrName string // e.g. "gridDisk" or "_gridRingInternal"
	goFile   string // relative file name
	goName   string // Go identifier
	goSig    string // full declaration line(s), one line
	kind     string // func, var, const, type, block, other
}

var wsRe = regexp.MustCompile(`\s+`)

func norm(s string) string { return strings.TrimSpace(wsRe.ReplaceAllString(s, " ")) }

func parseHeader(headerPath string) []cDecl {
	data, err := os.ReadFile(headerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read C header %s: %v\n(run `make -C testref h3-source` or pass -h3ver/-header)\n", headerPath, err)
		os.Exit(1)
	}
	text := string(data)
	// DECLSPEC <ret> H3_EXPORT(<name>)(<params>); possibly spanning lines.
	re := regexp.MustCompile(`(?s)DECLSPEC\s+([^;{}]+?)\s*H3_EXPORT\((\w+)\)\s*\(([^;]*?)\)\s*;`)
	matches := re.FindAllStringSubmatchIndex(text, -1)
	var out []cDecl
	for _, m := range matches {
		ret := norm(text[m[2]:m[3]])
		name := text[m[4]:m[5]]
		params := norm(text[m[6]:m[7]])
		sig := fmt.Sprintf("%s %s(%s)", ret, name, params)
		// Look backwards for the immediately preceding doc comment block to
		// detect deprecation markers.
		deprecated := false
		pre := text[:m[0]]
		if idx := strings.LastIndex(pre, "/**"); idx >= 0 {
			if end := strings.Index(pre[idx:], "*/"); end >= 0 {
				block := pre[idx : idx+end]
				between := pre[idx+end+2:]
				if strings.TrimSpace(stripLineComments(between)) == "" {
					deprecated = strings.Contains(strings.ToLower(block), "deprecated")
				}
			}
		}
		out = append(out, cDecl{name: name, signature: sig, deprecated: deprecated})
	}
	return out
}

func stripLineComments(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "*") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

var (
	attrRe = regexp.MustCompile(`//\s*Ported from H3 C:\s*([\w./-]+)::(\S+)`)
	// Some attributions wrap the name: "latLng.c::H3_EXPORT(edgeLengthM)".
	exportWrapRe = regexp.MustCompile(`^H3_EXPORT\((\w+)\)$`)
	funcRe       = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)`)
	varRe        = regexp.MustCompile(`^(?:var|const)\s+(\w+)`)
	typeRe       = regexp.MustCompile(`^type\s+(\w+)`)
	blockRe      = regexp.MustCompile(`^(?:var|const)\s*\($`)
)

func unwrapName(s string) string {
	s = strings.TrimRight(s, ".") // tolerate godot-added trailing periods
	if m := exportWrapRe.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return s
}

func parseGoFiles(repoRoot string) []goDecl {
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read repo root %s: %v\n", repoRoot, err)
		os.Exit(1)
	}
	var out []goDecl
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", name, err)
			os.Exit(1)
		}
		lines := strings.Split(string(data), "\n")
		var pending [][2]string
		inBlockComment := false
		for i := range lines {
			line := strings.TrimSpace(lines[i])
			if inBlockComment {
				if m := attrRe.FindStringSubmatch(line); m != nil {
					pending = append(pending, [2]string{m[1], unwrapName(m[2])})
				}
				if strings.Contains(line, "*/") {
					inBlockComment = false
				}
				continue
			}
			if strings.HasPrefix(line, "/*") {
				if m := attrRe.FindStringSubmatch(line); m != nil {
					pending = append(pending, [2]string{m[1], unwrapName(m[2])})
				}
				if !strings.Contains(line, "*/") {
					inBlockComment = true
				}
				continue
			}
			if strings.HasPrefix(line, "//") {
				if m := attrRe.FindStringSubmatch(line); m != nil {
					pending = append(pending, [2]string{m[1], unwrapName(m[2])})
				}
				continue
			}
			if line == "" {
				continue
			}
			if len(pending) > 0 {
				kind, ident, sig := classify(lines, i)
				for _, p := range pending {
					out = append(out, goDecl{
						attrFile: p[0], attrName: p[1],
						goFile: name, goName: ident, goSig: sig, kind: kind,
					})
				}
				pending = nil
			}
		}
	}
	return out
}

// classify identifies the declaration starting at lines[i] and returns its
// kind, identifier, and a one-line signature (for funcs: up to opening brace).
func classify(lines []string, i int) (kind, ident, sig string) {
	line := strings.TrimSpace(lines[i])
	if m := funcRe.FindStringSubmatch(line); m != nil {
		full := line
		for j := i; !strings.Contains(full, "{") && j+1 < len(lines) && j < i+10; j++ {
			full += " " + strings.TrimSpace(lines[j+1])
		}
		if idx := strings.Index(full, "{"); idx >= 0 {
			full = full[:idx]
		}
		return "func", m[1], norm(full)
	}
	if m := typeRe.FindStringSubmatch(line); m != nil {
		return "type", m[1], norm(line)
	}
	if m := varRe.FindStringSubmatch(line); m != nil {
		k := "var"
		if strings.HasPrefix(line, "const") {
			k = "const"
		}
		s := line
		if idx := strings.Index(s, "="); idx >= 0 {
			s = s[:idx]
		}
		return k, m[1], norm(s)
	}
	if blockRe.MatchString(line) {
		return "block", "", norm(line)
	}
	return "other", "", norm(line)
}

// omissions lists C public functions that intentionally have no dedicated
// public Go API, with the reason. The -verify mode accepts these instead of
// requiring an "H3 C API:" doc reference.
var omissions = map[string]string{
	"describeH3Error":           "surfaced as the Error() text of the sentinel errors (errors.go)",
	"degsToRads":                "replaced by the Angle type: Deg(x).Rad()",
	"radsToDegs":                "replaced by the Angle type: Rad(x).Deg()",
	"destroyLinkedMultiPolygon": "meaningless under garbage collection; CellsToMultiPolygon returns slices",
}

var apiLineRe = regexp.MustCompile(`(?m)//\s*H3 C API:\s*(.+)$`)
var wordRe = regexp.MustCompile(`[A-Za-z_]\w*`)

// collectAPIRefs gathers every identifier mentioned on an "H3 C API:" doc
// line in the non-test root Go files.
func collectAPIRefs(repoRoot string) map[string]bool {
	refs := map[string]bool{}
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return refs
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, name))
		if err != nil {
			continue
		}
		for _, m := range apiLineRe.FindAllStringSubmatch(string(data), -1) {
			for _, w := range wordRe.FindAllString(m[1], -1) {
				refs[w] = true
			}
		}
	}
	return refs
}

func main() {
	repoRoot := flag.String("repo", ".", "repository root containing the ported Go files")
	h3ver := flag.String("h3ver", "4.5.0", "upstream H3 version vendored under testref/")
	header := flag.String("header", "", "explicit path to h3api.h.in (overrides -h3ver)")
	verify := flag.Bool("verify", false, "verify completeness: every C public function must be ported AND either referenced by an 'H3 C API:' doc line or listed in the omissions table; exit 1 otherwise")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "apiinventory maps the H3 C public API to the Go port (CSV to stdout) and verifies completeness with -verify.")
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/apiinventory [flags] > docs/c-api-inventory.csv")
		flag.PrintDefaults()
	}
	flag.Parse()

	headerPath := *header
	if headerPath == "" {
		headerPath = filepath.Join(*repoRoot, "testref", "h3-"+*h3ver, "src", "h3lib", "include", "h3api.h.in")
	}

	cDecls := parseHeader(headerPath)
	goDecls := parseGoFiles(*repoRoot)

	if *verify {
		refs := collectAPIRefs(*repoRoot)
		ported := map[string]bool{}
		for _, g := range goDecls {
			ported[g.attrName] = true
		}
		failed := 0
		for _, c := range cDecls {
			if !ported[c.name] {
				fmt.Fprintf(os.Stderr, "VERIFY FAIL: C public function %s has no Go port (missing 'Ported from H3 C:' attribution)\n", c.name)
				failed++
			}
			if !refs[c.name] && omissions[c.name] == "" {
				fmt.Fprintf(os.Stderr, "VERIFY FAIL: C public function %s has neither an 'H3 C API:' doc reference nor an omissions entry\n", c.name)
				failed++
			}
		}
		if failed > 0 {
			fmt.Fprintf(os.Stderr, "verify: %d problems across %d C public functions\n", failed, len(cDecls))
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "verify: OK — all %d C public functions are ported and publicly represented\n", len(cDecls))
		return
	}

	cByName := map[string]cDecl{}
	for _, c := range cDecls {
		if prev, dup := cByName[c.name]; dup {
			fmt.Fprintf(os.Stderr, "WARNING: duplicate C decl %s (%s vs %s)\n", c.name, prev.signature, c.signature)
		}
		cByName[c.name] = c
	}

	goByAttr := map[string][]goDecl{}
	for _, g := range goDecls {
		goByAttr[g.attrName] = append(goByAttr[g.attrName], g)
	}

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	_ = w.Write([]string{"c_function", "c_signature", "go_file", "go_func", "go_signature", "is_c_public", "notes"})

	// 1) Every C public function, matched to Go ports.
	for _, c := range cDecls {
		notes := ""
		if c.deprecated {
			notes = "deprecated in header"
		}
		ports := goByAttr[c.name]
		if len(ports) == 0 {
			_ = w.Write([]string{c.name, c.signature, "", "", "", "true", strings.TrimSpace(notes + " MISSING")})
			continue
		}
		for _, g := range ports {
			n := notes
			if g.kind != "func" {
				n = strings.TrimSpace(n + " go decl kind: " + g.kind)
			}
			n = strings.TrimSpace(n + " attr: " + g.attrFile)
			_ = w.Write([]string{c.name, c.signature, g.goFile, g.goName, g.goSig, "true", n})
		}
	}

	// 2) Go decls attributed to C names NOT in the public header (internal helpers).
	var internals []goDecl
	for _, g := range goDecls {
		if _, ok := cByName[g.attrName]; !ok {
			internals = append(internals, g)
		}
	}
	sort.Slice(internals, func(a, b int) bool {
		if internals[a].attrFile != internals[b].attrFile {
			return internals[a].attrFile < internals[b].attrFile
		}
		return internals[a].attrName < internals[b].attrName
	})
	for _, g := range internals {
		_ = w.Write([]string{g.attrName, "(C-internal: " + g.attrFile + ")", g.goFile, g.goName, g.goSig, "false", "kind: " + g.kind})
	}

	// Summary to stderr.
	ported, missing := 0, 0
	var missingNames []string
	for _, c := range cDecls {
		if len(goByAttr[c.name]) > 0 {
			ported++
		} else {
			missing++
			missingNames = append(missingNames, c.name)
		}
	}
	fmt.Fprintf(os.Stderr, "C public functions: %d; ported: %d; missing: %d\n", len(cDecls), ported, missing)
	if missing > 0 {
		fmt.Fprintf(os.Stderr, "Missing: %s\n", strings.Join(missingNames, ", "))
	}
	fmt.Fprintf(os.Stderr, "Go attributions total: %d (public-matched + internal = %d + %d)\n",
		len(goDecls), len(goDecls)-len(internals), len(internals))
}
