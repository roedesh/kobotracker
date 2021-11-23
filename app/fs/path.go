// Package fs provides a helper function for working with paths relative
// to the executable.
package fs

import (
	"os"
	"path/filepath"
)

var (
	rootDir string // The directory in which the executable resides
)

func init() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootDir = filepath.Dir(executable)
}

// Given a path relative to the executable, returns an absolute path.
func GetAbsolutePath(relativePath string) string {
	return filepath.Join(rootDir, relativePath)
}
