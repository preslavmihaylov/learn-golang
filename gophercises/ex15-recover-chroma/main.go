package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/debug", debugHandler)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := string(debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, withLinksToSrc(stack, "/debug"))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func withLinksToSrc(stackTrace, endpoint string) string {
	fmt.Println("finding matches...")
	patt := regexp.MustCompile(`(((/[a-zA-Z0-9.-]+)+/[a-zA-Z0-9.-]+\.go):([0-9]+))`)
	return patt.ReplaceAllStringFunc(stackTrace, func(match string) string {
		toks := strings.Split(match, ":")
		file, line := toks[0], toks[1]

		return fmt.Sprintf("<a href=\"%s?file=%s&line=%s\">%s</a>", endpoint, file, line, match)
	})
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	if len(urlParams["file"]) != 1 {
		fmt.Fprintf(w, "url param 'file' is missing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if len(urlParams["line"]) != 1 {
		fmt.Fprintf(w, "url param 'line' is missing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file := urlParams["file"][0]
	line, err := strconv.Atoi(urlParams["line"][0])
	if err != nil {
		fmt.Fprintf(w, "failed to parse line: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(w, "couldn't read file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.Standalone(), html.HighlightLines([][2]int{[2]int{line, line}}))
	lexer := lexers.Get("go")

	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		fmt.Fprintf(w, "tokenise file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = formatter.Format(w, style, iterator)
	// err = quick.Highlight(w, string(contents), "go", "html", "monokai")
	if err != nil {
		fmt.Fprintf(w, "failed to render file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
