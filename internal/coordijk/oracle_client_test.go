//go:build oracle

package coordijk

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// oracleClient wraps the testref/h3ref binary and exposes helpers mirroring
// the coordijk API to reduce boilerplate in parity tests.
type oracleClient struct {
	t    *testing.T
	path string
}

// newOracle locates the h3ref binary (respects ORACLE_PATH override) and
// returns a ready client. Skips the calling test if not found.
func newOracle(t *testing.T) *oracleClient {
	t.Helper()
	// Allow explicit override
	if p := os.Getenv("ORACLE_PATH"); p != "" {
		if _, err := exec.LookPath(p); err == nil {
			return &oracleClient{t: t, path: p}
		}
	}
	// Resolve default location: repoRoot/testref/h3ref
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	cand := filepath.Join(root, "testref", "h3ref")
	if _, err := exec.LookPath(cand); err == nil {
		return &oracleClient{t: t, path: cand}
	}
	t.Skip("testref/h3ref not found; run 'make ref' or set ORACLE_PATH")
	return nil
}

// execOut runs the oracle command and returns trimmed stdout.
func (o *oracleClient) execOut(args ...string) string {
	o.t.Helper()
	out, err := exec.Command(o.path, args...).Output()
	if err != nil {
		o.t.Fatalf("oracle error for %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

// Distance mirrors Distance(a,b CoordIJK) int.
func (o *oracleClient) Distance(a, b CoordIJK) int {
	out := o.execOut("coordijk_distance",
		strconv.Itoa(a.I), strconv.Itoa(a.J), strconv.Itoa(a.K),
		strconv.Itoa(b.I), strconv.Itoa(b.J), strconv.Itoa(b.K),
	)
	v, err := strconv.Atoi(out)
	if err != nil {
		o.t.Fatalf("parse distance: %q: %v", out, err)
	}
	return v
}

// Rotate60CCW mirrors (*CoordIJK).Rotate60CCW but via oracle.
func (o *oracleClient) Rotate60CCW(v CoordIJK) CoordIJK {
	out := o.execOut("coordijk_rotate", "ccw",
		strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K))
	f := strings.Fields(out)
	if len(f) != 3 {
		o.t.Fatalf("unexpected rotate ccw output: %q", out)
	}
	i, _ := strconv.Atoi(f[0])
	j, _ := strconv.Atoi(f[1])
	k, _ := strconv.Atoi(f[2])
	return CoordIJK{i, j, k}
}

// Rotate60CW mirrors (*CoordIJK).Rotate60CW but via oracle.
func (o *oracleClient) Rotate60CW(v CoordIJK) CoordIJK {
	out := o.execOut("coordijk_rotate", "cw",
		strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K))
	f := strings.Fields(out)
	if len(f) != 3 {
		o.t.Fatalf("unexpected rotate cw output: %q", out)
	}
	i, _ := strconv.Atoi(f[0])
	j, _ := strconv.Atoi(f[1])
	k, _ := strconv.Atoi(f[2])
	return CoordIJK{i, j, k}
}

// Neighbor mirrors (*CoordIJK).Neighbor.
func (o *oracleClient) Neighbor(v CoordIJK, d Direction) CoordIJK {
	out := o.execOut("coordijk_neighbor", strconv.Itoa(int(d)),
		strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K))
	f := strings.Fields(out)
	if len(f) != 3 {
		o.t.Fatalf("unexpected neighbor output: %q", out)
	}
	i, _ := strconv.Atoi(f[0])
	j, _ := strconv.Atoi(f[1])
	k, _ := strconv.Atoi(f[2])
	return CoordIJK{i, j, k}
}

// Up/Down aperture transforms.
func (o *oracleClient) UpAp7(v CoordIJK) CoordIJK    { return o.ijkCmd("coordijk_up_ap7", v) }
func (o *oracleClient) UpAp7r(v CoordIJK) CoordIJK   { return o.ijkCmd("coordijk_up_ap7r", v) }
func (o *oracleClient) DownAp7(v CoordIJK) CoordIJK  { return o.ijkCmd("coordijk_down_ap7", v) }
func (o *oracleClient) DownAp7r(v CoordIJK) CoordIJK { return o.ijkCmd("coordijk_down_ap7r", v) }
func (o *oracleClient) DownAp3(v CoordIJK) CoordIJK  { return o.ijkCmd("coordijk_down_ap3", v) }
func (o *oracleClient) DownAp3r(v CoordIJK) CoordIJK { return o.ijkCmd("coordijk_down_ap3r", v) }

func (o *oracleClient) ijkCmd(cmd string, v CoordIJK) CoordIJK {
	out := o.execOut(cmd, strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K))
	f := strings.Fields(out)
	if len(f) != 3 {
		o.t.Fatalf("unexpected %s output: %q", cmd, out)
	}
	i, _ := strconv.Atoi(f[0])
	j, _ := strconv.Atoi(f[1])
	k, _ := strconv.Atoi(f[2])
	return CoordIJK{i, j, k}
}

// Hex2d conversions.
func (o *oracleClient) IJKToHex2D(v CoordIJK) Vec2d {
	out := o.execOut("coordijk_hex2d",
		strconv.Itoa(v.I), strconv.Itoa(v.J), strconv.Itoa(v.K))
	f := strings.Fields(out)
	if len(f) != 2 {
		o.t.Fatalf("unexpected coordijk_hex2d output: %q", out)
	}
	x, err1 := strconv.ParseFloat(f[0], 64)
	y, err2 := strconv.ParseFloat(f[1], 64)
	if err1 != nil || err2 != nil {
		o.t.Fatalf("parse hex2d: %q", out)
	}
	return Vec2d{X: x, Y: y}
}

func (o *oracleClient) Hex2DToIJK(x, y float64) CoordIJK {
	out := o.execOut("coordijk_from_hex2d", strconv.FormatFloat(x, 'g', -1, 64), strconv.FormatFloat(y, 'g', -1, 64))
	f := strings.Fields(out)
	if len(f) != 3 {
		o.t.Fatalf("unexpected coordijk_from_hex2d output: %q", out)
	}
	i, _ := strconv.Atoi(f[0])
	j, _ := strconv.Atoi(f[1])
	k, _ := strconv.Atoi(f[2])
	return CoordIJK{i, j, k}
}
