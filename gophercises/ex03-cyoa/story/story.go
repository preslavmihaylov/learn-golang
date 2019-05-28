// Package story encapsulates routines for managing a Create Your Own Adventure story
package story

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
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []struct {
		Text  string `json:"text"`
		Title string `json:"arc"`
	} `json:"options"`
}

// FromJSONFile receives a json filename as argument, attempts to parse it and return
// a struct of type Story.
// In case of an issue with the file or json contents, an error is returned.
func FromJSONFile(jsonFilename string) (*Story, error) {
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

// ChapterByID attempts to find a chapter in the story by the provided id.
// In case one is not found, an error is returned.
func (st *Story) ChapterByID(id string) (Chapter, error) {
	chap, ok := st.chapterMap[id]
	if !ok {
		return Chapter{}, fmt.Errorf("Chapter with id %s not found", id)
	}

	return chap, nil
}

// IntroChapter attempts to find a chapter in the story by the intro id of the story.
// In case one is not found, an error is returned.
func (st *Story) IntroChapter() (Chapter, error) {
	chap, ok := st.chapterMap[st.intro]
	if !ok {
		return Chapter{}, fmt.Errorf("Intro Chapter not found. Intro ID: %s", st.intro)
	}

	return chap, nil
}

func newStory() *Story {
	return &Story{
		chapterMap: make(map[string]Chapter),
		intro:      "intro",
	}
}
