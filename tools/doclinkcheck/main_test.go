package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFixtureUnresolvedLinkDetected(t *testing.T) {
	t.Parallel()
	problems, err := run(filepath.Join("testdata", "badpkg"))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 {
		t.Fatalf("problems = %v, want exactly one (NoSuchSymbol)", problems)
	}
	if !strings.Contains(problems[0], "[NoSuchSymbol]") {
		t.Fatalf("problem %q does not name NoSuchSymbol", problems[0])
	}
}

func TestRootPackageDocLinksResolve(t *testing.T) {
	t.Parallel()
	problems, err := run(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("root package has unresolved doc links:\n%s", strings.Join(problems, "\n"))
	}
}
