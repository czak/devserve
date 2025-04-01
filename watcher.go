package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"time"

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
		logger.Debug("Watching %s", dir)
		watcher.Add(dir)
	})

	debounce := debouncer{
		duration: 50 * time.Millisecond,
	}

	for {
		select {
		case event := <-watcher.Events:
			if !event.Has(fsnotify.Write) {
				continue
			}

			relpath, err := filepath.Rel(basepath, event.Name)
			if err != nil {
				log.Fatal(err)
			}

			debounce.then(func() {
				logger.Debug("Change event: %s", relpath)
				ps.publish(relpath)
			})

		case err := <-watcher.Errors:
			logger.Error("Watcher error: %v", err)
		}
	}
}

func walkDirsRecursive(root string, dirfn func(string)) {
	walk(root, root, dirfn)
}

func walk(root string, sym string, dirfn func(string)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Warn("Unable to enter %s: %v", path, err)
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
