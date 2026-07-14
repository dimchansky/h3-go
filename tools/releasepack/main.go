// Command releasepack is the single authoritative release-build entry point
// (invoked as `make release-dist VERSION=<tag> OUT=<dir>`). The CI build
// job, the CI verify-reproducible job, and the local release procedure all
// run exactly this implementation, so the archives they produce must be
// byte-identical (full SHA-256 equality of every archive is a release gate).
//
// It validates the release invariants, cross-compiles the h3 CLI for every
// supported platform with a normalized environment, and packs deterministic
// archives:
//
//	preconditions   canonical tag syntax (vX.Y.Z or vX.Y.Z-rc.N per
//	                docs/versioning.md); clean worktree and index; the tag
//	                exists and points at HEAD; `go env GOVERSION` equals the
//	                pinned toolchain; the output directory is empty.
//	build           CGO_ENABLED=0 go build -trimpath
//	                -ldflags "-s -w -buildid= -X …internal/cli.buildVersion=<tag>"
//	                with GOTOOLCHAIN=local, GOENV=off, GOWORK=off, empty
//	                GOFLAGS/GOEXPERIMENT, TZ=UTC, LC_ALL=C, and per-target
//	                architecture baselines (GOAMD64=v1, GOARM64=v8.0) —
//	                inherited host settings cannot leak into the binaries.
//	                SOURCE_DATE_EPOCH is derived from the tagged commit,
//	                never read from the environment.
//	postconditions  every binary's `go version -m` reports the module path,
//	                the tagged commit (vcs.revision), vcs.modified=false,
//	                the target GOOS/GOARCH, and no host paths; the
//	                host-runnable binary is executed and must print
//	                "h3 4.4.0 (<tag>)".
//	output          h3-<tag>-<os>-<arch>.tar.gz (.zip for windows) per
//	                platform, each containing the binary, LICENSE, NOTICE,
//	                and README.md (from cmd/h3/README-archive.md), plus a
//	                sha256sum-compatible SHA256SUMS. Archive metadata is
//	                normalized (see pack.go) so archives — not just
//	                binaries — reproduce bit-for-bit.
//
// The Go toolchain pin lives in requiredGoVersion below and must match the
// release-builds workflow; docs/releasing.md records the bump procedure.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// requiredGoVersion is the exact toolchain releases are built with; the
// release-builds workflow pins the same version for setup-go.
const requiredGoVersion = "go1.26.5"

