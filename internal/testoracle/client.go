//go:build oracle

package testoracle

import (
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strconv"
    "strings"
    "testing"
)

// Client is a thin wrapper over testref/h3ref for oracle-backed tests.
type Client struct {
    t    *testing.T
    path string
}

// New returns a client, resolving h3ref path from ORACLE_PATH or defaulting
// to repo-root/testref/h3ref. Skips the test if not found.
func New(t *testing.T) *Client {
    t.Helper()
    if p := os.Getenv("ORACLE_PATH"); p != "" {
        if _, err := exec.LookPath(p); err == nil {
            return &Client{t: t, path: p}
        }
    }
    _, thisFile, _, _ := runtime.Caller(0)
    root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
    cand := filepath.Join(root, "testref", "h3ref")
    if _, err := exec.LookPath(cand); err == nil {
        return &Client{t: t, path: cand}
    }
    t.Skip("testref/h3ref not found; run 'make ref' or set ORACLE_PATH")
    return nil
}

func (c *Client) out(args ...string) string {
    c.t.Helper()
    b, err := exec.Command(c.path, args...).Output()
    if err != nil {
        c.t.Fatalf("oracle error for %v: %v", args, err)
    }
    return strings.TrimSpace(string(b))
}

// DistanceIJK mirrors Distance on IJK triples.
func (c *Client) DistanceIJK(a, b [3]int) int {
    out := c.out("coordijk_distance",
        strconv.Itoa(a[0]), strconv.Itoa(a[1]), strconv.Itoa(a[2]),
        strconv.Itoa(b[0]), strconv.Itoa(b[1]), strconv.Itoa(b[2]))
    v, err := strconv.Atoi(out)
    if err != nil { c.t.Fatalf("parse distance: %q: %v", out, err) }
    return v
}

func (c *Client) RotateIJKCCW(v [3]int) [3]int {
    out := c.out("coordijk_rotate", "ccw",
        strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected rotate ccw output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) RotateIJKCW(v [3]int) [3]int {
    out := c.out("coordijk_rotate", "cw",
        strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected rotate cw output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) NeighborIJK(v [3]int, d int) [3]int {
    out := c.out("coordijk_neighbor", strconv.Itoa(d),
        strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected neighbor output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) UpAp7IJK(v [3]int) [3]int    { return c.ijkCmd("coordijk_up_ap7", v) }
func (c *Client) UpAp7rIJK(v [3]int) [3]int   { return c.ijkCmd("coordijk_up_ap7r", v) }
func (c *Client) DownAp7IJK(v [3]int) [3]int  { return c.ijkCmd("coordijk_down_ap7", v) }
func (c *Client) DownAp7rIJK(v [3]int) [3]int { return c.ijkCmd("coordijk_down_ap7r", v) }
func (c *Client) DownAp3IJK(v [3]int) [3]int  { return c.ijkCmd("coordijk_down_ap3", v) }
func (c *Client) DownAp3rIJK(v [3]int) [3]int { return c.ijkCmd("coordijk_down_ap3r", v) }

func (c *Client) ijkCmd(cmd string, v [3]int) [3]int {
    out := c.out(cmd, strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected %s output: %q", cmd, out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) IJKToHex2D(v [3]int) (float64, float64) {
    out := c.out("coordijk_hex2d", strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 2 { c.t.Fatalf("unexpected coordijk_hex2d output: %q", out) }
    x, err1 := strconv.ParseFloat(f[0], 64)
    y, err2 := strconv.ParseFloat(f[1], 64)
    if err1 != nil || err2 != nil { c.t.Fatalf("parse hex2d: %q", out) }
    return x, y
}

func (c *Client) Hex2DToIJK(x, y float64) [3]int {
    out := c.out("coordijk_from_hex2d", strconv.FormatFloat(x, 'g', -1, 64), strconv.FormatFloat(y, 'g', -1, 64))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected coordijk_from_hex2d output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

