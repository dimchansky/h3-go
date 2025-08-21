//go:build oracle

package testoracle

import (
    "fmt"
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

// IJKDistance mirrors Distance on IJK triples.
func (c *Client) IJKDistance(a, b [3]int) int {
    out := c.out("ijkDistance",
        strconv.Itoa(a[0]), strconv.Itoa(a[1]), strconv.Itoa(a[2]),
        strconv.Itoa(b[0]), strconv.Itoa(b[1]), strconv.Itoa(b[2]))
    v, err := strconv.Atoi(out)
    if err != nil { c.t.Fatalf("parse distance: %q: %v", out, err) }
    return v
}

func (c *Client) IJKRotate60ccw(v [3]int) [3]int {
    out := c.out("ijkRotate60ccw",
        strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected rotate ccw output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) IJKRotate60cw(v [3]int) [3]int {
    out := c.out("ijkRotate60cw",
        strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected rotate cw output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) Neighbor(v [3]int, d int) [3]int {
    out := c.out("neighbor", strconv.Itoa(d),
        strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected neighbor output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) UpAp7(v [3]int) [3]int    { return c.ijkCmd("upAp7", v) }
func (c *Client) UpAp7r(v [3]int) [3]int   { return c.ijkCmd("upAp7r", v) }
func (c *Client) DownAp7(v [3]int) [3]int  { return c.ijkCmd("downAp7", v) }
func (c *Client) DownAp7r(v [3]int) [3]int { return c.ijkCmd("downAp7r", v) }
func (c *Client) DownAp3(v [3]int) [3]int  { return c.ijkCmd("downAp3", v) }
func (c *Client) DownAp3r(v [3]int) [3]int { return c.ijkCmd("downAp3r", v) }

func (c *Client) ijkCmd(cmd string, v [3]int) [3]int {
    out := c.out(cmd, strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected %s output: %q", cmd, out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

func (c *Client) IJKToHex2d(v [3]int) (float64, float64) {
    out := c.out("ijkToHex2d", strconv.Itoa(v[0]), strconv.Itoa(v[1]), strconv.Itoa(v[2]))
    f := strings.Fields(out)
    if len(f) != 2 { c.t.Fatalf("unexpected ijkToHex2d output: %q", out) }
    x, err1 := strconv.ParseFloat(f[0], 64)
    y, err2 := strconv.ParseFloat(f[1], 64)
    if err1 != nil || err2 != nil { c.t.Fatalf("parse hex2d: %q", out) }
    return x, y
}

func (c *Client) Hex2dToCoordIJK(x, y float64) [3]int {
    out := c.out("hex2dToCoordIJK", strconv.FormatFloat(x, 'g', -1, 64), strconv.FormatFloat(y, 'g', -1, 64))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected hex2dToCoordIJK output: %q", out) }
    i, _ := strconv.Atoi(f[0]); j, _ := strconv.Atoi(f[1]); k, _ := strconv.Atoi(f[2])
    return [3]int{i, j, k}
}

// Rotate an H3 index 60 degrees CCW via oracle.
func (c *Client) H3Rotate60ccw(h uint64) uint64 {
    out := c.out("h3Rotate60ccw", strconv.FormatUint(h, 16))
    v, err := strconv.ParseUint(strings.TrimPrefix(out, "0x"), 16, 64)
    if err != nil { c.t.Fatalf("parse rotate60ccw: %q: %v", out, err) }
    return v
}

// Rotate an H3 index 60 degrees CW via oracle.
func (c *Client) H3Rotate60cw(h uint64) uint64 {
    out := c.out("h3Rotate60cw", strconv.FormatUint(h, 16))
    v, err := strconv.ParseUint(strings.TrimPrefix(out, "0x"), 16, 64)
    if err != nil { c.t.Fatalf("parse rotate60cw: %q: %v", out, err) }
    return v
}

// GeoToFaceIJK returns (face,i,j,k) for given lat,lng (degrees) and res.
func (c *Client) GeoToFaceIjk(lat, lng float64, res int) (int, int, int, int) {
    out := c.out("geoToFaceIjk", strconv.FormatFloat(lat, 'g', -1, 64), strconv.FormatFloat(lng, 'g', -1, 64), strconv.Itoa(res))
    var face, i, j, k int
    _, err := fmt.Sscanf(out, "%d %d %d %d", &face, &i, &j, &k)
    if err != nil {
        c.t.Fatalf("parse geotofaceijk: %q: %v", out, err)
    }
    return face, i, j, k
}


// H3FromFaceIJK converts Face+IJK+res to an H3 index via oracle.
func (c *Client) FaceIjkToH3(face, i, j, k, res int) uint64 {
    out := c.out("faceIjkToH3", strconv.Itoa(face), strconv.Itoa(i), strconv.Itoa(j), strconv.Itoa(k), strconv.Itoa(res))
    // out like 0x8928308280fffff
    v, err := strconv.ParseUint(strings.TrimPrefix(out, "0x"), 16, 64)
    if err != nil { c.t.Fatalf("parse faceijk h3: %q: %v", out, err) }
    return v
}

// H3FromLatLng converts latitude/longitude (degrees) + res to an H3 index via oracle.
// Returns (h3, errCode) where errCode==0 on success.
func (c *Client) LatLngToCell(lat, lng float64, res int) (uint64, int) {
    out := c.out("latLngToCell",
        strconv.FormatFloat(lat, 'g', -1, 64),
        strconv.FormatFloat(lng, 'g', -1, 64),
        strconv.Itoa(res),
    )
    // out like: "0x8928308280fffff 0" or "0x0 <err>"
    f := strings.Fields(out)
    if len(f) != 2 { c.t.Fatalf("unexpected latlng output: %q", out) }
    h3, err := strconv.ParseUint(strings.TrimPrefix(f[0], "0x"), 16, 64)
    if err != nil { c.t.Fatalf("parse latlng h3: %q: %v", out, err) }
    code, err2 := strconv.Atoi(f[1])
    if err2 != nil { c.t.Fatalf("parse latlng errcode: %q: %v", out, err2) }
    return h3, code
}

// GetResolution returns the resolution for an H3 index via oracle.
func (c *Client) GetResolution(h uint64) int {
    out := c.out("getResolution", strconv.FormatUint(h, 16))
    v, err := strconv.Atoi(out)
    if err != nil { c.t.Fatalf("parse getResolution: %q: %v", out, err) }
    return v
}

// GetBaseCellNumber returns the base cell number for an H3 index via oracle.
func (c *Client) GetBaseCellNumber(h uint64) int {
    out := c.out("getBaseCellNumber", strconv.FormatUint(h, 16))
    v, err := strconv.Atoi(out)
    if err != nil { c.t.Fatalf("parse getBaseCellNumber: %q: %v", out, err) }
    return v
}

// IsBaseCellPentagon checks whether a base cell number is a pentagon.
func (c *Client) IsBaseCellPentagon(base int) bool {
    out := c.out("isBaseCellPentagon", strconv.Itoa(base))
    v, err := strconv.Atoi(out)
    if err != nil { c.t.Fatalf("parse isBaseCellPentagon: %q: %v", out, err) }
    return v != 0
}

// GeoToHex2d wraps _geoToHex2d and returns (face, x, y).
func (c *Client) GeoToHex2d(lat, lng float64, res int) (int, float64, float64) {
    out := c.out("geoToHex2d",
        strconv.FormatFloat(lat, 'g', -1, 64),
        strconv.FormatFloat(lng, 'g', -1, 64),
        strconv.Itoa(res))
    f := strings.Fields(out)
    if len(f) != 3 { c.t.Fatalf("unexpected geoToHex2d output: %q", out) }
    face, err1 := strconv.Atoi(f[0])
    x, err2 := strconv.ParseFloat(f[1], 64)
    y, err3 := strconv.ParseFloat(f[2], 64)
    if err1 != nil || err2 != nil || err3 != nil {
        c.t.Fatalf("parse geoToHex2d: %q: %v %v %v", out, err1, err2, err3)
    }
    return face, x, y
}
