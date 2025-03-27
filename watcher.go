package main

import (
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

func watchFiles(dir string, ps *pubsub) {
	w := watcher.New()
	defer w.Close()

	if err := w.AddRecursive(dir); err != nil {
		log.Fatal(err)
	}

	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					ps.publish(event.Name())
				}

			case err := <-w.Error:
				log.Fatal(err)
			}
		}
	}()

	w.Start(100 * time.Millisecond)
}
