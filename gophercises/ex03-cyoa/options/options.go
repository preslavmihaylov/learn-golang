// Package options encapsulates routines for parsing the program's command line arguments
package options

import "flag"

// O encapsulates the parsed command line arguments
type O struct {
	JSONFilename         string
	HTMLTemplateFilename string
}

// ParseArgs parses the command line arguments and returns a struct O, containing the results
func ParseArgs() *O {
	opts := O{}

	opts.JSONFilename = *flag.String("json", "gopher.json", "a JSON file with an encoded CYOA story")
	opts.HTMLTemplateFilename = *flag.String(
		"template", "story.html", "an HTML Template used to render the story in the server")
	flag.Parse()

	return &opts
}
