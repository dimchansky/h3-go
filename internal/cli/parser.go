package cli

// This file reimplements the upstream applib argument parser
// (src/apps/applib/lib/args.c) rather than using the flag package, because
// the C contract differs from Go conventions in observable ways:
// case-insensitive option matching, multiple aliases per option, and
// help/parse errors that do not fail the process (Run turns them into exit
// code 0).

import (
	"io"
	"strconv"
	"strings"
)

// optionSpec describes one option of a command: its aliases (e.g. "-r" and
// "--resolution"), whether it consumes a following value, whether it is
// required, and its help line. names[0] is the canonical key runners use
// with parsedArgs regardless of which alias the user typed.
type optionSpec struct {
	names     []string
	value     bool
	required  bool
	valueName string
	help      string
}

// parsedArgs holds parsed option values keyed by each option's canonical
// name (optionSpec.names[0]). found records presence separately so that
// value-less options and empty values are distinguishable.
type parsedArgs struct {
	values map[string]string
	found  map[string]bool
}

// opt builds a value-carrying option spec from space-separated aliases,
// e.g. opt("-c --cell", true). The first alias is the canonical key.
func opt(names string, required bool) optionSpec {
	return optionSpec{names: strings.Fields(names), value: true, required: required}
}

func (p parsedArgs) has(name string) bool   { return p.found[name] }
func (p parsedArgs) get(name string) string { return p.values[name] }

func (p parsedArgs) integer(name string) (int, error) {
	v, err := strconv.ParseInt(p.get(name), 10, 32)
	return int(v), err
}

func (p parsedArgs) int64(name string) (int64, error) {
	return strconv.ParseInt(p.get(name), 10, 64)
}

func (p parsedArgs) float(name string) (float64, error) {
	return strconv.ParseFloat(p.get(name), 64)
}

// parseOptions parses a command's argument vector against its option specs.
// The boolean result tells Run whether to execute the command; false covers
// both help requests (printed to out) and parse errors (help printed to
// errOut with the error prepended) — the two cases upstream treats as
// "handled, exit 0". Matching is case-insensitive, every option implicitly
// gains -h/--help, repeated options and missing values or required options
// are errors, exactly as in upstream args.c.
func parseOptions(program, description string, argv []string, specs []optionSpec, out, errOut io.Writer) (parsedArgs, bool) {
	help := optionSpec{names: []string{"-h", "--help"}, help: "Show this help message."}
	specs = append([]optionSpec{help}, specs...)
	parsed := parsedArgs{values: map[string]string{}, found: map[string]bool{}}
	canonical := map[string]int{}
	for i := range specs {
		for _, name := range specs[i].names {
			canonical[strings.ToLower(name)] = i
		}
	}
	fail := func(message, detail string) (parsedArgs, bool) {
		if detail != "" {
			message += ": " + detail
		}
		printHelp(errOut, program, description, specs, message)
		return parsed, false
	}
	for i := 0; i < len(argv); i++ {
		idx, ok := canonical[strings.ToLower(argv[i])]
		if !ok {
			return fail("Unknown argument", "")
		}
		key := specs[idx].names[0]
		if parsed.found[key] {
			return fail("Argument specified multiple times", argv[i])
		}
		parsed.found[key] = true
		if specs[idx].value {
			i++
			if i == len(argv) {
				return fail("Argument value not present", argv[i-1])
			}
			parsed.values[key] = argv[i]
		}
	}
	if parsed.has("-h") {
		printHelp(out, program, description, specs, "")
		return parsed, false
	}
	for _, spec := range specs {
		if spec.required && !parsed.has(spec.names[0]) {
			return fail("Required argument missing", spec.names[0])
		}
	}
	return parsed, true
}

func printHelp(w io.Writer, program, description string, specs []optionSpec, parseError string) {
	if parseError != "" {
		writef(w, "%s: %s\n", program, parseError)
	}
	writef(w, "%s: %s\nH3 4.4.0\n\n", program, description)
	for _, spec := range specs {
		writeText(w, "\t", strings.Join(spec.names, ", "))
		if spec.value {
			name := spec.valueName
			if name == "" {
				name = "value"
			}
			writef(w, " <%s>", name)
		}
		writeText(w, "\t")
		if spec.required {
			writeText(w, "Required. ")
		}
		writeln(w, spec.help)
	}
}
