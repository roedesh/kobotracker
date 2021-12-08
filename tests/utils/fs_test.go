package main

import (
	"kobotracker/app/utils"
	"testing"
)

func TestGetAbsolutePath(t *testing.T) {
	path := utils.GetAbsolutePath("path/to/file.txt")
	if path != "/mnt/onboard/.adds/kobotracker/path/to/file.txt" {
		t.Errorf("Path was incorrect, got: %s, want: /mnt/onboard/.adds/kobotracker/path/to/file.txt.", path)
	}
}
