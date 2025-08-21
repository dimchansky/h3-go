//go:build oracle

// Package testoracle exposes a thin Go client around the external H3 C oracle
// binary built under testref/, enabling parity tests without cgo.
package testoracle

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Max returns ORACLE_MAX (default 200), controlling test sweep sizes.
func Max() int {
	if v := os.Getenv("ORACLE_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 200
}

// Seed returns ORACLE_SEED (default 1337). 0 means time-based.
func Seed() int64 {
	if v := os.Getenv("ORACLE_SEED"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return 1337
}

// NewRand constructs a rand.Rand using Seed (or time-based if Seed() == 0).
func NewRand() *rand.Rand {
	s := Seed()
	if s == 0 {
		s = time.Now().UnixNano()
	}
	return rand.New(rand.NewSource(s))
}

// RandIJK returns a random IJK triple in [min,max].
func RandIJK(r *rand.Rand, min, max int) [3]int {
	if min > max {
		min, max = max, min
	}
	return [3]int{r.Intn(max-min+1) + min, r.Intn(max-min+1) + min, r.Intn(max-min+1) + min}
}

// RandIJKPlus returns a random IJK+ triple in [0,max].
func RandIJKPlus(r *rand.Rand, max int) [3]int {
	return [3]int{r.Intn(max + 1), r.Intn(max + 1), r.Intn(max + 1)}
}
