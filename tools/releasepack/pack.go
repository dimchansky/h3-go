package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Deterministic archive rules (release plan D-B): fixed entry order
// (directory first, then files sorted by name), uid/gid 0 with empty
// user/group names, modes 0755 (directories, binaries) and 0644 (docs),
// every mtime equal to SOURCE_DATE_EPOCH (the tagged commit time), a gzip
// header with zero MTIME and no file name, and zip entries with the fixed
// timestamp and no extra fields. Two runs over identical inputs must
// produce byte-identical archives.

const (
	dirMode = 0o755
	docMode = 0o644
	binMode = 0o755
)

// stagedFile is one file placed below the archive's top-level directory.
type stagedFile struct {
	name string      // name inside the archive, below topDir
	src  string      // filesystem path the content is read from
	mode fs.FileMode // binMode for executables, docMode for documents
}

// sortedFiles returns the files in canonical (name-sorted) archive order.
func sortedFiles(files []stagedFile) []stagedFile {
	out := slices.Clone(files)
	slices.SortFunc(out, func(a, b stagedFile) int { return strings.Compare(a.name, b.name) })
	return out
}

// writeTarGz writes topDir/ and its files as a deterministic .tar.gz.
func writeTarGz(outPath, topDir string, files []stagedFile, epoch time.Time) (err error) {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	gz, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	// The gzip.Header zero value is already deterministic: MTIME 0, no
	// name, no comment, OS byte 255 (unknown).
	tw := tar.NewWriter(gz)

	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeDir,
		Name:     topDir + "/",
		Mode:     dirMode,
		ModTime:  epoch,
		Format:   tar.FormatUSTAR,
	}); err != nil {
		return err
	}
	for _, sf := range sortedFiles(files) {
		data, err := os.ReadFile(sf.src)
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     topDir + "/" + sf.name,
			Mode:     int64(sf.mode),
			Size:     int64(len(data)),
			ModTime:  epoch,
			Format:   tar.FormatUSTAR,
		}); err != nil {
			return err
		}
		if _, err := tw.Write(data); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gz.Close()
}

// writeZip writes topDir/ and its files as a deterministic .zip.
func writeZip(outPath, topDir string, files []stagedFile, epoch time.Time) (err error) {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	zw := zip.NewWriter(f)
	dosDate, dosTime := dosDateTime(epoch)
	dir := &zip.FileHeader{
		Name:   topDir + "/",
		Method: zip.Store,
	}
	// Setting Modified makes archive/zip append a 9-byte Extended Timestamp
	// extra field; the release spec demands extra-free entries, so only the
	// legacy MS-DOS fields (interpreted as UTC by the reader) are set.
	dir.ModifiedDate, dir.ModifiedTime = dosDate, dosTime //nolint:staticcheck // deliberate: avoids the extra field Modified would add
	dir.SetMode(fs.ModeDir | dirMode)
	if _, err := zw.CreateHeader(dir); err != nil {
		return err
	}
	for _, sf := range sortedFiles(files) {
		data, err := os.ReadFile(sf.src)
		if err != nil {
			return err
		}
		hdr := &zip.FileHeader{
			Name:   topDir + "/" + sf.name,
			Method: zip.Deflate,
		}
		hdr.ModifiedDate, hdr.ModifiedTime = dosDate, dosTime //nolint:staticcheck // deliberate: avoids the extra field Modified would add
		hdr.SetMode(sf.mode)
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return zw.Close()
}

// dosDateTime converts a UTC time to the legacy MS-DOS date/time fields
// (2-second resolution; zip readers interpret them as UTC).
func dosDateTime(t time.Time) (dosDate, dosTime uint16) {
	t = t.UTC()
	dosDate = uint16(t.Day()) | uint16(t.Month())<<5 | uint16(t.Year()-1980)<<9
	dosTime = uint16(t.Second()/2) | uint16(t.Minute())<<5 | uint16(t.Hour())<<11
	return dosDate, dosTime
}

// writeSHA256SUMS writes a sha256sum-compatible manifest ("<hex>  <name>",
// names sorted) covering the given files in dir.
func writeSHA256SUMS(dir string, names []string) error {
	sorted := slices.Clone(names)
	slices.Sort(sorted)
	var b strings.Builder
	for _, name := range sorted {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "%x  %s\n", sha256.Sum256(data), name)
	}
	return os.WriteFile(filepath.Join(dir, "SHA256SUMS"), []byte(b.String()), docMode)
}
