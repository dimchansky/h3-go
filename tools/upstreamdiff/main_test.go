package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTree materializes a minimal synthetic upstream tree (src/h3lib/lib +
// an empty src/h3lib/include) with the given lib files.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	lib := filepath.Join(root, "src", "h3lib", "lib")
	inc := filepath.Join(root, "src", "h3lib", "include")
	for _, d := range []string{lib, inc} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(lib, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func scan(t *testing.T, root string) map[string]symbol {
	t.Helper()
	syms, err := scanTree(root)
	if err != nil {
		t.Fatal(err)
	}
	return syms
}

const baseSrc = `/** Frobnicates the cell.
 * @param h cell index
 */
H3Error frobnicate(H3Index h) {
    return doWork(h); // interior note
}
`

func TestDiffLeadingCommentOnlyChange(t *testing.T) {
	t.Parallel()
	oldTree := scan(t, writeTree(t, map[string]string{"demo.c": baseSrc}))
	newTree := scan(t, writeTree(t, map[string]string{"demo.c": `/** Frobnicates the cell. NEW CONTRACT NOTE.
 * @param h cell index
 */
H3Error frobnicate(H3Index h) {
    return doWork(h); // interior note
}
`}))

	// Default mode: leading-comment-only changes are invisible.
	if rows := diffSymbols(oldTree, newTree, false); len(rows) != 0 {
		t.Fatalf("default mode reported %d changes for a leading-comment-only diff, want 0: %+v", len(rows), rows)
	}
	// -comments mode: reported, flagged comment-only.
	rows := diffSymbols(oldTree, newTree, true)
	if len(rows) != 1 {
		t.Fatalf("comments mode reported %d changes, want 1: %+v", len(rows), rows)
	}
	r := rows[0]
	if r.status != "changed" || r.key != "demo.c::frobnicate" ||
		r.bodyChanged || !r.commentChanged {
		t.Fatalf("unexpected change row: %+v", r)
	}
	if got := changeDims(r); got != "comment-only" {
		t.Fatalf("changeDims = %q, want comment-only", got)
	}
}

func TestDiffInteriorBodyCommentChange(t *testing.T) {
	t.Parallel()
	oldTree := scan(t, writeTree(t, map[string]string{"demo.c": baseSrc}))
	newTree := scan(t, writeTree(t, map[string]string{"demo.c": `/** Frobnicates the cell.
 * @param h cell index
 */
H3Error frobnicate(H3Index h) {
    return doWork(h); // CHANGED interior note
}
`}))

	// A comment inside the captured body is a body change in both modes —
	// pinning the tool's historical behavior.
	for _, includeComments := range []bool{false, true} {
		rows := diffSymbols(oldTree, newTree, includeComments)
		if len(rows) != 1 {
			t.Fatalf("comments=%v: %d changes, want 1: %+v", includeComments, len(rows), rows)
		}
		r := rows[0]
		if !r.bodyChanged || r.commentChanged {
			t.Fatalf("comments=%v: unexpected dimensions: %+v", includeComments, r)
		}
	}
}

func TestDiffBodyAndLeadingCommentChange(t *testing.T) {
	t.Parallel()
	oldTree := scan(t, writeTree(t, map[string]string{"demo.c": baseSrc}))
	newTree := scan(t, writeTree(t, map[string]string{"demo.c": `/** Frobnicates the cell. NEW CONTRACT NOTE.
 * @param h cell index
 */
H3Error frobnicate(H3Index h) {
    return doWorkDifferently(h);
}
`}))

	// Both dimensions must be visible — the comment change is not folded
	// into "body changed".
	rows := diffSymbols(oldTree, newTree, true)
	if len(rows) != 1 {
		t.Fatalf("comments mode reported %d changes, want 1: %+v", len(rows), rows)
	}
	r := rows[0]
	if !r.bodyChanged || !r.commentChanged {
		t.Fatalf("want both dimensions flagged, got %+v", r)
	}
	if got := changeDims(r); got != "body+comment" {
		t.Fatalf("changeDims = %q, want body+comment", got)
	}
	// ...and default mode still reports it (as a plain body change).
	if rows := diffSymbols(oldTree, newTree, false); len(rows) != 1 || rows[0].commentChanged {
		t.Fatalf("default mode rows: %+v", rows)
	}
}

func TestLeadingCommentAttachment(t *testing.T) {
	t.Parallel()
	syms := scan(t, writeTree(t, map[string]string{"demo.c": `/*
 * License header — separated by a blank line, must not attach.
 */

// Table of frobnication factors.
static const int FACTORS[2] = {1, 2};

#include "other.h"
H3Error bare(H3Index h) {
    return 0;
}
`}))
	table, ok := syms["demo.c::FACTORS"]
	if !ok {
		t.Fatal("FACTORS not extracted")
	}
	if want := norm("// Table of frobnication factors."); table.comment != want {
		t.Fatalf("FACTORS comment = %q, want %q", table.comment, want)
	}
	bare, ok := syms["demo.c::bare"]
	if !ok {
		t.Fatal("bare not extracted")
	}
	if bare.comment != "" {
		t.Fatalf("bare should have no leading comment (preprocessor breaks attachment), got %q", bare.comment)
	}
}
