package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
)

func watchFiles(dir string, ps *pubsub) {
	w := watcher.New()
	defer w.Close()

	basepath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	if err := w.AddRecursive(basepath); err != nil {
		log.Fatal(err)
	}

	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.IsDir() {
					continue
				}

				relpath, err := filepath.Rel(basepath, event.Path)
				if err != nil {
					log.Fatal(err)
				}

				ps.publish(relpath)

			case err := <-w.Error:
				log.Fatal(err)
			}
		}
	}()

	w.Start(100 * time.Millisecond)
}
