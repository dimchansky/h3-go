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

// locateH3Ref attempts to find the oracle binary built at testref/h3ref.
func locateH3Ref(t *testing.T) (string, bool) {
    _, thisFile, _, _ := runtime.Caller(0)
    // repo root = two levels up from this file's directory
    root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
    cand := filepath.Join(root, "testref", "h3ref")
    if _, err := exec.LookPath(cand); err == nil {
        return cand, true
    }
    return "", false
}

func TestOracle_CoordIJKDistance(t *testing.T) {
    h3ref, ok := locateH3Ref(t)
    if !ok {
        t.Skip("testref/h3ref not found; run 'make ref' to build oracle")
    }

    cases := []struct{ a, b CoordIJK }{
        {CoordIJK{0,0,0}, CoordIJK{0,0,0}},
        {CoordIJK{1,0,0}, CoordIJK{0,0,0}},
        {CoordIJK{0,2,0}, CoordIJK{0,0,0}},
        {CoordIJK{1,0,1}, CoordIJK{0,0,0}},
        {CoordIJK{1,0,1}, CoordIJK{0,2,0}},
        {CoordIJK{2,-1,1}, CoordIJK{-1,2,-1}},
    }
    for _, tc := range cases {
        got := Distance(tc.a, tc.b)
        cmd := exec.Command(h3ref, "coordijk_distance",
            strconv.Itoa(tc.a.I), strconv.Itoa(tc.a.J), strconv.Itoa(tc.a.K),
            strconv.Itoa(tc.b.I), strconv.Itoa(tc.b.J), strconv.Itoa(tc.b.K))
        out, err := cmd.Output()
        if err != nil { t.Fatalf("oracle error: %v", err) }
        want, _ := strconv.Atoi(strings.TrimSpace(string(out)))
        if got != want {
            t.Fatalf("Distance(%v,%v) = %d, want %d", tc.a, tc.b, got, want)
        }
    }
}

func TestOracle_CoordIJKRotate(t *testing.T) {
    h3ref, ok := locateH3Ref(t)
    if !ok { t.Skip("testref/h3ref not found; run 'make ref' to build oracle") }
    cases := []CoordIJK{{0,0,0},{1,0,0},{0,1,0},{0,0,1},{1,1,0},{1,0,1},{0,1,1},{2,-1,1},{-1,2,-1}}
    for _, v := range cases {
        got := v; got.Rotate60CCW()
        out, err := exec.Command(h3ref, "coordijk_rotate", "ccw", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)).Output()
        if err != nil { t.Fatalf("oracle error: %v", err) }
        f := strings.Fields(string(out)); wi, _ := strconv.Atoi(f[0]); wj, _ := strconv.Atoi(f[1]); wk, _ := strconv.Atoi(f[2])
        if got != (CoordIJK{wi,wj,wk}) { t.Fatalf("Rotate60CCW(%v) = %v, want %v", v, got, CoordIJK{wi,wj,wk}) }

        got = v; got.Rotate60CW()
        out, err = exec.Command(h3ref, "coordijk_rotate", "cw", strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)).Output()
        if err != nil { t.Fatalf("oracle error: %v", err) }
        f = strings.Fields(string(out)); wi, _ = strconv.Atoi(f[0]); wj, _ = strconv.Atoi(f[1]); wk, _ = strconv.Atoi(f[2])
        if got != (CoordIJK{wi,wj,wk}) { t.Fatalf("Rotate60CW(%v) = %v, want %v", v, got, CoordIJK{wi,wj,wk}) }
    }
}

func TestOracle_CoordIJKNeighbor(t *testing.T) {
    h3ref, ok := locateH3Ref(t)
    if !ok { t.Skip("testref/h3ref not found; run 'make ref' to build oracle") }
    v := CoordIJK{2,1,-1}
    for d := Direction(0); d < NumDigits; d++ {
        got := v; got.Neighbor(d)
        out, err := exec.Command(h3ref, "coordijk_neighbor", strconv.Itoa(int(d)), strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K)).Output()
        if err != nil { t.Fatalf("oracle error: %v", err) }
        f := strings.Fields(string(out)); wi, _ := strconv.Atoi(f[0]); wj, _ := strconv.Atoi(f[1]); wk, _ := strconv.Atoi(f[2])
        if got != (CoordIJK{wi,wj,wk}) { t.Fatalf("Neighbor(%d) = %v, want %v", d, got, CoordIJK{wi,wj,wk}) }
    }
}

