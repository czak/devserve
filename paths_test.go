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
	// ├── missing
	// └── present
	//     └── index.html
	dir, _ := os.OpenRoot(t.TempDir())
	dir.Mkdir("missing", 0755)
	dir.Mkdir("present", 0755)
	dir.Create("present/index.html")

	want := dir.Name() + "/missing"
	got := resolveFilepath(dir.Name(), "/missing/")
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	want = dir.Name() + "/present/index.html"
	got = resolveFilepath(dir.Name(), "/present/")
	if got != want {
		t.Errorf("want %v, got %v", want, got)
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

	got := resolveFilepath(dir.Name(), "/about")
	want := dir.Name() + "/about.html"
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	got = resolveFilepath(dir.Name(), "/fileconflict")
	want = dir.Name() + "/fileconflict"
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	got = resolveFilepath(dir.Name(), "/posts")
	want = dir.Name() + "/posts.html"
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	got = resolveFilepath(dir.Name(), "/posts/")
	want = dir.Name() + "/posts"
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
