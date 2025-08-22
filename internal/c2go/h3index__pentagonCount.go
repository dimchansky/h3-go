package c2go

import "github.com/dimchansky/h3-go/internal/tables"

// pentagonCount returns the number of pentagon base cells (parity with C).
func pentagonCount() int { return tables.NumPentagons }
