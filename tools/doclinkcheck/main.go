// Command doclinkcheck verifies that every [Symbol] / [Type.Method]
// doc-link candidate in the root package's doc comments resolves to a
// symbol declared in that package.
//
// go/doc/comment leaves an unresolvable link candidate as plain text rather
// than emitting a DocLink node, so a checker that only walks parsed DocLink
// nodes cannot detect broken links. This tool instead scans the syntactic
// candidates itself and resolves each against the package's declarations.
//
//	go run ./tools/doclinkcheck            # check the repository root package (make check-docs)
//	go run ./tools/doclinkcheck -dir DIR   # check another package directory
//
// Only candidates that look like package-local exported references are
// checked ([Cell], [Cell.GridDisk]): bracketed text not starting with an
// uppercase letter (ranges like "[0, Len())", GeoJSON's "[lng, lat]") and
// package-qualified references (e.g. [strconv.ErrSyntax], whose qualifier
// is lowercase) are ignored. Test files are skipped. It exits 1 and lists
// each unresolved reference as "file:line: message".
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

var candRe = regexp.MustCompile(`\[([A-Z][A-Za-z0-9_]*(?:\.[A-Z][A-Za-z0-9_]*)?)\]`)

func main() {
	dir := flag.String("dir", ".", "package directory to check")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "doclinkcheck verifies [Symbol] doc links in a package's doc comments.")
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/doclinkcheck [-dir DIR]")
		flag.PrintDefaults()
	}
	flag.Parse()

	problems, err := run(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, p := range problems {
		fmt.Println(p)
	}
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "doclinkcheck: %d unresolved doc link(s)\n", len(problems))
		os.Exit(1)
	}
}

// run checks the package files in dir (non-test .go files, no recursion)
// and returns one "file:line: message" entry per unresolved doc-link
// candidate, sorted by position.
func run(dir string) ([]string, error) {
	files, fset, err := parsePackageFiles(dir)
	if err != nil {
		return nil, err
	}

	declared := declaredNames(files)
	var problems []string
	for _, file := range files {
		for _, cg := range docComments(file) {
			pos := fset.Position(cg.Pos())
			for _, m := range candRe.FindAllStringSubmatch(cg.Text(), -1) {
				name := m[1]
				if !declared[name] {
					problems = append(problems,
						fmt.Sprintf("%s:%d: unresolved doc link [%s]", pos.Filename, pos.Line, name))
				}
			}
		}
	}
	sort.Strings(problems)
	return problems, nil
}

// parsePackageFiles parses every non-test .go file directly in dir. It
// deliberately avoids the deprecated parser.ParseDir/ast.Package pair.
func parsePackageFiles(dir string) ([]*ast.File, *token.FileSet, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, n), nil, parser.ParseComments)
		if err != nil {
			return nil, nil, err
		}
		files = append(files, f)
	}
	return files, fset, nil
}

// declaredNames collects the package's top-level identifiers plus
// "Type.Method" entries for methods (pointer receivers included).
func declaredNames(files []*ast.File) map[string]bool {
	names := map[string]bool{}
	for _, file := range files {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil || len(d.Recv.List) == 0 {
					names[d.Name.Name] = true
					continue
				}
				if recv := receiverTypeName(d.Recv.List[0].Type); recv != "" {
					names[recv+"."+d.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						names[s.Name.Name] = true
					case *ast.ValueSpec:
						for _, n := range s.Names {
							names[n.Name] = true
						}
					}
				}
			}
		}
	}
	return names
}

func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.IndexExpr: // generic receiver
		return receiverTypeName(t.X)
	}
	return ""
}

// docComments returns the comment groups godoc renders: the package/file
// doc plus declaration, spec, and struct-field docs (leading and trailing).
func docComments(file *ast.File) []*ast.CommentGroup {
	var out []*ast.CommentGroup
	add := func(cgs ...*ast.CommentGroup) {
		for _, cg := range cgs {
			if cg != nil {
				out = append(out, cg)
			}
		}
	}
	add(file.Doc)
	ast.Inspect(file, func(n ast.Node) bool {
		switch d := n.(type) {
		case *ast.FuncDecl:
			add(d.Doc)
		case *ast.GenDecl:
			add(d.Doc)
		case *ast.TypeSpec:
			add(d.Doc, d.Comment)
		case *ast.ValueSpec:
			add(d.Doc, d.Comment)
		case *ast.Field:
			add(d.Doc, d.Comment)
		}
		return true
	})
	return out
}
