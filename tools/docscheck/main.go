// Command docscheck verifies the repository's Markdown documentation:
// every relative link must point to an existing file or directory, and
// every fragment (#anchor) into a Markdown file must match a heading
// anchor as GitHub generates them.
//
// Usage:
//
//	go run ./tools/docscheck            # scan the repository (make check-docs)
//	go run ./tools/docscheck -root DIR  # scan another tree
//
// The tool is read-only and dependency-free. It scans all *.md files
// except generated/downloaded trees (testref/, .gocache/, .git/, .idea/),
// skips fenced code blocks and inline code spans, and ignores absolute
// URLs (http/https/mailto). It exits 1 and lists each broken link as
// "file:line: message" when problems are found.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

var (
	linkRe       = regexp.MustCompile(`!?\[[^\]]*\]\(([^()]+)\)`)
	codeSpanRe   = regexp.MustCompile("`[^`]*`")
	headingRe    = regexp.MustCompile(`^#{1,6}\s+(.*?)\s*#*\s*$`)
	headingLink  = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	skippedDirs  = map[string]bool{".git": true, ".gocache": true, ".idea": true, "testref": true, "node_modules": true}
	anchorsCache = map[string][]string{}
)

func main() {
	root := flag.String("root", ".", "repository root to scan")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "docscheck verifies relative links and #anchors in the repository's Markdown files.")
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/docscheck [-root dir]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var mdFiles []string
	err := filepath.WalkDir(*root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skippedDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".md") {
			mdFiles = append(mdFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "docscheck: %v\n", err)
		os.Exit(1)
	}

	problems := 0
	for _, file := range mdFiles {
		for _, p := range checkFile(*root, file) {
			fmt.Println(p)
			problems++
		}
	}
	if problems > 0 {
		fmt.Fprintf(os.Stderr, "docscheck: %d broken link(s) in %d Markdown file(s)\n", problems, len(mdFiles))
		os.Exit(1)
	}
	fmt.Printf("docscheck: OK — %d Markdown files checked\n", len(mdFiles))
}

func checkFile(root, file string) []string {
	data, err := os.ReadFile(file)
	if err != nil {
		return []string{fmt.Sprintf("%s: %v", file, err)}
	}
	var problems []string
	inFence := false
	for i, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		for _, m := range linkRe.FindAllStringSubmatch(codeSpanRe.ReplaceAllString(line, ""), -1) {
			if p := checkLink(root, file, m[1]); p != "" {
				problems = append(problems, fmt.Sprintf("%s:%d: %s", file, i+1, p))
			}
		}
	}
	return problems
}

func checkLink(root, file, rawTarget string) string {
	target := strings.Trim(strings.Fields(rawTarget)[0], "<>")
	if strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
		return ""
	}
	path, frag, _ := strings.Cut(target, "#")
	resolved := file
	if path != "" {
		if decoded, err := url.PathUnescape(path); err == nil {
			path = decoded
		}
		if strings.HasPrefix(path, "/") {
			resolved = filepath.Join(root, path)
		} else {
			resolved = filepath.Join(filepath.Dir(file), path)
		}
		if _, err := os.Stat(resolved); err != nil {
			return fmt.Sprintf("broken link %q (%s does not exist)", target, resolved)
		}
	}
	if frag == "" || !strings.EqualFold(filepath.Ext(resolved), ".md") {
		return ""
	}
	anchors, ok := anchorsCache[resolved]
	if !ok {
		anchors = headingAnchors(resolved)
		anchorsCache[resolved] = anchors
	}
	if slices.Contains(anchors, strings.ToLower(frag)) {
		return ""
	}
	return fmt.Sprintf("broken anchor %q (no heading generates #%s in %s)", target, frag, resolved)
}

// headingAnchors extracts the anchors GitHub generates for a file's ATX
// headings: strip markdown decoration, lowercase, drop everything except
// letters/digits/spaces/hyphens/underscores, then turn spaces into hyphens;
// repeated anchors get -1, -2, ... suffixes.
func headingAnchors(file string) []string {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	var anchors []string
	seen := map[string]int{}
	inFence := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		m := headingRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		text := headingLink.ReplaceAllString(m[1], "$1")
		text = strings.NewReplacer("`", "", "*", "", "~", "").Replace(text)
		var b strings.Builder
		for _, r := range strings.ToLower(text) {
			switch {
			case r == ' ':
				b.WriteRune('-')
			case r == '-' || r == '_' ||
				('a' <= r && r <= 'z') || ('0' <= r && r <= '9') ||
				(r > 127 && (unicode.IsLetter(r) || unicode.IsDigit(r))):
				b.WriteRune(r)
			}
		}
		anchor := b.String()
		if n, dup := seen[anchor]; dup {
			seen[anchor] = n + 1
			anchor = fmt.Sprintf("%s-%d", anchor, n)
		} else {
			seen[anchor] = 1
		}
		anchors = append(anchors, anchor)
	}
	return anchors
}
