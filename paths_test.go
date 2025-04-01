package main

import (
	"os"
	"testing"
)

func TestResolveFilepathRegular(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		path string
		want string
	}{
		{
			"simple",
			"basedir",
			"path/to.css",
			"basedir/path/to.css",
		},
		{
			"path with leading slash",
			"basedir",
			"/path/to.css",
			"basedir/path/to.css",
		},
		{
			"traverse up",
			"basedir",
			"../to.css",
			"to.css",
		},
		{
			"absolute base",
			"/tmp/basedir",
			"../to.css",
			"/tmp/to.css",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFilepath(tt.dir, tt.path)
			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestResolveFilepathIndex(t *testing.T) {
	// dir
	// ├── index.html
	// ├── missing
	// └── present
	//     └── index.html
	dir, _ := os.OpenRoot(t.TempDir())
	dir.Create("index.html")
	dir.Mkdir("missing", 0755)
	dir.Mkdir("present", 0755)
	dir.Create("present/index.html")

	tests := []struct {
		name string
		dir  string
		path string
		want string
	}{
		{
			name: "root with index.html",
			dir:  dir.Name(),
			path: "/",
			want: dir.Name() + "/index.html",
		},
		{
			name: "directory with no index",
			dir:  dir.Name(),
			path: "/missing/",
			want: dir.Name() + "/missing",
		},
		{
			name: "directory with index.html",
			dir:  dir.Name(),
			path: "/present/",
			want: dir.Name() + "/present/index.html",
		},
		{
			name: "directory with index.html, no trailing slash",
			dir:  dir.Name(),
			path: "/present",
			want: dir.Name() + "/present",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFilepath(tt.dir, tt.path)
			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestResolveFilepathPretty(t *testing.T) {
	// dir
	// ├── about.html
	// ├── fileconflict
	// ├── fileconflict.html
	// ├── posts/
	// │   └── ...
	// └── posts.html
	dir, _ := os.OpenRoot(t.TempDir())
	dir.Create("about.html")
	dir.Create("fileconflict")
	dir.Create("fileconflict.html")
	dir.Mkdir("posts", 0755)
	dir.Create("posts.html")

	tests := []struct {
		name string
		dir  string
		path string
		want string
	}{
		{
			name: "pretty url",
			dir:  dir.Name(),
			path: "/about",
			want: dir.Name() + "/about.html",
		},
		{
			name: "conflict with actual file",
			dir:  dir.Name(),
			path: "/fileconflict",
			want: dir.Name() + "/fileconflict",
		},
		{
			name: "same name as sibling directory",
			dir:  dir.Name(),
			path: "/posts",
			want: dir.Name() + "/posts.html",
		},
		{
			name: "directory with trailing slash",
			dir:  dir.Name(),
			path: "/posts/",
			want: dir.Name() + "/posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFilepath(tt.dir, tt.path)
			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}
