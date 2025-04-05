package version

import (
	_ "embed"

	"fmt"
	"strings"
)

var (
	//go:embed version.txt
	version string
)

// Print returns current CLI tool version.
func Print() string {
	return fmt.Sprintf("v%s", strings.TrimSpace(version))
}
