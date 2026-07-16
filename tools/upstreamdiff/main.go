// Command upstreamdiff compares two upstream H3 C source trees at the
// symbol level and maps every changed C symbol to this repository's Go port
// via the "// Ported from H3 C: <file>::<name>" attribution comments.
//
// It exists because a file-level diff ("h3Index.c changed") is useless for
// reviewing an upgrade of a function-by-function port, and an API-inventory
// check alone only proves that public functions exist — not that their
// implementations match the new upstream version.
//
//	go run ./tools/upstreamdiff -from testref/h3-4.3.0 -to testref/h3-4.4.0
//	go run ./tools/upstreamdiff -from ... -to ... -strict   # exit 1 on unmapped changes
//	go run ./tools/upstreamdiff -from ... -to ... -comments # + leading-comment drift
//
// Output: a Markdown report (stdout) with three sections — library symbol
// changes (functions, tables, macros, types) each mapped to its Go file,
// public-header (h3api.h.in) symbol changes, and upstream test-file changes
// cross-checked against docs/upstream-test-inventory.csv. Symbols whose Go mapping is
// missing or ambiguous are flagged for human review; the tool never edits
// Go code.
//
// By default only symbol bodies are compared (leading comments are
// invisible to the diff; comments inside bodies count as body changes).
// With -comments, the leading documentation comment of every extracted
// symbol is compared as an independent dimension: the symbol table gains a
// Change column (body / comment-only / body+comment) and a derived
// "Leading-comment-only changes" section lists symbols whose body is
// unchanged. This exists because public GoDoc mirrors upstream contract
// comments, so comment drift is documentation drift (see the sync workflow
// in CONTRIBUTING.md). -strict is unaffected by -comments: it still exits 1
// only for changed symbols that lack a Go mapping.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// C symbol extraction
// ---------------------------------------------------------------------------

type symbol struct {
	kind    string // func, table, macro, type
	name    string
	file    string // base name, e.g. "h3Index.c"
	body    string // normalized source text of the definition
	comment string // normalized leading comment block ("" if none)
}

var (
	// H3_EXPORT(name) or plain identifier immediately before '('.
	funcNameRe = regexp.MustCompile(`(?:H3_EXPORT\(\s*(\w+)\s*\)|\b(\w+))\s*\($`)
	defineRe   = regexp.MustCompile(`^#define\s+(\w+)`)
	typedefRe  = regexp.MustCompile(`^\}\s*(\w+)\s*;`)
	// static const Foo bar[...] = {   /  const int foo = ...
	tableRe = regexp.MustCompile(`^(?:static\s+)?const\s+[\w\s\*]+?\b(\w+)\s*(?:\[[^\]]*\])+\s*=`)
	constRe = regexp.MustCompile(`^(?:static\s+)?const\s+[\w\s\*]+?\b(\w+)\s*=`)
)

