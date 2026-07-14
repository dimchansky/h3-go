package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// fixedEpoch stands in for SOURCE_DATE_EPOCH; an even second keeps the
// 2-second MS-DOS zip timestamp resolution exact.
var fixedEpoch = time.Unix(1752500000, 0).UTC()

// stageInput writes deterministic fake inputs and returns them in
// deliberately unsorted order to prove the writers sort.
func stageInput(t *testing.T) []stagedFile {
	t.Helper()
	dir := t.TempDir()
	write := func(name, content string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	return []stagedFile{
		{name: "NOTICE", src: write("NOTICE", "notice\n"), mode: docMode},
		{name: "h3", src: write("h3", "\x7fELF fake binary"), mode: binMode},
		{name: "LICENSE", src: write("LICENSE", "license\n"), mode: docMode},
		{name: "README.md", src: write("README.md", "# readme\n"), mode: docMode},
	}
}

// wantOrder is the canonical entry order: top directory, then sorted names.
var wantOrder = []string{"top/", "top/LICENSE", "top/NOTICE", "top/README.md", "top/h3"}

func TestTarGzDeterministicBytes(t *testing.T) {
	files := stageInput(t)
	out := t.TempDir()
	a, b := filepath.Join(out, "a.tar.gz"), filepath.Join(out, "b.tar.gz")
	if err := writeTarGz(a, "top", files, fixedEpoch); err != nil {
		t.Fatal(err)
	}
	if err := writeTarGz(b, "top", files, fixedEpoch); err != nil {
		t.Fatal(err)
	}
	assertIdenticalFiles(t, a, b)
}

func TestZipDeterministicBytes(t *testing.T) {
	files := stageInput(t)
	out := t.TempDir()
	a, b := filepath.Join(out, "a.zip"), filepath.Join(out, "b.zip")
	if err := writeZip(a, "top", files, fixedEpoch); err != nil {
		t.Fatal(err)
	}
	if err := writeZip(b, "top", files, fixedEpoch); err != nil {
		t.Fatal(err)
	}
	assertIdenticalFiles(t, a, b)
}

func TestTarGzStructure(t *testing.T) {
	out := filepath.Join(t.TempDir(), "x.tar.gz")
	if err := writeTarGz(out, "top", stageInput(t), fixedEpoch); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	if !gz.ModTime.IsZero() {
		t.Errorf("gzip MTIME = %v, want zero", gz.ModTime)
	}
	if gz.Name != "" || gz.Comment != "" {
		t.Errorf("gzip name/comment = %q/%q, want empty", gz.Name, gz.Comment)
	}
	if gz.OS != 255 {
		t.Errorf("gzip OS byte = %d, want 255 (unknown)", gz.OS)
	}

	tr := tar.NewReader(gz)
	var got []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, hdr.Name)
		if hdr.Uid != 0 || hdr.Gid != 0 || hdr.Uname != "" || hdr.Gname != "" {
			t.Errorf("%s: ownership %d/%d %q/%q, want 0/0 with empty names", hdr.Name, hdr.Uid, hdr.Gid, hdr.Uname, hdr.Gname)
		}
		if !hdr.ModTime.Equal(fixedEpoch) {
			t.Errorf("%s: mtime %v, want %v", hdr.Name, hdr.ModTime, fixedEpoch)
		}
		if hdr.Format != tar.FormatUSTAR {
			t.Errorf("%s: format %v, want USTAR", hdr.Name, hdr.Format)
		}
		wantMode := int64(docMode)
		switch hdr.Name {
		case "top/":
			wantMode = dirMode
		case "top/h3":
			wantMode = binMode
		}
		if hdr.Mode != wantMode {
			t.Errorf("%s: mode %o, want %o", hdr.Name, hdr.Mode, wantMode)
		}
	}
	assertOrder(t, got)
}

func TestZipStructure(t *testing.T) {
	out := filepath.Join(t.TempDir(), "x.zip")
	if err := writeZip(out, "top", stageInput(t), fixedEpoch); err != nil {
		t.Fatal(err)
	}
	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = zr.Close() }()

	var got []string
	for _, f := range zr.File {
		got = append(got, f.Name)
		if len(f.Extra) != 0 {
			t.Errorf("%s: %d bytes of extra fields, want none", f.Name, len(f.Extra))
		}
		m := f.Modified.UTC()
		if m.Year() != fixedEpoch.Year() || m.Month() != fixedEpoch.Month() || m.Day() != fixedEpoch.Day() ||
			m.Hour() != fixedEpoch.Hour() || m.Minute() != fixedEpoch.Minute() || m.Second() != fixedEpoch.Second() {
			t.Errorf("%s: modified %v, want %v", f.Name, m, fixedEpoch)
		}
		perm := f.Mode().Perm()
		wantPerm := os.FileMode(docMode)
		switch f.Name {
		case "top/":
			wantPerm = dirMode
		case "top/h3":
			wantPerm = binMode
		}
		if perm != wantPerm {
			t.Errorf("%s: mode %o, want %o", f.Name, perm, wantPerm)
		}
	}
	assertOrder(t, got)
}

func TestSHA256SUMSManifest(t *testing.T) {
	dir := t.TempDir()
	// Unsorted creation order; the manifest must come out sorted.
	if err := os.WriteFile(filepath.Join(dir, "b.zip"), []byte("bb"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.tar.gz"), []byte("aa"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeSHA256SUMS(dir, []string{"b.zip", "a.tar.gz"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("%x  a.tar.gz\n%x  b.zip\n", sha256.Sum256([]byte("aa")), sha256.Sum256([]byte("bb")))
	if string(got) != want {
		t.Errorf("SHA256SUMS:\n%s\nwant:\n%s", got, want)
	}
}

func assertOrder(t *testing.T, got []string) {
	t.Helper()
	if len(got) != len(wantOrder) {
		t.Fatalf("entries = %v, want %v", got, wantOrder)
	}
	for i := range wantOrder {
		if got[i] != wantOrder[i] {
			t.Fatalf("entry[%d] = %q, want %q (full order %v)", i, got[i], wantOrder[i], got)
		}
	}
}

func assertIdenticalFiles(t *testing.T, a, b string) {
	t.Helper()
	da, err := os.ReadFile(a)
	if err != nil {
		t.Fatal(err)
	}
	db, err := os.ReadFile(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(da, db) {
		t.Fatalf("archives differ: %d vs %d bytes", len(da), len(db))
	}
}
