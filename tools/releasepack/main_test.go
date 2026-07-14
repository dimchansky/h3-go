package main

import (
	"strings"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	accept := []string{"v0.3.0", "v0.3.0-rc.1", "v0.3.0-rc.10", "v1.2.3", "v10.20.30"}
	reject := []string{
		"",
		"0.3.0",          // missing v
		"v0.4.4.0",       // four components (H3-axis encoding is forbidden)
		"v1garbage",      // not SemVer
		"v01.2.3",        // leading zero
		"v0.3.0-rc.0",    // rc counting starts at 1
		"v0.3.0-rc.01",   // leading zero in rc number
		"v0.3.0-alpha.1", // only the documented -rc.N prerelease form
		"v0.3.0+h3.4.4",  // build metadata rejected by project policy
		"v0.3.0-h3.4.5.0",
		"v0.3.0 ",
	}
	for _, v := range accept {
		if err := checkVersion(v); err != nil {
			t.Errorf("checkVersion(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range reject {
		if err := checkVersion(v); err == nil {
			t.Errorf("checkVersion(%q) = nil, want error", v)
		}
	}
}

func TestCheckGoVersion(t *testing.T) {
	if err := checkGoVersion(requiredGoVersion); err != nil {
		t.Errorf("exact pin rejected: %v", err)
	}
	for _, v := range []string{"go1.24.11", "go1.27.0", "", "devel +abc123"} {
		if err := checkGoVersion(v); err == nil {
			t.Errorf("checkGoVersion(%q) = nil, want error", v)
		}
	}
}

// buildInfo assembles a plausible `go version -m` output.
func buildInfo(commit, goos, goarch string, modified bool, extra string) string {
	mod := "false"
	if modified {
		mod = "true"
	}
	return strings.Join([]string{
		"bin/h3: go1.26.5",
		"\tpath\tgithub.com/dimchansky/h3-go/cmd/h3",
		"\tmod\tgithub.com/dimchansky/h3-go\t(devel)",
		"\tbuild\t-buildmode=exe",
		"\tbuild\t-trimpath=true",
		"\tbuild\tCGO_ENABLED=0",
		"\tbuild\tGOOS=" + goos,
		"\tbuild\tGOARCH=" + goarch,
		"\tbuild\tvcs=git",
		"\tbuild\tvcs.revision=" + commit,
		"\tbuild\tvcs.modified=" + mod,
		extra,
	}, "\n")
}

func TestCheckBuildInfo(t *testing.T) {
	const commit = "0123456789abcdef0123456789abcdef01234567"
	const repo = "/home/user/src/h3-go"

	if err := checkBuildInfo(buildInfo(commit, "linux", "amd64", false, ""), commit, "linux", "amd64", repo); err != nil {
		t.Errorf("valid build info rejected: %v", err)
	}

	bad := []struct {
		name string
		info string
	}{
		{"wrong commit", buildInfo(strings.Repeat("f", 40), "linux", "amd64", false, "")},
		{"modified vcs state", buildInfo(commit, "linux", "amd64", true, "")},
		{"wrong GOOS", buildInfo(commit, "windows", "amd64", false, "")},
		{"wrong GOARCH", buildInfo(commit, "linux", "arm64", false, "")},
		{"leaked repo path", buildInfo(commit, "linux", "amd64", false, "\tbuild\t-ldflags="+repo)},
		{"wrong module path", strings.ReplaceAll(buildInfo(commit, "linux", "amd64", false, ""), "dimchansky/h3-go", "other/module")},
	}
	for _, tc := range bad {
		if err := checkBuildInfo(tc.info, commit, "linux", "amd64", repo); err == nil {
			t.Errorf("%s: checkBuildInfo = nil, want error", tc.name)
		}
	}
}

func TestBuildEnvNormalization(t *testing.T) {
	t.Setenv("GOFLAGS", "-tags=evil")
	t.Setenv("GOAMD64", "v4")
	t.Setenv("GOEXPERIMENT", "loopvar")
	t.Setenv("GOWORK", "/tmp/go.work")

	env := buildEnv("linux", "arm64")
	got := map[string]string{}
	for _, kv := range env {
		k, v, _ := strings.Cut(kv, "=")
		if _, dup := got[k]; dup {
			t.Errorf("duplicate env key %q", k)
		}
		got[k] = v
	}
	want := map[string]string{
		"GOOS": "linux", "GOARCH": "arm64",
		"CGO_ENABLED": "0", "GOFLAGS": "", "GOEXPERIMENT": "",
		"GOTOOLCHAIN": "local", "GOENV": "off", "GOWORK": "off",
		"TZ": "UTC", "LC_ALL": "C",
		"GOAMD64": "v1", "GOARM64": "v8.0",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("env %s = %q, want %q", k, got[k], v)
		}
	}

	hostEnv := buildEnv("", "")
	for _, kv := range hostEnv {
		if strings.HasPrefix(kv, "GOOS=") || strings.HasPrefix(kv, "GOARCH=") {
			t.Errorf("host-variant env must not force GOOS/GOARCH, found %q", kv)
		}
	}
}
