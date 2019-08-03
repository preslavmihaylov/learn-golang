package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// program renames all somename_XXX.txt to XXX_somename.txt
	var renamePatt = regexp.MustCompile(`^[a-zA-Z]+\_(?P<num>[0-9]+)\.txt$`)
	numPatt := regexp.MustCompile(`[0-9]+`)
	filepath.Walk("./sample", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("failed traversing path %s: %s", path, err)
			return nil
		}

		filename := filepath.Base(path)
		if !info.IsDir() && renamePatt.MatchString(filename) {
			matches := numPatt.FindSubmatch([]byte(filename))
			renamed := fmt.Sprintf("%s/%s_%s.txt",
				filepath.Dir(path), string(matches[0]), strings.Split(filename, "_")[0])

			os.Rename(path, renamed)
		}

		return nil
	})
}
