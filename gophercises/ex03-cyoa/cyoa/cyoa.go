// Package cyoa encapsulates routines for managing a Create Your Own Adventure story
package cyoa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Story encapsulates a set of chapters, indexed by string id
type Story struct {
	chapterMap map[string]Chapter
	intro      string
}

// Chapter encapsulates a chapter of a story
type Chapter struct {
	Title   string   `json:"title"`
	Text    []string `json:"story"`
	Options []struct {
		Text  string `json:"text"`
		Title string `json:"arc"`
	} `json:"options"`
}

// ParseStory receives a json filename as argument, attempts to parse it and return
// a struct of type Story.
// In case of an issue with the file or json contents, an error is returned.
func ParseStory(jsonFilename string) (*Story, error) {
	jsonBytes, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read json file: %s", err)
	}

	story := newStory()
	storyRawJSON := map[string]*json.RawMessage{}
	err = json.Unmarshal(jsonBytes, &storyRawJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json file: %s", err)
	}

	for chapterID, ChapterRawJSON := range storyRawJSON {
		currChapter := Chapter{}
		err = json.Unmarshal([]byte(*ChapterRawJSON), &currChapter)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal chapter %s: %s", chapterID, err)
		}

		story.chapterMap[chapterID] = currChapter
	}

	if _, ok := story.chapterMap[story.intro]; !ok {
		return nil, fmt.Errorf(
			"story is missing an intro chapter. Expected intro chapter with id: %s",
			story.intro)
	}

	return story, nil
}

func newStory() *Story {
	return &Story{
		chapterMap: make(map[string]Chapter),
		intro:      "intro",
	}
}
