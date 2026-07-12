package cli

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type upstreamCLICase struct {
	name, command, expected string
}

func TestUpstreamCLICompatibility(t *testing.T) {
	cases := loadUpstreamCLICases(t)
	if len(cases) != 170 {
		t.Fatalf("loaded %d upstream CLI cases; want 170", len(cases))
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			argv, stdin, merge, transform := prepareInvocation(t, tc.command)
			var stdout, stderr bytes.Buffer
			code := Run(argv, bytes.NewReader(stdin), &stdout, &stderr)
			actual := stdout.String()
			if merge {
				actual += stderr.String()
			} else if stderr.Len() != 0 {
				t.Fatalf("unexpected stderr: %q", stderr.String())
			}
			actual = strings.TrimRight(actual, "\r\n")
			switch transform {
			case "commas":
				actual = strings.NewReplacer("\r\n", ",", "\n", ",", "\r", ",").Replace(stdout.String())
			case "float3":
				value, err := strconv.ParseFloat(strings.TrimSpace(stdout.String()), 64)
				if err != nil {
					t.Fatal(err)
				}
				actual = fmt.Sprintf("%.3f", value)
			}
			if actual != tc.expected && !equivalentScalarFloat(actual, tc.expected) {
				t.Fatalf("output mismatch\ncommand: %s\n got: %q\nwant: %q\nstderr: %q", tc.command, actual, tc.expected, stderr.String())
			}
			wantCode := expectedExit(tc.expected)
			if code != wantCode {
				t.Fatalf("exit code = %d; want %d (stderr %q)", code, wantCode, stderr.String())
			}
		})
	}
}

func equivalentScalarFloat(actual, expected string) bool {
	a, errA := strconv.ParseFloat(actual, 64)
	b, errB := strconv.ParseFloat(expected, 64)
	if errA != nil || errB != nil {
		return false
	}
	tolerance := math.Max(1e-12, math.Abs(b)*1e-12)
	return math.Abs(a-b) <= tolerance
}

func loadUpstreamCLICases(t *testing.T) []upstreamCLICase {
	t.Helper()
	f, err := os.Open("testdata/upstream-cli-cases.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	out := make([]upstreamCLICase, 0, len(records)-1)
	for i, row := range records {
		if i == 0 {
			continue
		}
		out = append(out, upstreamCLICase{name: row[1], command: row[2], expected: row[3]})
	}
	return out
}

func expectedExit(expected string) int {
	if strings.HasPrefix(expected, "Error ") {
		fields := strings.Fields(expected)
		code, _ := strconv.Atoi(strings.TrimSuffix(fields[1], ":"))
		return code
	}
	if strings.HasPrefix(expected, "Only two pairs") {
		return 1
	}
	return 0
}

func prepareInvocation(t *testing.T, commandLine string) (argv []string, stdin []byte, merge bool, transform string) {
	t.Helper()
	commandLine = strings.ReplaceAll(commandLine, "${PROJECT_SOURCE_DIR}/tests/inputfiles", "testdata/fixtures")
	for strings.Contains(commandLine, "\\`") {
		commandLine = strings.ReplaceAll(commandLine, "\\`", "`")
	}
	merge = strings.Contains(commandLine, "2>&1")
	commandLine = strings.ReplaceAll(commandLine, "2>&1", "")
	if i := strings.Index(commandLine, " | "); i >= 0 {
		pipeline := commandLine[i+3:]
		commandLine = commandLine[:i]
		if strings.Contains(pipeline, "xargs printf") {
			transform = "float3"
		} else {
			transform = "commas"
		}
	}
	fields := shellFields(t, commandLine)
	for i := 0; i < len(fields); i++ {
		if fields[i] == "<" && i+1 < len(fields) {
			var err error
			stdin, err = os.ReadFile(fields[i+1])
			if err != nil {
				t.Fatal(err)
			}
			fields = append(fields[:i], fields[i+2:]...)
			break
		}
	}
	return fields, stdin, merge, transform
}

func shellFields(t *testing.T, input string) []string {
	t.Helper()
	var fields []string
	var current strings.Builder
	quote := byte(0)
	flush := func() {
		if current.Len() > 0 {
			fields = append(fields, current.String())
			current.Reset()
		}
	}
	for i := 0; i < len(input); i++ {
		c := input[i]
		if quote == '`' {
			end := strings.IndexByte(input[i:], '`')
			if end < 0 {
				t.Fatalf("unterminated backtick in %q", input)
			}
			inside := input[i : i+end]
			parts := strings.Fields(inside)
			if len(parts) != 2 || parts[0] != "cat" {
				t.Fatalf("unsupported substitution %q", inside)
			}
			data, err := os.ReadFile(filepath.Clean(parts[1]))
			if err != nil {
				t.Fatal(err)
			}
			current.Write(data)
			i += end
			quote = 0
			continue
		}
		if quote != 0 {
			if quote == '"' && c == '`' {
				end := strings.IndexByte(input[i+1:], '`')
				if end < 0 {
					t.Fatalf("unterminated backtick in %q", input)
				}
				inside := input[i+1 : i+1+end]
				parts := strings.Fields(inside)
				if len(parts) != 2 || parts[0] != "cat" {
					t.Fatalf("unsupported substitution %q", inside)
				}
				data, err := os.ReadFile(filepath.Clean(parts[1]))
				if err != nil {
					t.Fatal(err)
				}
				current.Write(data)
				i += end + 1
				continue
			}
			if c == quote {
				quote = 0
			} else {
				current.WriteByte(c)
			}
			continue
		}
		switch c {
		case '\'', '"', '`':
			quote = c
		case ' ', '\t', '\r', '\n':
			flush()
		default:
			current.WriteByte(c)
		}
	}
	flush()
	return fields
}

func TestParserAndTopLevelExitContract(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		code           int
		stdout, stderr string
	}{
		{name: "no args", code: 1, stdout: "Please use h3 --help"},
		{name: "unknown", args: []string{"unknown"}, code: 1, stdout: "Please use h3 --help"},
		{name: "general help", args: []string{"--help"}, code: 0, stdout: "cellToLatLng"},
		{name: "version", args: []string{"--version"}, code: 0, stdout: "h3 4.4.0"},
		{name: "command help", args: []string{"gridDisk", "--help"}, code: 0, stdout: "H3 4.4.0"},
		{name: "missing required exits zero", args: []string{"gridDisk", "-k", "1"}, code: 0, stderr: "Required argument missing"},
		{name: "duplicate exits zero", args: []string{"getNumCells", "-r", "1", "--resolution", "2"}, code: 0, stderr: "Argument specified multiple times"},
		{name: "unknown option exits zero", args: []string{"pentagonCount", "--wat"}, code: 0, stderr: "Unknown argument"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(tc.args, strings.NewReader(""), &stdout, &stderr)
			if code != tc.code || !strings.Contains(stdout.String(), tc.stdout) || !strings.Contains(stderr.String(), tc.stderr) {
				t.Fatalf("code/stdout/stderr = %d, %q, %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, fmt.Errorf("broken pipe") }

func TestOutputFailureReturnsNonzero(t *testing.T) {
	if code := Run([]string{"getRes0Cells"}, strings.NewReader(""), failingWriter{}, &bytes.Buffer{}); code != 1 {
		t.Fatalf("exit code = %d; want 1", code)
	}
}
