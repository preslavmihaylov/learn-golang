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
	urlShortenerMux := http.NewServeMux()
	urlShortenerMux.HandleFunc("/", defaultHandler)

	for url, RedirectURL := range redirs {
		urlShortenerMux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, RedirectURL, 301)
		})
	}

	return urlShortenerMux.ServeHTTP
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "No such route found")
}
