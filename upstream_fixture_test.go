package h3

// Pure-Go equivalents of the input-driven H3 v4.5.0 test programs:
//   src/apps/testapps/testLatLngToCell.c
//   src/apps/testapps/testCellToLatLng.c
//   src/apps/testapps/testCellToBoundary.c
//
// The large upstream inputs remain under testref/ and are never vendored.
// Run with `make test-upstream-fixtures` after fetching the reference tree.

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func upstreamFixtureRoot(t *testing.T) string {
	t.Helper()
	root := os.Getenv("H3_UPSTREAM_FIXTURE_ROOT")
	if root == "" {
		t.Skip("set H3_UPSTREAM_FIXTURE_ROOT (or run make test-upstream-fixtures)")
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("upstream fixture root %q: %v", root, err)
	}
	return root
}

func fixtureFiles(t *testing.T, root string, patterns ...string) []string {
	t.Helper()
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(root, pattern))
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		t.Fatalf("no upstream fixtures matched %v under %s", patterns, root)
	}
	return files
}

func TestUpstreamLatLngToCellFixtures(t *testing.T) {
	root := upstreamFixtureRoot(t)
	for _, path := range fixtureFiles(t, root, "bc*centers.txt", "rand*centers.txt") {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			in := bufio.NewReader(f)
			for record := 1; ; record++ {
				var text string
				var latDeg, lngDeg float64
				_, err := fmt.Fscan(in, &text, &latDeg, &lngDeg)
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Fatalf("record %d: %v", record, err)
				}
				want, err := ParseCell(text)
				if err != nil {
					t.Fatalf("record %d: ParseCell(%q): %v", record, text, err)
				}
				got, err := LatLngToCell(LatLngDegs(latDeg, lngDeg), want.Resolution())
				if err != nil || got != want {
					t.Fatalf("record %d: LatLngToCell(%g, %g, %d) = %v, %v; want %v",
						record, latDeg, lngDeg, want.Resolution(), got, err, want)
				}
			}
		})
	}
}

func TestUpstreamCellToLatLngFixtures(t *testing.T) {
	root := upstreamFixtureRoot(t)
	const epsilon = 0.000001 * RadPerDeg
	for _, path := range fixtureFiles(t, root, "res*ic.txt") {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			in := bufio.NewReader(f)
			for record := 1; ; record++ {
				var text string
				var latDeg, lngDeg float64
				_, err := fmt.Fscan(in, &text, &latDeg, &lngDeg)
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Fatalf("record %d: %v", record, err)
				}
				cell, err := ParseCell(text)
				if err != nil {
					t.Fatalf("record %d: ParseCell(%q): %v", record, text, err)
				}
				got, err := cell.LatLng()
				if err != nil {
					t.Fatalf("record %d: Cell.LatLng(%v): %v", record, cell, err)
				}
				want := LatLngDegs(latDeg, lngDeg)
				if !geoAlmostEqualThreshold(&got, &want, epsilon) {
					t.Fatalf("record %d: Cell.LatLng(%v) = %v; want %v", record, cell, got, want)
				}
				back, err := LatLngToCell(got, cell.Resolution())
				if err != nil || back != cell {
					t.Fatalf("record %d: center round trip = %v, %v; want %v", record, back, err, cell)
				}
			}
		})
	}
}

func TestUpstreamCellToBoundaryFixtures(t *testing.T) {
	root := upstreamFixtureRoot(t)
	for _, path := range fixtureFiles(t, root, "*cells.txt") {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			scan := bufio.NewScanner(f)
			scan.Split(bufio.ScanWords)
			for record := 1; scan.Scan(); record++ {
				cell, errC := stringToH3(scan.Text())
				if errC != eSuccess {
					t.Fatalf("record %d: stringToH3: %v", record, errC)
				}
				if !scan.Scan() || scan.Text() != "{" {
					t.Fatalf("record %d: missing boundary opening brace", record)
				}
				var want []LatLng
				for scan.Scan() && scan.Text() != "}" {
					latDeg, err := strconv.ParseFloat(scan.Text(), 64)
					if err != nil || !scan.Scan() {
						t.Fatalf("record %d: malformed latitude: %v", record, err)
					}
					lngDeg, err := strconv.ParseFloat(scan.Text(), 64)
					if err != nil {
						t.Fatalf("record %d: malformed longitude: %v", record, err)
					}
					want = append(want, LatLngDegs(latDeg, lngDeg))
				}
				var got CellBoundary
				if errC := cellToBoundary(cell, &got); errC != eSuccess {
					t.Fatalf("record %d: cellToBoundary(%v): %v", record, cell, errC)
				}
				if got.Len() != len(want) {
					t.Fatalf("record %d: boundary vertex count %d; want %d", record, got.Len(), len(want))
				}
				for i := range want {
					if !geoAlmostEqual(&got.verts[i], &want[i]) {
						t.Fatalf("record %d vertex %d: got %v; want %v", record, i, got.verts[i], want[i])
					}
				}
			}
			if err := scan.Err(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
