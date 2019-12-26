package middleware

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

func Development(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				stack := string(debug.Stack())
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, withLinksToSrc(stack, "/debug"))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func withLinksToSrc(stackTrace, endpoint string) string {
	patt := regexp.MustCompile(`(((/[a-zA-Z0-9.-]+)+/[a-zA-Z0-9.-]+\.go):([0-9]+))`)
	return patt.ReplaceAllStringFunc(stackTrace, func(match string) string {
		toks := strings.Split(match, ":")
		file, line := toks[0], toks[1]

		return fmt.Sprintf("<a href=\"%s?file=%s&line=%s\">%s</a>", endpoint, file, line, match)
	})
}
