// Command cliinventory discovers and verifies the H3 v4.4.0 command-line contract.
// It extracts the registered subcommand set from h3.c and every balanced
// add_h3_cli_test(...) registration from tests/cli/*.txt.
package main

import (
	"crypto/sha256"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type cliCase struct {
	source, name, command, expected, digest string
}

type sourceFile struct{ path, role, digest string }

var commandRe = regexp.MustCompile(`SUBCOMMAND_INDEX\(([A-Za-z0-9_]+)\)`)

func main() {
	upstream := flag.String("upstream", "testref/h3-4.4.0", "upstream H3 tree")
	registry := flag.String("registry", "docs/cli-test-inventory.csv", "committed case registry")
	contract := flag.String("contract", "docs/cli-contract.csv", "committed semantic command registry")
	fixturesRegistry := flag.String("fixtures", "docs/cli-fixture-inventory.csv", "committed fixture registry")
	sourcesRegistry := flag.String("sources", "docs/cli-source-inventory.csv", "committed source registry")
	emit := flag.Bool("emit-cases", false, "write discovered case CSV to stdout")
	emitFixtures := flag.Bool("emit-fixtures", false, "write referenced fixture CSV to stdout")
	emitSources := flag.Bool("emit-sources", false, "write CLI support-source CSV to stdout")
	updateEcosystem := flag.Bool("update-ecosystem-inventory", false, "mark discovered CLI cases full in docs/upstream-test-inventory.csv")
	updateContractMetadata := flag.Bool("update-contract-metadata", false, "add source/input/exit/test metadata to the semantic contract")
	verify := flag.Bool("verify", false, "fail on command/case/source drift")
	flag.Parse()

	commands, sourceDigest, err := scanCommands(*upstream)
	check(err)
	cases, err := scanCases(*upstream)
	check(err)
	fixtures, err := scanFixtures(*upstream, cases)
	check(err)
	sources, err := scanSources(*upstream)
	check(err)
	if *emit {
		writeCases(cases)
		return
	}
	if *emitFixtures {
		writeSources(fixtures)
		return
	}
	if *emitSources {
		writeSources(sources)
		return
	}
	if *updateEcosystem {
		check(updateEcosystemInventory("docs/upstream-test-inventory.csv", cases))
		return
	}
	if *updateContractMetadata {
		check(updateContract(*contract))
		return
	}

	problems := verifyCases(*registry, cases)
	problems = append(problems, verifyContract(*contract, commands, sourceDigest, cases)...)
	problems = append(problems, verifySources(*fixturesRegistry, fixtures)...)
	problems = append(problems, verifySources(*sourcesRegistry, sources)...)
	fmt.Printf("CLI contract: %d commands, %d scenarios.\n", len(commands), len(cases))
	if len(problems) == 0 {
		fmt.Println("OK: CLI command and scenario registries match upstream.")
		return
	}
	for _, problem := range problems {
		fmt.Println("-", problem)
	}
	if *verify {
		os.Exit(1)
	}
}

func scanFixtures(root string, cases []cliCase) ([]sourceFile, error) {
	const marker = "${PROJECT_SOURCE_DIR}/tests/inputfiles/"
	names := map[string]bool{}
	for _, c := range cases {
		for rest := c.command; ; {
			i := strings.Index(rest, marker)
			if i < 0 {
				break
			}
			rest = rest[i+len(marker):]
			end := strings.IndexAny(rest, " \\\"`<>")
			if end < 0 {
				end = len(rest)
			}
			names[rest[:end]] = true
			rest = rest[end:]
		}
	}
	var out []sourceFile
	for name := range names {
		path := filepath.ToSlash(filepath.Join("tests/inputfiles", name))
		digest, err := fileDigest(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return nil, err
		}
		out = append(out, sourceFile{path, "CLI fixture", digest})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].path < out[j].path })
	return out, nil
}

func scanSources(root string) ([]sourceFile, error) {
	roles := map[string]string{
		"CMakeLists.txt":                 "target and installed binary registration",
		"src/apps/filters/h3.c":          "commands, formats, I/O, and exit behavior",
		"src/apps/applib/include/args.h": "argument parser contract",
		"src/apps/applib/lib/args.c":     "argument parser implementation",
	}
	var out []sourceFile
	for path, role := range roles {
		digest, err := fileDigest(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return nil, err
		}
		out = append(out, sourceFile{path, role, digest})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].path < out[j].path })
	return out, nil
}

