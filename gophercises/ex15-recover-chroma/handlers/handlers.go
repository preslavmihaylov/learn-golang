package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func PanicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Oh no!")
}

func PanicAfterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	panic("Oh no!")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	file, line, err := extractFileAndLineFromURL(r)
	if err != nil {
		fmt.Fprintf(w, "failed to get url params: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = outputSourceFile(w, file, line)
	if err != nil {
		fmt.Fprintf(w, "failed to output source file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractFileAndLineFromURL(r *http.Request) (string, int, error) {
	urlParams := r.URL.Query()
	if len(urlParams["file"]) != 1 {
		return "", 0, fmt.Errorf("url param 'file' is missing")
	} else if len(urlParams["line"]) != 1 {
		return "", 0, fmt.Errorf("url param 'line' is missing")
	}

	file := urlParams["file"][0]
	line, err := strconv.Atoi(urlParams["line"][0])
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse line: %s", err)
	}

	return file, line, nil
}

func outputSourceFile(w http.ResponseWriter, file string, line int) error {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	formatter := html.New(html.Standalone(),
		html.HighlightLines([][2]int{[2]int{line, line}}),
		html.WithLineNumbers())
	iterator, err := lexers.Get("go").Tokenise(nil, string(contents))
	if err != nil {
		return fmt.Errorf("failed to tokenise file: %s", err)
	}

	err = formatter.Format(w, styles.Get("monokai"), iterator)
	if err != nil {
		return fmt.Errorf("failed to render file: %s", err)
	}

	return nil
}
