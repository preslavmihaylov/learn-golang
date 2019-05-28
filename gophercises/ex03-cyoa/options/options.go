// Package options encapsulates routines for parsing the program's command line arguments
package options

import "flag"

// O encapsulates the parsed command line arguments
type O struct {
	JSONFilename         string
	HTMLTemplateFilename string
}

var opts O

func init() {
	opts = O{}
	opts.JSONFilename = *flag.String("json", "gopher.json", "a JSON file with an encoded CYOA story")
	opts.HTMLTemplateFilename = *flag.String(
		"template", "story.html", "an HTML Template used to render the story in the server")
	flag.Parse()
}

// JSONFilename returns the parsed cmd line option "json"
func JSONFilename() string {
	return opts.JSONFilename
}

// HTMLTemplateFilename returns the parsed cmd line option "template"
func HTMLTemplateFilename() string {
	return opts.HTMLTemplateFilename
}
