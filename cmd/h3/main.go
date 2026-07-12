// Command h3 provides the upstream-compatible H3 command-line utility.
package main

import (
	"os"

	"github.com/dimchansky/h3-go/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
