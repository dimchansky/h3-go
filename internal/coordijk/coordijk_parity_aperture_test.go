//go:build oracle

package coordijk

import (
    "os/exec"
    "path/filepath"
    "runtime"
    "strconv"
    "strings"
    "testing"
)

func locateH3RefAperture(t *testing.T) (string, bool) {
    _, thisFile, _, _ := runtime.Caller(0)
    root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
    cand := filepath.Join(root, "testref", "h3ref")
    if _, err := exec.LookPath(cand); err == nil { return cand, true }
    return "", false
}

func TestOracle_ApertureTransforms(t *testing.T) {
    h3ref, ok := locateH3RefAperture(t)
    if !ok { t.Skip("testref/h3ref not found; run 'make ref' to build oracle") }
    inputs := []CoordIJK{{0,0,0},{1,0,0},{0,1,0},{0,0,1},{2,1,0},{7,0,0},{0,7,0}}
    run := func(args ...string) CoordIJK {
        b, err := exec.Command(h3ref, args...).Output()
        if err != nil { t.Fatalf("oracle error: %v", err) }
        f := strings.Fields(string(b)); if len(f) != 3 { t.Fatalf("unexpected oracle output: %q", string(b)) }
        i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
        return CoordIJK{i,j,k}
    }
    for _, v := range inputs {
        got := v; got.UpAp7();  want := run("coordijk_up_ap7", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)); if got != want { t.Fatalf("UpAp7(%v)=%v, want %v", v, got, want) }
        got = v; got.UpAp7r(); want = run("coordijk_up_ap7r", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)); if got != want { t.Fatalf("UpAp7r(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp7();  want = run("coordijk_down_ap7", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)); if got != want { t.Fatalf("DownAp7(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp7r(); want = run("coordijk_down_ap7r", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)); if got != want { t.Fatalf("DownAp7r(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp3();  want = run("coordijk_down_ap3", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)); if got != want { t.Fatalf("DownAp3(%v)=%v, want %v", v, got, want) }
        got = v; got.DownAp3r(); want = run("coordijk_down_ap3r", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)); if got != want { t.Fatalf("DownAp3r(%v)=%v, want %v", v, got, want) }
    }
}

