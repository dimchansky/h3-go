package cli

import (
	"io"
	"strconv"
	"strings"
)

type optionSpec struct {
	names     []string
	value     bool
	required  bool
	valueName string
	help      string
}

type parsedArgs struct {
	values map[string]string
	found  map[string]bool
}

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
