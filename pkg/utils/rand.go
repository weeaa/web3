package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var ExecPath = getExecPath()

// getExecPath returns the executable's absolute path.
func getExecPath() string {
	ex, err := os.Executable()
	if err != nil {
		slog.Error("error getting exec path", err)
	}
	return filepath.Dir(ex)
}

func FirstLastFour(input string) string {
	return fmt.Sprintf("%s...%s", input[:4], input[len(input)-4:])
}