// versionRE accepts canonical module SemVer release tags plus the -rc.N
// prereleases documented in docs/versioning.md (rc counting starts at 1).
// Build metadata is rejected as project policy: it would not survive as a
// distinct canonical Go module version and must not encode a second axis.
var versionRE = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.[1-9][0-9]*)?$`)

// target is one supported GOOS/GOARCH release platform.
type target struct{ goos, goarch string }

// targets are the supported release platforms, in build order.
var targets = []target{
	{"linux", "amd64"}, {"linux", "arm64"},
	{"darwin", "amd64"}, {"darwin", "arm64"},
	{"windows", "amd64"}, {"windows", "arm64"},
}

func main() {
	version := flag.String("version", "", "release tag to build (vX.Y.Z or vX.Y.Z-rc.N; must point at HEAD)")
	out := flag.String("out", "", "empty output directory for the archives and SHA256SUMS (outside the repository)")
	repo := flag.String("repo", ".", "repository root")
	flag.Parse()
	if *version == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "usage: releasepack -version vX.Y.Z -out <empty dir> [-repo <root>]")
		os.Exit(2)
	}
	if err := run(*version, *out, *repo); err != nil {
		fmt.Fprintln(os.Stderr, "releasepack:", err)
		os.Exit(1)
	}
}

func run(version, out, repo string) error {
	repo, err := filepath.Abs(repo)
	if err != nil {
		return err
	}
	if err := checkVersion(version); err != nil {
		return err
	}

	// Preconditions on the repository state.
	status, err := gitOutput(repo, "status", "--porcelain")
	if err != nil {
		return err
	}
	if status != "" {
		return fmt.Errorf("worktree/index not clean; refusing to build a release from it:\n%s", status)
	}
	head, err := gitOutput(repo, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	tagCommit, err := gitOutput(repo, "rev-parse", "--verify", version+"^{commit}")
	if err != nil {
		return fmt.Errorf("tag %s not found: %w", version, err)
	}
	if tagCommit != head {
		return fmt.Errorf("tag %s points at %s, but HEAD is %s", version, tagCommit, head)
	}

	// Toolchain pin (GOTOOLCHAIN=local in the build env prevents silent
	// toolchain switching, so the installed go must BE the pinned one).
	goVersion, err := commandOutput(repo, buildEnv("", ""), "go", "env", "GOVERSION")
	if err != nil {
		return err
	}
	if err := checkGoVersion(goVersion); err != nil {
		return err
	}

	// SOURCE_DATE_EPOCH: derived from the tagged commit, never inherited.
	epochStr, err := gitOutput(repo, "log", "-1", "--format=%ct", "HEAD")
	if err != nil {
		return err
	}
	epochSec, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return fmt.Errorf("commit timestamp %q: %w", epochStr, err)
	}
	epoch := time.Unix(epochSec, 0).UTC()

	// Output directory must exist empty (or not exist yet).
	if err := os.MkdirAll(out, dirMode); err != nil {
		return err
	}
	if entries, err := os.ReadDir(out); err != nil {
		return err
	} else if len(entries) > 0 {
		return fmt.Errorf("output directory %s is not empty", out)
	}

	stage, err := os.MkdirTemp("", "releasepack-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(stage) }()

	ldflags := "-s -w -buildid= -X github.com/dimchansky/h3-go/internal/cli.buildVersion=" + version
	hostVerified := false
	var archives []string
	for _, t := range targets {
		name := fmt.Sprintf("h3-%s-%s-%s", version, t.goos, t.goarch)
		suffix := ""
		if t.goos == "windows" {
			suffix = ".exe"
		}
		bin := filepath.Join(stage, name+suffix)
		cmd := exec.Command("go", "build", "-trimpath", "-ldflags", ldflags, "-o", bin, "./cmd/h3")
		cmd.Dir = repo
		cmd.Env = buildEnv(t.goos, t.goarch)
		if outB, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build %s/%s: %w\n%s", t.goos, t.goarch, err, outB)
		}

		// Postcondition: embedded build info is exactly what the release
		// claims — module path, tagged commit, unmodified VCS state, the
		// target platform, and no host paths baked in.
		info, err := commandOutput(repo, buildEnv("", ""), "go", "version", "-m", bin)
		if err != nil {
			return err
		}
		if err := checkBuildInfo(info, head, t.goos, t.goarch, repo); err != nil {
			return fmt.Errorf("%s/%s: %w", t.goos, t.goarch, err)
		}

		// Postcondition: the injected CLI version, proven by execution
		// (not build-info inspection) on the host-runnable target.
		if t.goos == runtime.GOOS && t.goarch == runtime.GOARCH {
			verOut, err := commandOutput(repo, nil, bin, "--version")
			if err != nil {
				return fmt.Errorf("%s --version: %w", bin, err)
			}
			want := fmt.Sprintf("h3 4.4.0 (%s)", version)
			if verOut != want {
				return fmt.Errorf("%s --version = %q, want %q", bin, verOut, want)
			}
			hostVerified = true
		}

		files := []stagedFile{
			{name: "h3" + suffix, src: bin, mode: binMode},
			{name: "LICENSE", src: filepath.Join(repo, "LICENSE"), mode: docMode},
			{name: "NOTICE", src: filepath.Join(repo, "NOTICE"), mode: docMode},
			{name: "README.md", src: filepath.Join(repo, "cmd", "h3", "README-archive.md"), mode: docMode},
		}
		var archive string
		if t.goos == "windows" {
			archive = name + ".zip"
			err = writeZip(filepath.Join(out, archive), name, files, epoch)
		} else {
			archive = name + ".tar.gz"
			err = writeTarGz(filepath.Join(out, archive), name, files, epoch)
		}
		if err != nil {
			return fmt.Errorf("archive %s: %w", archive, err)
		}
		archives = append(archives, archive)
		fmt.Printf("packed %s\n", archive)
	}
	if !hostVerified {
		return fmt.Errorf("no target matched the host (%s/%s); the injected CLI version was never execution-verified", runtime.GOOS, runtime.GOARCH)
	}

	if err := writeSHA256SUMS(out, archives); err != nil {
		return err
	}
	fmt.Printf("releasepack: %d archives + SHA256SUMS in %s (SOURCE_DATE_EPOCH=%d, %s)\n",
		len(archives), out, epochSec, requiredGoVersion)
	return nil
}

// checkVersion enforces the canonical tag syntax.
func checkVersion(version string) error {
	if !versionRE.MatchString(version) {
		return fmt.Errorf("version %q violates the tag policy in docs/versioning.md (want vX.Y.Z or vX.Y.Z-rc.N)", version)
	}
	return nil
}

// checkGoVersion enforces the exact toolchain pin.
func checkGoVersion(goVersion string) error {
	if goVersion != requiredGoVersion {
		return fmt.Errorf("go toolchain is %s, but releases are pinned to %s (see docs/releasing.md for the bump procedure)", goVersion, requiredGoVersion)
	}
	return nil
}

// checkBuildInfo validates `go version -m` output for one built binary.
func checkBuildInfo(info, commit, goos, goarch, repo string) error {
	checks := []struct{ want, what string }{
		{"\tmod\tgithub.com/dimchansky/h3-go\t", "module path"},
		{"vcs.revision=" + commit, "tagged commit"},
		{"vcs.modified=false", "clean VCS state"},
		{"GOOS=" + goos, "target GOOS"},
		{"GOARCH=" + goarch, "target GOARCH"},
	}
	for _, c := range checks {
		if !strings.Contains(info, c.want) {
			return fmt.Errorf("build info lacks %s (%q):\n%s", c.what, c.want, info)
		}
	}
	if strings.Contains(info, repo) {
		return fmt.Errorf("build info leaks the repository path %s (is -trimpath in effect?):\n%s", repo, info)
	}
	return nil
}

// buildEnv returns the normalized build environment: host settings that
// could alter the binaries are explicitly overridden, never inherited.
// Empty goos/goarch yields the host-target variant used for `go env` and
// `go version -m` invocations.
func buildEnv(goos, goarch string) []string {
	overrides := map[string]string{
		"CGO_ENABLED":  "0",
		"GOFLAGS":      "",
		"GOEXPERIMENT": "",
		"GOTOOLCHAIN":  "local",
		"GOENV":        "off",
		"GOWORK":       "off",
		"TZ":           "UTC",
		"LC_ALL":       "C",
		"GOAMD64":      "v1",
		"GOARM64":      "v8.0",
	}
	if goos != "" {
		overrides["GOOS"] = goos
	}
	if goarch != "" {
		overrides["GOARCH"] = goarch
	}
	env := os.Environ()
	out := env[:0]
	for _, kv := range env {
		key, _, _ := strings.Cut(kv, "=")
		if _, replaced := overrides[key]; !replaced {
			out = append(out, kv)
		}
	}
	for key, val := range overrides {
		out = append(out, key+"="+val)
	}
	return out
}

// gitOutput runs git in the repository and returns trimmed stdout.
func gitOutput(repo string, args ...string) (string, error) {
	return commandOutput(repo, nil, "git", args...)
}

// commandOutput runs a command and returns trimmed stdout (stderr is
// attached to the error on failure).
func commandOutput(dir string, env []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = env
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, stderr.String())
	}
	return strings.TrimSpace(string(out)), nil
}
