package urlshort

import (
	"fmt"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then a default handler
// will be called.
func MapHandler(redirs map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Path
		if redir, ok := redirs[u]; ok {
			http.Redirect(w, r, redir, 301)
		} else {
			defaultHandler(w, r)
		}
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "No such route found")
}
