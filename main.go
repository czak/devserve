package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed event.js
var eventScript string

type config struct {
	addr string
	dir  string
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.addr, "addr", ":8080", "network address")
	flag.StringVar(&cfg.dir, "dir", ".", "directory to serve from")
	flag.Parse()

	ps := new(pubsub)
	go watchFiles(cfg.dir, ps)

	http.HandleFunc("/", serveFiles(cfg.dir))
	http.HandleFunc("/events", handleEvents(ps))

	log.Printf("Serving files from %#v on %#v\n", cfg.dir, cfg.addr)

	log.Fatal(http.ListenAndServe(cfg.addr, nil))
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

// Serve a HTML file with injected <script>
func serveHtml(w http.ResponseWriter, r *http.Request, fp string) {
	html, err := os.ReadFile(fp)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	html = bytes.Replace(html, []byte("</body>"), []byte("<script>"+eventScript+"</script></body>"), 1)

	w.Write(html)
}
