package main

import (
	"os"
	"slices"
	"testing"
)

func TestWalkDirsRecursive(t *testing.T) {
	// root
	// ├── ext
	// │   └── ext2
	// └── start
	//     ├── dir
	//     │   └── dir2
	//     └── link -> ../ext
	root := t.TempDir()
	os.MkdirAll(root+"/ext/ext2", 0755)
	os.MkdirAll(root+"/start/dir/dir2", 0755)
	os.Symlink("../ext", root+"/start/link")

	start := root + "/start"

	dirs := []string{}
	expected := []string{
		start,
		start + "/dir",
		start + "/dir/dir2",
		start + "/link",
		start + "/link/ext2",
	}

	walkDirsRecursive(start, func(dir string) {
		dirs = append(dirs, dir)
	})

	if !slices.Equal(dirs, expected) {
		t.Errorf("expected %v, got %v", expected, dirs)
	}
}
