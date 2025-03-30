package main

import (
	"io/fs"
	"log"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func watchFiles(root string, ps *pubsub) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	basepath, err := filepath.Abs(root)
	if err != nil {
		log.Fatal(err)
	}

	walkDirsRecursive(basepath, func(dir string) {
		slog.Debug("Watching", "dir", dir)
		watcher.Add(dir)
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if !event.Has(fsnotify.Write) {
				continue
			}

			relpath, err := filepath.Rel(basepath, event.Name)
			if err != nil {
				log.Fatal(err)
			}

			slog.Debug("Change event", "path", relpath)

			ps.publish(relpath)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("Watcher error", "error", err)
		}
	}
}

func walkDirsRecursive(root string, dirfn func(string)) {
	walk(root, root, dirfn)
}

func walk(root string, sym string, dirfn func(string)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Warn("Unable to enter", "path", path)
			return filepath.SkipDir
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		path = filepath.Join(sym, rel)

		if d.IsDir() {
			if isHidden(path) {
				return filepath.SkipDir
			}

			dirfn(path)
		}

		if d.Type()&fs.ModeSymlink == fs.ModeSymlink {
			realpath, _ := filepath.EvalSymlinks(path)
			return walk(realpath, path, dirfn)
		}

		return nil
	})
}

func isHidden(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}
