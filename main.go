package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed event.js
var eventScript string

type config struct {
	addr  string
	dir   string
	debug bool
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.addr, "addr", ":8080", "network address")
	flag.StringVar(&cfg.dir, "dir", ".", "directory to serve from")
	flag.BoolVar(&cfg.debug, "debug", false, "verbose logging")
	flag.Parse()

	if cfg.debug {
		logger.level = LevelDebug
	}

	ps := new(pubsub)
	go watchFiles(cfg.dir, ps)

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFiles(cfg.dir))
	mux.HandleFunc("/events", handleEvents(ps))

	logger.Info("Starting server from %#v at %#v", cfg.dir, cfg.addr)

	err := http.ListenAndServe(cfg.addr, logRequest(mux))
	logger.Error("Failed to start: %s", err)
	os.Exit(1)
}

func serveFiles(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		fp := resolveFilepath(dir, r.URL.Path)

		if filepath.Ext(fp) == ".html" {
			serveHtml(w, r, fp)
		} else {
			http.ServeFile(w, r, fp)
		}
	}
}

func handleEvents(ps *pubsub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.(http.Flusher).Flush()

		done := r.Context().Done()
		ch := ps.subscribe()
		defer ps.unsubscribe(ch)

		for {
			select {
			case <-done:
				return

			case event := <-ch:
				fmt.Fprintf(w, "event: change\ndata: %s\n\n", event)
				w.(http.Flusher).Flush()
			}
		}
	}
}

// Serve a HTML file with injected <script>
func serveHtml(w http.ResponseWriter, r *http.Request, fp string) {
	html, err := os.ReadFile(fp)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	html = bytes.Replace(
		html,
		[]byte("</body>"),
		[]byte("<script>"+eventScript+"</script></body>"),
		1,
	)

	w.Write(html)
}
