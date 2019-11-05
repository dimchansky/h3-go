// +build tools

package tools

import (
	_ "golang.org/x/lint/golint"           // tool
	_ "golang.org/x/tools/cmd/goimports"   // tool
	_ "honnef.co/go/tools/cmd/staticcheck" // tool
)