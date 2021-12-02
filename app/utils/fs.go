package utils

import (
	"path/filepath"
)

func GetAbsolutePath(relativePath string) string {
	return filepath.Join("/mnt/onboard/.adds/kobotracker", relativePath)
}
