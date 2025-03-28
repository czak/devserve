package main

import (
	"os"
	"path/filepath"
)

// Get local file path of the requested file.
// If request is for a directory, append "index.html".
//
// Note the path is not guaranteed to be for an existing file.
func resolveFilepath(dir, path string) string {
	fp := filepath.Join(dir, path)

	if path[len(path)-1] == '/' {
		ip := filepath.Join(fp, "index.html")
		if _, err := os.Stat(ip); err == nil {
			fp = ip
		}
	}

	return fp
}
