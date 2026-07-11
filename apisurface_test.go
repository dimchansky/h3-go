package h3

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestAPISurface locks the exported surface of the package against
// docs/api-surface.txt. Any accidental export or removal fails this test.
// Regenerate the golden file with:
//
//	UPDATE_API_SURFACE=1 go test -run TestAPISurface .
func TestAPISurface(t *testing.T) {
	t.Parallel()

	got := strings.Join(exportedSurface(t), "\n") + "\n"

	const golden = "docs/api-surface.txt"
	if os.Getenv("UPDATE_API_SURFACE") != "" {
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("%s updated", golden)
		return
	}

	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read %s: %v (regenerate with UPDATE_API_SURFACE=1)", golden, err)
	}
	if got != string(want) {
		t.Errorf("exported API surface changed; if intentional, regenerate %s with UPDATE_API_SURFACE=1\n--- got ---\n%s\n--- want ---\n%s",
			golden, got, want)
	}
}

// exportedSurface renders one sorted line per exported declaration of the
// production package: files in the repository root, excluding tests and the
// cgo parity harness (build constraints mentioning c2go).
func exportedSurface(t *testing.T) []string {
	t.Helper()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var lines []string
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if constraintExcluded(string(src)) {
			continue
		}
		af, err := parser.ParseFile(fset, f, src, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range af.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				if !d.Name.IsExported() {
					continue
				}
				if d.Recv != nil {
					recv := recvTypeName(d.Recv)
					if !ast.IsExported(strings.TrimPrefix(recv, "*")) {
						continue
					}
					lines = append(lines, fmt.Sprintf("method (%s) %s", recv, d.Name.Name))
				} else {
					lines = append(lines, "func "+d.Name.Name)
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() {
							kind := "type"
							if s.Assign != token.NoPos {
								kind = "type-alias"
							}
							lines = append(lines, kind+" "+s.Name.Name)
						}
					case *ast.ValueSpec:
						for _, n := range s.Names {
							if n.IsExported() {
								lines = append(lines, d.Tok.String()+" "+n.Name)
							}
						}
					}
				}
			}
		}
	}
	sort.Strings(lines)
	return lines
}

func recvTypeName(fl *ast.FieldList) string {
	if len(fl.List) == 0 {
		return "?"
	}
	switch e := fl.List[0].Type.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return "*" + id.Name
		}
	}
	return "?"
}

func constraintExcluded(src string) bool {
	for line := range strings.SplitSeq(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//go:build") {
			return strings.Contains(trimmed, "c2go")
		}
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			return false // reached package clause without constraint
		}
	}
	return false
}
