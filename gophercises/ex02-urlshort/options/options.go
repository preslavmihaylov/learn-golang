// Package options encapsulates routines for parsing the program's command line arguments
package options

import "flag"

// O encapsulates the parsed command line arguments
type O struct {
	YAMLFilename string
}

// ParseArgs parses the command line arguments and returns a struct O, containing the results
func ParseArgs() *O {
	opts := O{}
	flag.StringVar(&opts.YAMLFilename, "yaml", "default.yaml",
		"a yaml file with fields in the format \n\t- path: {path}\n\t  url: {url}\n")
	flag.Parse()

	return &opts
}