func fileDigest(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func scanCommands(root string) ([]string, string, error) {
	path := filepath.Join(root, "src/apps/filters/h3.c")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	var commands []string
	for _, match := range commandRe.FindAllStringSubmatch(string(data), -1) {
		if match[1] == "s" { // SUBCOMMAND_INDEX macro parameter
			continue
		}
		commands = append(commands, match[1])
	}
	return commands, fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func scanCases(root string) ([]cliCase, error) {
	dir := filepath.Join(root, "tests/cli")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []cliCase
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		digest := fmt.Sprintf("%x", sha256.Sum256(data))
		calls, err := extractCalls(string(data))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		for _, call := range calls {
			args, err := cmakeArgs(call)
			if err != nil || len(args) != 3 {
				return nil, fmt.Errorf("%s: add_h3_cli_test has %d args: %v", entry.Name(), len(args), err)
			}
			out = append(out, cliCase{"tests/cli/" + entry.Name(), args[0], args[1], args[2], digest})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out, nil
}

func extractCalls(source string) ([]string, error) {
	const marker = "add_h3_cli_test("
	var calls []string
	for offset := 0; ; {
		i := strings.Index(source[offset:], marker)
		if i < 0 {
			return calls, nil
		}
		start := offset + i + len(marker)
		depth, quoted, escaped := 1, false, false
		for j := start; j < len(source); j++ {
			c := source[j]
			if escaped {
				escaped = false
				continue
			}
			if quoted && c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				quoted = !quoted
				continue
			}
			if quoted {
				continue
			}
			switch c {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					calls = append(calls, source[start:j])
					offset = j + 1
					goto next
				}
			}
		}
		return nil, fmt.Errorf("unterminated add_h3_cli_test call")
	next:
	}
}

func cmakeArgs(call string) ([]string, error) {
	var args []string
	for i := 0; i < len(call); {
		for i < len(call) && (call[i] == ' ' || call[i] == '\n' || call[i] == '\r' || call[i] == '\t') {
			i++
		}
		if i == len(call) {
			break
		}
		if call[i] == '"' {
			start := i
			i++
			escaped := false
			for i < len(call) {
				if escaped {
					escaped = false
					i++
					continue
				}
				if call[i] == '\\' {
					escaped = true
					i++
					continue
				}
				if call[i] == '"' {
					i++
					break
				}
				i++
			}
			value, err := strconv.Unquote(call[start:i])
			if err != nil {
				return nil, err
			}
			args = append(args, value)
			continue
		}
		start := i
		for i < len(call) && !strings.ContainsRune(" \n\r\t", rune(call[i])) {
			i++
		}
		args = append(args, call[start:i])
	}
	return args, nil
}

func writeCases(cases []cliCase) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"source", "name", "command", "expected", "source_sha256"})
	for _, c := range cases {
		_ = w.Write([]string{c.source, c.name, c.command, c.expected, c.digest})
	}
	w.Flush()
}

func writeSources(files []sourceFile) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"source", "role", "source_sha256"})
	for _, file := range files {
		_ = w.Write([]string{file.path, file.role, file.digest})
	}
	w.Flush()
}

func verifySources(path string, want []sourceFile) []string {
	f, err := os.Open(path)
	if err != nil {
		return []string{err.Error()}
	}
	defer func() { _ = f.Close() }()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []string{err.Error()}
	}
	got := map[string][]string{}
	for i, row := range records {
		if i > 0 && len(row) == 3 {
			got[row[0]] = row
		}
	}
	var problems []string
	for _, file := range want {
		row, ok := got[file.path]
		if !ok || row[1] != file.role || row[2] != file.digest {
			problems = append(problems, "unreviewed or changed CLI source: "+file.path)
		}
		delete(got, file.path)
	}
	for path := range got {
		problems = append(problems, "stale CLI source: "+path)
	}
	return problems
}

