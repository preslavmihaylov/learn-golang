// Package quiz implements a set of methods for working with quizzes,
// containing simple questions (description + single answer)
package quiz

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

// Quiz encapsulates a set of questions
type Quiz struct {
	Questions []Question
}

// Question represents a single question, consisting of a description and an answer
type Question struct {
	Description string
	Answer      string
}

type quizResultParallel struct {
	question Question
	err      error
}

// New returns a new *Quiz after reading a set of questions from the provided CSV File.
// It returns an error in case there is a problem with the csv file or its format
func New(csvFilename string) (*Quiz, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("Caught error while opening file\n\t %s", err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Caught error while reading CSV file\n\t %s", err)
	}

	quiz := Quiz{}
	for rIndx, row := range records {
		if len(row) != 2 {
			return nil, fmt.Errorf(
				"Invalid CSV format. Expected 2 columns per row, Got %d on row %d", len(row), rIndx)
		}

		quiz.Questions = append(quiz.Questions, Question{Description: row[0], Answer: row[1]})
	}

	return &quiz, nil
}

// Run executes an interactive session of the provided quiz,
// by prompting an user with the set of questions on the provided writer.
// It returns an error in case of a problem with the provided reader or writer.
func (quiz *Quiz) Run(reader io.Reader, writer io.Writer) (cntCorrent int, err error) {
	return quiz.RunTimedQuiz(reader, writer, 30*time.Second)
}

// RunTimedQuiz executes a quiz which ends after the user finishes or a timer expires.
// It returns an error in case of a problem with the provided reader or writer.
func (quiz *Quiz) RunTimedQuiz(reader io.Reader, writer io.Writer, duration time.Duration) (cntCorrent int, err error) {
	resultsChan := make(chan quizResultParallel)
	go quiz.runQuizParallel(reader, writer, resultsChan)

	cntCorrect := 0
	quizComplete := false

	timer := time.NewTimer(duration)
	for !quizComplete {
		select {
		case res, ok := <-resultsChan:
			if res.err != nil {
				return cntCorrect, res.err
			}

			if !ok {
				quizComplete = true
			}

			cntCorrect++
		case <-timer.C:
			return cntCorrect, nil
		}
	}

	return cntCorrect, nil
}

// RunUntimedQuiz executes a quiz, which runs until the user finishes.
// It returns an error in case of a problem with the provided reader or writer.
func (quiz *Quiz) RunUntimedQuiz(reader io.Reader, writer io.Writer) (cntCorrent int, err error) {
	resultsChan := make(chan quizResultParallel)
	go quiz.runQuizParallel(reader, writer, resultsChan)

	cntCorrect := 0
	for res := range resultsChan {
		if res.err != nil {
			return cntCorrect, res.err
		}

		cntCorrect++
	}

	return cntCorrect, nil
}

func (quiz *Quiz) runQuizParallel(reader io.Reader, writer io.Writer, resultsChan chan quizResultParallel) {
	for i, q := range quiz.Questions {
		_, err := fmt.Fprintf(writer, "Problem #%d: %s = ", i+1, q.Description)
		if err != nil {
			resultsChan <- quizResultParallel{Question{},
				fmt.Errorf("Caught error while writing to provided writer\n\t %s", err)}
			close(resultsChan)
			return
		}

		var givenAnswer string
		_, err = fmt.Fscanln(reader, &givenAnswer)
		if err != nil {
			resultsChan <- quizResultParallel{Question{},
				fmt.Errorf("Caught error while reading from provided reader\n\t %s", err)}
			close(resultsChan)
			return
		}

		if givenAnswer == q.Answer {
			resultsChan <- quizResultParallel{q, nil}
		}
	}

	close(resultsChan)
	return
}
