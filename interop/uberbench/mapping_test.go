package uberbench

import (
	"encoding/csv"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMappingSymbolsExist verifies the uber/h3-go side of the comparison
// matrix (docs/comparison-uber-h3-go.csv) against the pinned binding
// release: every symbol the matrix claims the binding has must actually be
// exported by it. The offline half of the verification — this library's
// side, row completeness against the C API inventory, and the generated
// doc tables — is tools/ubercompare in the root module; this test needs
// the dependency and therefore lives here.
//
// It fails when a binding release removes or renames an API the matrix
// still lists — the signal to refresh docs/comparison-uber-h3-go.csv.
func TestMappingSymbolsExist(t *testing.T) {
	exported := uberExportedSymbols(t)

	f, err := os.Open("../../docs/comparison-uber-h3-go.csv")
	if err != nil {
		t.Fatalf("open matrix: %v", err)
	}
	defer f.Close()
	rec, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("parse matrix: %v", err)
	}
	if len(rec) < 2 || rec[0][3] != "uber_api" || rec[0][0] != "c_function" {
		t.Fatalf("unexpected matrix header: %v", rec[0])
	}

	checked := 0
	for _, row := range rec[1:] {
		cfunc, cell := row[0], row[3]
		for entry := range strings.SplitSeq(cell, ";") {
			entry = strings.TrimSpace(entry)
			if entry == "" || strings.HasPrefix(entry, "—") {
				continue
			}
			if i := strings.Index(entry, " ("); i >= 0 {
				entry = entry[:i]
			}
			if !exported[entry] {
				t.Errorf("%s: matrix lists uber/h3-go symbol %q, not found in the pinned binding", cfunc, entry)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no uber/h3-go symbols found in the matrix")
	}
	t.Logf("verified %d uber/h3-go symbol references", checked)
}

// uberExportedSymbols parses the pinned binding's sources from the module
// cache and returns its exported surface as a set of "Name" and
// "Recv.Name" strings.
func uberExportedSymbols(t *testing.T) map[string]bool {
	t.Helper()
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/uber/h3-go/v4").Output()
	if err != nil {
		t.Fatalf("locate uber/h3-go module dir: %v", err)
	}
	dir := strings.TrimSpace(string(out))

	// parser.ParseDir is deprecated for not handling build tags, which is
	// fine here: the binding is one flat package and this test only needs
	// its exported names; x/tools/go/packages would add a dependency for
	// nothing.
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse uber/h3-go sources: %v", err)
	}

	syms := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if !d.Name.IsExported() {
						continue
					}
					if d.Recv == nil || len(d.Recv.List) == 0 {
						syms[d.Name.Name] = true
						continue
					}
					if recv := receiverTypeName(d.Recv.List[0].Type); recv != "" && ast.IsExported(recv) {
						syms[recv+"."+d.Name.Name] = true
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.ValueSpec:
							for _, name := range s.Names {
								if name.IsExported() {
									syms[name.Name] = true
								}
							}
						case *ast.TypeSpec:
							if s.Name.IsExported() {
								syms[s.Name.Name] = true
							}
						}
					}
				}
			}
		}
	}
	if len(syms) == 0 {
		t.Fatalf("no exported symbols parsed from %s", dir)
	}
	return syms
}

func receiverTypeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return receiverTypeName(e.X)
	case *ast.IndexExpr: // generic receiver
		return receiverTypeName(e.X)
	}
	return ""
}
