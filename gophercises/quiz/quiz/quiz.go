// Package quiz implements a set of methods for working with quizzes,
// containing simple questions (description + single answer)
package quiz

import "io"

// Quiz encapsulates a set of questions
type Quiz struct {
	Questions []Question
}

// Question represents a single question, consisting of a description and an answer
type Question struct {
	Description string
	Answer      string
}

// New returns a new *Quiz after reading a set of questions from the provided CSV File.
// It returns an error in case there is a problem with the csv file or its format
func New(csvFilename string) (*Quiz, error) {
	return &Quiz{}, nil
}

// Run executes an interactive session of the provided quiz,
// by prompting an user with the set of questions on the provided writer.
// TODO: Complete error description
func (quiz *Quiz) Run(writer io.Writer) (cntCorrent int, err error) {
	return 0, nil
}
