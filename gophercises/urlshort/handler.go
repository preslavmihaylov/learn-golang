package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	urlShortenerMux := http.NewServeMux()
	for url, RedirectURL := range pathsToUrls {
		urlShortenerMux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, RedirectURL, 301)
		})
	}

	return urlShortenerMux.ServeHTTP
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathToURLMap, err := parseYAMLRedirectDirectives(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(pathToURLMap, fallback), nil
}

func parseYAMLRedirectDirectives(yml []byte) (map[string]string, error) {
	type redirectDirective struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}

	redirs := []redirectDirective{}
	err := yaml.Unmarshal(yml, &redirs)
	if err != nil {
		return nil, fmt.Errorf("Failed unmarshaling yaml. Received err: %s", err)
	}

	pathToURL := map[string]string{}
	for _, redir := range redirs {
		pathToURL[redir.Path] = redir.URL
	}

	return pathToURL, nil
}
