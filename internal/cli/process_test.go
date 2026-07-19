package cli

// Process-level tests: TestBinaryProcessContract builds the real cmd/h3
// binary and checks what only a process can show (pipes, stderr routing,
// exit statuses). TestDifferentialWithCCLI additionally replays every
// registered scenario against the compiled upstream C executable; it is
// opt-in via the H3_CLI_C_BINARY env var (`make test-cli-diff` builds the C
// binary and sets it).

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	h3 "github.com/dimchansky/h3-go"
)

func TestBinaryProcessContract(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	binary := filepath.Join(t.TempDir(), "h3")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command("go", "build", "-o", binary, "./cmd/h3")
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build h3: %v\n%s", err, output)
	}

	t.Run("stdin and stdout", func(t *testing.T) {
		cmd := exec.Command(binary, "greatCircleDistanceKm", "-i", "--")
		cmd.Stdin = strings.NewReader("[[0, 0], [0, 1]]")
		output, err := cmd.Output()
		if err != nil || !strings.HasPrefix(string(output), "111.195") {
			t.Fatalf("output/error = %q, %v", output, err)
		}
	})

	t.Run("operation exit and stderr", func(t *testing.T) {
		cmd := exec.Command(binary, "cellAreaKm2", "-c", "115283473fffffff")
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		err := cmd.Run()
		exit, ok := err.(*exec.ExitError)
		if !ok || exit.ExitCode() != 5 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "Error 5:") {
			t.Fatalf("exit/stdout/stderr = %v, %q, %q", err, stdout.String(), stderr.String())
		}
	})

	t.Run("parser error exit", func(t *testing.T) {
		cmd := exec.Command(binary, "gridDisk", "-k", "1")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil || !strings.Contains(stderr.String(), "Required argument missing") {
			t.Fatalf("error/stderr = %v, %q", err, stderr.String())
		}
	})
}

func TestDifferentialWithCCLI(t *testing.T) {
	cBinary := os.Getenv("H3_CLI_C_BINARY")
	if cBinary == "" {
		t.Skip("set H3_CLI_C_BINARY to the upstream C h3 executable")
	}
	for _, tc := range loadUpstreamCLICases(t) {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			compareDifferential(t, cBinary, tc.command)
		})
	}
	coordinates := []h3.LatLng{
		h3.LatLngDegs(89.9999, 179.9999),
		h3.LatLngDegs(-89.9999, -179.9999),
		h3.LatLngDegs(0, 180),
	}
	for res := 0; res <= 15; res++ {
		compareDifferential(t, cBinary, fmt.Sprintf("getNumCells -r %d", res))
		compareDifferential(t, cBinary, fmt.Sprintf("getPentagons -r %d", res))
		pentagons, err := h3.Pentagons(res)
		if err != nil {
			t.Fatal(err)
		}
		compareDifferential(t, cBinary, fmt.Sprintf("gridDisk -c %s -k 2", pentagons[0]))
		for _, ll := range coordinates {
			command := fmt.Sprintf("latLngToCell -r %d --lat %.10f --lng %.10f", res, ll.Lat.Deg(), ll.Lng.Deg())
			compareDifferential(t, cBinary, command)
			cell, err := h3.LatLngToCell(ll, res)
			if err != nil {
				t.Fatal(err)
			}
			compareDifferential(t, cBinary, fmt.Sprintf("cellToBoundary -c %s", cell))
		}
	}
	var batch strings.Builder
	for cell := range h3.CellsAtRes(2) {
		batch.WriteString(cell.String())
		batch.WriteByte('\n')
	}
	batchPath := filepath.Join(t.TempDir(), "large-cell-batch.txt")
	if err := os.WriteFile(batchPath, []byte(batch.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	compareDifferential(t, cBinary, "compactCells -i "+batchPath)
}

func compareDifferential(t *testing.T, cBinary, commandLine string) {
	t.Helper()
	argv, stdin, merge, transform := prepareInvocation(t, commandLine)
	var goOut, goErr bytes.Buffer
	goCode := Run(argv, bytes.NewReader(stdin), &goOut, &goErr)
	cmd := exec.Command(cBinary, argv...)
	cmd.Stdin = bytes.NewReader(stdin)
	var cOut, cErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &cOut, &cErr
	err := cmd.Run()
	cCode := 0
	if exit, ok := err.(*exec.ExitError); ok {
		cCode = exit.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	goText := normalizedScenarioOutput(goOut.String(), goErr.String(), merge, transform)
	cText := normalizedScenarioOutput(cOut.String(), cErr.String(), merge, transform)
	if goCode != cCode || !equivalentOutput(goText, cText) {
		t.Fatalf("%s: Go=(%d,%q,%q) C=(%d,%q,%q)", commandLine, goCode, goOut.String(), goErr.String(), cCode, cOut.String(), cErr.String())
	}
}

func equivalentOutput(a, b string) bool {
	if a == b || equivalentScalarFloat(a, b) {
		return true
	}
	var av, bv any
	if json.Unmarshal([]byte(a), &av) != nil || json.Unmarshal([]byte(b), &bv) != nil {
		return false
	}
	return equivalentJSON(av, bv)
}

func equivalentJSON(a, b any) bool {
	switch av := a.(type) {
	case float64:
		bv, ok := b.(float64)
		return ok && math.Abs(av-bv) <= math.Max(5e-8, math.Abs(bv)*1e-12)
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	case nil:
		return b == nil
	case []any:
		bv, ok := b.([]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !equivalentJSON(av[i], bv[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func normalizedScenarioOutput(stdout, stderr string, merge bool, transform string) string {
	actual := stdout
	if merge {
		actual += stderr
	}
	switch transform {
	case "commas":
		return strings.NewReplacer("\r\n", ",", "\n", ",", "\r", ",").Replace(stdout)
	case "float3":
		return strings.TrimSpace(stdout)
	case "linecount":
		return strconv.Itoa(strings.Count(stdout, "\n"))
	default:
		return strings.TrimRight(actual, "\r\n")
	}
}
