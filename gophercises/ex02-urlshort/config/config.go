package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ReadYAMLConfig reads the contents of the provided yaml config file and populates the provided mapping
// in the form [Path] -> Url.
// It expects the provided yaml config file to be in the format:
// 		- path: {path}
//		  url: {url}
// It returns an error in case of an invalid yaml format or io error.
func ReadYAMLConfig(yamlFilename string, redirs map[string]string) error {
	yml, err := ioutil.ReadFile(yamlFilename)
	if err != nil {
		return fmt.Errorf("received error while reading yaml file: %s", err)
	}

	err = parseYAML(yml, redirs)
	if err != nil {
		return fmt.Errorf("received error while parsing yaml: %s", err)
	}

	return nil
}

func parseYAML(yml []byte, redirs map[string]string) error {
	type redirectDirective struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}

	parsedRedirs := []redirectDirective{}
	err := yaml.Unmarshal(yml, &parsedRedirs)
	if err != nil {
		return fmt.Errorf("Failed unmarshaling yaml. Received err: %s", err)
	}

	for _, redir := range parsedRedirs {
		redirs[redir.Path] = redir.URL
	}

	return nil
}
