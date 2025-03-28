package main

import (
	"os"
	"path/filepath"
)

// Get local file path of the requested file.
//
// If request is for a directory, try "index.html" inside it.
// If request is for a missing file, try with ".html".
//
// Note the returned path is not guaranteed to be for an existing
// file or directory.
func resolveFilepath(dir, path string) string {
	fp := filepath.Join(dir, path)

	// Try index.html for directory requests (only if trailing '/'!)
	if path[len(path)-1] == '/' {
		index := filepath.Join(fp, "index.html")
		return ifExists(index, fp)
	}

	// Try .html for missing files (aka "pretty" urls)
	if info, err := os.Stat(fp); err != nil || info.IsDir() {
		pretty := fp + ".html"
		return ifExists(pretty, fp)
	}

	return fp
}

func ifExists(path, fallback string) string {
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return fallback
}