// extractSymbols parses one C file into top-level symbols. It relies on the
// clang-format layout of the H3 sources (definitions start at column 0,
// bodies use balanced braces) rather than a full C parser.
func extractSymbols(path string) ([]symbol, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	base := filepath.Base(path)
	lines := strings.Split(string(data), "\n")
	var out []symbol

	i := 0
	inComment := false
	var pending []string // most recent contiguous leading comment block
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Track /* ... */ block comments; their continuation lines may start
		// at column 0 without a leading '*' (free-form prose in H3 sources).
		// The block accumulates as the pending leading comment of whatever
		// symbol immediately follows it.
		if inComment {
			pending = append(pending, trimmed)
			if strings.Contains(trimmed, "*/") {
				inComment = false
			}
			i++
			continue
		}
		if strings.HasPrefix(trimmed, "/*") && !strings.Contains(trimmed, "*/") {
			inComment = true
			pending = []string{trimmed}
			i++
			continue
		}

		// #define — single or backslash-continued.
		if m := defineRe.FindStringSubmatch(trimmed); m != nil {
			start := i
			for strings.HasSuffix(strings.TrimSpace(lines[i]), `\`) && i+1 < len(lines) {
				i++
			}
			out = append(out, symbol{kind: "macro", name: m[1], file: base,
				body:    norm(strings.Join(lines[start:i+1], "\n")),
				comment: norm(strings.Join(pending, "\n"))})
			pending = nil
			i++
			continue
		}

		// Single-line comments accumulate as pending; a blank line, other
		// preprocessor directive, or stray closing brace breaks attachment.
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") ||
			strings.HasPrefix(trimmed, "*") {
			pending = append(pending, trimmed)
			i++
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "}") {
			pending = nil
			i++
			continue
		}

		// Top-level definitions only (column 0, identifier-ish start).
		c := line[0]
		isIdentStart := c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		if !isIdentStart {
			pending = nil
			i++
			continue
		}

		// Collect the full statement/definition: until braces balance AND we
		// hit a ';' (declarations, tables) or the closing '}' of a function.
		start := i
		depth := 0
		opened := false
		end := i
		for j := i; j < len(lines); j++ {
			for _, ch := range lines[j] {
				switch ch {
				case '{':
					depth++
					opened = true
				case '}':
					depth--
				}
			}
			t := strings.TrimSpace(lines[j])
			if opened && depth <= 0 {
				end = j
				break
			}
			if !opened && strings.HasSuffix(t, ";") {
				end = j
				break
			}
			end = j
		}
		block := strings.Join(lines[start:end+1], "\n")

		// Classify + name the block.
		header := block
		if idx := strings.Index(block, "{"); idx >= 0 {
			header = block[:idx]
		}
		headerNorm := norm(header)

		name, kind := "", ""
		switch {
		case strings.HasPrefix(trimmed, "typedef"):
			kind = "type"
			if m := typedefRe.FindStringSubmatch(strings.TrimSpace(lines[end])); m != nil {
				name = m[1]
			} else {
				// typedef X Y;
				fields := strings.Fields(strings.TrimSuffix(norm(block), ";"))
				if len(fields) > 1 {
					name = fields[len(fields)-1]
				}
			}
		case tableRe.MatchString(headerNorm):
			kind, name = "table", tableRe.FindStringSubmatch(headerNorm)[1]
		case opened && strings.Contains(header, "("):
			kind = "func"
			// name = identifier (or H3_EXPORT arg) right before the first
			// top-level '(' of the header.
			h := headerNorm
			if idx := parenIdx(h); idx >= 0 {
				if m := funcNameRe.FindStringSubmatch(h[:idx+1]); m != nil {
					if m[1] != "" {
						name = m[1]
					} else {
						name = m[2]
					}
				}
			}
		case constRe.MatchString(headerNorm):
			kind, name = "table", constRe.FindStringSubmatch(headerNorm)[1]
		}

		if name != "" && kind != "" {
			out = append(out, symbol{kind: kind, name: name, file: base,
				body: norm(block), comment: norm(strings.Join(pending, "\n"))})
		}
		pending = nil
		i = end + 1
	}
	return out, nil
}

// parenIdx returns the index of the '(' that starts the parameter list —
// the first '(' unless preceded by H3_EXPORT, in which case the one after
// the macro's closing paren.
func parenIdx(header string) int {
	if idx := strings.Index(header, "H3_EXPORT("); idx >= 0 {
		if closeIdx := strings.Index(header[idx:], ")"); closeIdx >= 0 {
			if p := strings.Index(header[idx+closeIdx:], "("); p >= 0 {
				return idx + closeIdx + p
			}
		}
	}
	return strings.Index(header, "(")
}

var wsRe = regexp.MustCompile(`\s+`)

func norm(s string) string { return strings.TrimSpace(wsRe.ReplaceAllString(s, " ")) }

// ---------------------------------------------------------------------------
// Symbol diff
// ---------------------------------------------------------------------------

// change is one changed/added/removed symbol with its change dimensions
// tracked independently: bodyChanged and commentChanged can both be true.
type change struct {
	status         string // added, changed, removed
	kind           string
	key            string // file::name
	bodyChanged    bool
	commentChanged bool
}

// diffSymbols compares two scanned trees. Without includeComments the
// result is the historical body-only diff (leading comments ignored). With
// includeComments, symbols whose leading documentation comment changed are
// also reported, and body/comment dimensions are set independently so a
// symbol with both changes is visible on both axes.
func diffSymbols(oldSyms, newSyms map[string]symbol, includeComments bool) []change {
	var out []change
	for k, ns := range newSyms {
		olds, ok := oldSyms[k]
		if !ok {
			out = append(out, change{status: "added", kind: ns.kind, key: k, bodyChanged: true})
			continue
		}
		bodyChanged := olds.body != ns.body
		commentChanged := includeComments && olds.comment != ns.comment
		if bodyChanged || commentChanged {
			out = append(out, change{status: "changed", kind: ns.kind, key: k,
				bodyChanged: bodyChanged, commentChanged: commentChanged})
		}
	}
	for k, olds := range oldSyms {
		if _, ok := newSyms[k]; !ok {
			out = append(out, change{status: "removed", kind: olds.kind, key: k, bodyChanged: true})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].key < out[j].key })
	return out
}

// changeDims renders the change-dimension cell for -comments mode.
func changeDims(r change) string {
	switch {
	case r.status != "changed":
		return "—"
	case r.bodyChanged && r.commentChanged:
		return "body+comment"
	case r.commentChanged:
		return "comment-only"
	default:
		return "body"
	}
}

// ---------------------------------------------------------------------------
// Tree scanning
// ---------------------------------------------------------------------------

func scanTree(root string) (map[string]symbol, error) {
	dirs := []string{
		filepath.Join(root, "src", "h3lib", "lib"),
		filepath.Join(root, "src", "h3lib", "include"),
	}
	syms := map[string]symbol{} // key: file::name
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("%s: %w (run `make -C testref H3_VERSION=<ver> h3-source`)", dir, err)
		}
		for _, e := range entries {
			n := e.Name()
			if e.IsDir() || (!strings.HasSuffix(n, ".c") && !strings.HasSuffix(n, ".h") && !strings.HasSuffix(n, ".h.in")) {
				continue
			}
			if n == "h3api.h" { // generated duplicate of h3api.h.in
				continue
			}
			fileSyms, err := extractSymbols(filepath.Join(dir, n))
			if err != nil {
				return nil, err
			}
			for _, s := range fileSyms {
				syms[s.file+"::"+s.name] = s
			}
		}
	}
	return syms, nil
}

func listTests(root string) (map[string]string, error) {
	dir := filepath.Join(root, "src", "apps", "testapps")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".c") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		out[e.Name()] = string(data)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Attribution mapping (Go side)
// ---------------------------------------------------------------------------

var attrRe = regexp.MustCompile(`//\s*Ported from H3 C:\s*([\w./-]+)::(\S+)`)
var exportWrapRe = regexp.MustCompile(`^H3_EXPORT\((\w+)\)$`)

// goAttributions returns C-name -> go files (a symbol may be referenced by
// several Go declarations, e.g. a constant block).
func goAttributions(repoRoot string) (map[string][]string, error) {
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return nil, err
	}
	m := map[string][]string{}
	add := func(cname, gofile string) {
		if !slices.Contains(m[cname], gofile) {
			m[cname] = append(m[cname], gofile)
		}
	}
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, n))
		if err != nil {
			return nil, err
		}
		for _, mt := range attrRe.FindAllStringSubmatch(string(data), -1) {
			name := strings.TrimRight(mt[2], ".,")
			if w := exportWrapRe.FindStringSubmatch(name); w != nil {
				name = w[1]
			}
			add(mt[1]+"::"+name, n)
			add(name, n) // fallback: name-only lookup
		}
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// Report
// ---------------------------------------------------------------------------

func main() {
	from := flag.String("from", "", "path to the old upstream tree (e.g. testref/h3-4.3.0)")
	to := flag.String("to", "", "path to the new upstream tree (e.g. testref/h3-4.4.0)")
	repo := flag.String("repo", ".", "repository root containing the Go port")
	portedTests := flag.String("ported-tests", "docs/upstream-test-inventory.csv", "reviewed upstream test inventory")
	strict := flag.Bool("strict", false, "exit 1 if any changed lib symbol lacks a Go mapping")
	comments := flag.Bool("comments", false, "also diff each symbol's leading documentation comment (documentation-drift review)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "upstreamdiff diffs two upstream H3 trees at the symbol level and maps changes to the Go port (Markdown to stdout).")
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/upstreamdiff -from <oldtree> -to <newtree> [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *from == "" || *to == "" {
		fmt.Fprintln(os.Stderr, "usage: upstreamdiff -from <oldtree> -to <newtree>")
		os.Exit(2)
	}

	oldSyms, err := scanTree(*from)
	check(err)
	newSyms, err := scanTree(*to)
	check(err)
	attrs, err := goAttributions(*repo)
	check(err)

	rows := diffSymbols(oldSyms, newSyms, *comments)

	fmt.Printf("# Upstream diff: %s -> %s\n\n", filepath.Base(*from), filepath.Base(*to))
	fmt.Printf("Symbols scanned: %d old, %d new. Changed/added/removed: %d.\n\n",
		len(oldSyms), len(newSyms), len(rows))

	unmapped := 0
	fmt.Println("## Library and header symbol changes")
	fmt.Println()
	if *comments {
		fmt.Println("| Status | Kind | C symbol | Change | Go mapping | Review |")
		fmt.Println("|---|---|---|---|---|---|")
	} else {
		fmt.Println("| Status | Kind | C symbol | Go mapping | Review |")
		fmt.Println("|---|---|---|---|---|")
	}
	for _, r := range rows {
		name := r.key[strings.Index(r.key, "::")+2:]
		goFiles := attrs[r.key]
		if len(goFiles) == 0 {
			goFiles = attrs[name]
		}
		mapping, note := "", ""
		switch {
		case len(goFiles) > 0:
			mapping = strings.Join(goFiles, ", ")
			note = "port the diff to the mapped file(s)"
		case r.status == "added":
			mapping = "—"
			note = "NEW: port + attribute + test"
		default:
			mapping = "**UNMAPPED**"
			note = "no attribution found — review by hand"
			unmapped++
		}
		if r.status == "removed" {
			note = "removed upstream — check Go side for retirement"
		}
		if *comments {
			if r.status == "changed" && r.commentChanged {
				note += "; review the leading-comment diff for GoDoc drift"
			}
			fmt.Printf("| %s | %s | `%s` | %s | %s | %s |\n",
				r.status, r.kind, r.key, changeDims(r), mapping, note)
		} else {
			fmt.Printf("| %s | %s | `%s` | %s | %s |\n", r.status, r.kind, r.key, mapping, note)
		}
	}

	if *comments {
		fmt.Println()
		fmt.Println("### Leading-comment-only changes")
		fmt.Println()
		n := 0
		for _, r := range rows {
			if r.status == "changed" && r.commentChanged && !r.bodyChanged {
				fmt.Printf("- `%s`\n", r.key)
				n++
			}
		}
		if n == 0 {
			fmt.Println("(none)")
		}
	}

	// Test files.
	oldTests, err := listTests(*from)
	check(err)
	newTests, err := listTests(*to)
	check(err)
	portedList, _ := os.ReadFile(filepath.Join(*repo, *portedTests))

	fmt.Println()
	fmt.Println("## Upstream test-file changes (src/apps/testapps)")
	fmt.Println()
	fmt.Println("| Status | File | Ported to Go? | Review |")
	fmt.Println("|---|---|---|---|")
	var testNames []string
	for n := range newTests {
		testNames = append(testNames, n)
	}
	sort.Strings(testNames)
	for _, n := range testNames {
		oldBody, existed := oldTests[n]
		status := ""
		switch {
		case !existed:
			status = "added"
		case oldBody != newTests[n]:
			status = "changed"
		default:
			continue
		}
		ported := "no"
		if strings.Contains(string(portedList), n) {
			ported = "yes"
		}
		note := "review the test diff; port new/changed cases"
		if status == "added" && ported == "no" {
			note = "NEW upstream test suite — port it"
		}
		fmt.Printf("| %s | `%s` | %s | %s |\n", status, n, ported, note)
	}
	for n := range oldTests {
		if _, ok := newTests[n]; !ok {
			fmt.Printf("| removed | `%s` | — | removed upstream |\n", n)
		}
	}

	fmt.Println()
	fmt.Printf("Unmapped changed symbols: %d\n", unmapped)
	if *strict && unmapped > 0 {
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
