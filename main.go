package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const eventScript = `<script>
new EventSource('/events').addEventListener('change', () => location.reload())
</script>`

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get cwd")
	}

	ps := new(pubsub)
	go watchFiles(cwd, ps)

	http.HandleFunc("/", serveFiles(cwd))
	http.HandleFunc("/events", handleEvents(ps))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveFiles(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

	html = bytes.Replace(html, []byte("</body>"), []byte(eventScript+"</body>"), 1)

	w.Write(html)
}
