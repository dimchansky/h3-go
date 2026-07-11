//go:build cgo && c2go

// This file provides an isolated compilation unit for utility.c
// to avoid duplicate symbols when combined with other C modules.

#include "utility.c"