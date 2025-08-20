//go:build oracle

package coordijk

import (
    "math/rand"
    "os"
    "strconv"
    "time"
)

func oracleMax() int {
    if v := os.Getenv("ORACLE_MAX"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            return n
        }
    }
    return 200
}

func oracleSeed() int64 {
    if v := os.Getenv("ORACLE_SEED"); v != "" {
        if n, err := strconv.ParseInt(v, 10, 64); err == nil {
            return n
        }
    }
    return 1337
}

func newRand() *rand.Rand {
    s := oracleSeed()
    if s == 0 {
        s = time.Now().UnixNano()
    }
    return rand.New(rand.NewSource(s))
}

func randInt(r *rand.Rand, min, max int) int {
    if min > max { min, max = max, min }
    return r.Intn(max-min+1) + min
}

func randIJK(r *rand.Rand, min, max int) CoordIJK {
    return CoordIJK{randInt(r, min, max), randInt(r, min, max), randInt(r, min, max)}
}

func randIJKPlus(r *rand.Rand, max int) CoordIJK {
    return CoordIJK{randInt(r, 0, max), randInt(r, 0, max), randInt(r, 0, max)}
}