func updateEcosystemInventory(path string, cases []cliCase) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	records, err := csv.NewReader(f).ReadAll()
	_ = f.Close()
	if err != nil {
		return err
	}
	known := map[string]bool{}
	for _, c := range cases {
		known[c.name] = true
	}
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) != 7 || row[0] != "cli" || !known[row[2]] {
			continue
		}
		row[3] = "full"
		row[4] = "TestUpstreamCLICompatibility/" + row[2]
		row[5] = "Executed against internal/cli with upstream arguments, fixtures, output, stderr, and exit status; selected workflows also run as a real process and the full suite is opt-in differential-tested against C"
	}
	tmp := path + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := csv.NewWriter(out)
	err = w.WriteAll(records)
	w.Flush()
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func verifyCases(path string, want []cliCase) []string {
	f, err := os.Open(path)
	if err != nil {
		return []string{err.Error()}
	}
	defer func() { _ = f.Close() }()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []string{err.Error()}
	}
	got := map[string][]string{}
	for i, row := range records {
		if i > 0 && len(row) == 5 {
			got[row[1]] = row
		}
	}
	var problems []string
	for _, c := range want {
		row, ok := got[c.name]
		if !ok {
			problems = append(problems, "unreviewed CLI scenario: "+c.name)
			continue
		}
		if row[0] != c.source || row[2] != c.command || row[3] != c.expected || row[4] != c.digest {
			problems = append(problems, "changed CLI scenario: "+c.name)
		}
		delete(got, c.name)
	}
	for name := range got {
		problems = append(problems, "stale CLI scenario: "+name)
	}
	return problems
}

func verifyContract(path string, commands []string, digest string, cases []cliCase) []string {
	f, err := os.Open(path)
	if err != nil {
		return []string{err.Error()}
	}
	defer func() { _ = f.Close() }()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []string{err.Error()}
	}
	counts := map[string]int{}
	for _, c := range cases {
		counts[strings.Fields(c.command)[0]]++
	}
	got := map[string][]string{}
	for i, row := range records {
		if i > 0 && len(row) == 11 {
			got[row[0]] = row
		}
	}
	var problems []string
	for _, name := range commands {
		row, ok := got[name]
		if !ok {
			problems = append(problems, "unreviewed CLI command: "+name)
			continue
		}
		if row[4] != strconv.Itoa(counts[name]) {
			problems = append(problems, "CLI scenario count changed: "+name)
		}
		if row[5] != digest {
			problems = append(problems, "h3.c changed; review CLI command: "+name)
		}
		if row[6] != "src/apps/filters/h3.c::SUBCOMMAND("+name+")" ||
			row[8] != "0 success; H3 error code 1-19; recognized parser error 0" ||
			row[9] != "TestUpstreamCLICompatibility" || row[10] == "" {
			problems = append(problems, "CLI metadata requires review: "+name)
		}
		delete(got, name)
	}
	for name := range got {
		problems = append(problems, "stale CLI command: "+name)
	}
	return problems
}

func updateContract(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	records, err := csv.NewReader(f).ReadAll()
	_ = f.Close()
	if err != nil {
		return err
	}
	outRows := [][]string{{"command", "syntax", "formats", "implementation", "test_count", "source_sha256", "upstream_source", "input", "exit_status", "go_test", "difference"}}
	for i, row := range records {
		if i == 0 || len(row) < 6 {
			continue
		}
		input := "flags"
		if strings.Contains(row[1], "--file") {
			input = "inline, file, or stdin (-i --)"
		}
		implementation := row[3]
		if !strings.Contains(implementation, "/") {
			files := map[string]string{
				"indexingCommands": "commands_indexing.go", "gridCommands": "commands_grid.go",
				"hierarchyCommands": "commands_grid.go", "regionCommands": "commands_region_edge.go",
				"edgeCommands": "commands_region_edge.go", "vertexCommands": "commands_region_edge.go",
				"miscCommands": "commands_misc.go",
			}
			implementation = "internal/cli/" + files[implementation] + "::" + implementation
		}
		base := append([]string(nil), row[:6]...)
		base[3] = implementation
		outRows = append(outRows, append(base,
			"src/apps/filters/h3.c::SUBCOMMAND("+row[0]+")", input,
			"0 success; H3 error code 1-19; recognized parser error 0",
			"TestUpstreamCLICompatibility", "none"))
	}
	tmp := path + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := csv.NewWriter(out)
	err = w.WriteAll(outRows)
	w.Flush()
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
