// Package quiz implements a set of methods for working with quizzes,
// containing simple questions (description + single answer)
package quiz

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

// Quiz encapsulates a set of questions
type Quiz struct {
	questions []question
	duration  time.Duration
}

// question represents a single question, consisting of a description and an answer
type question struct {
	description string
	answer      string
}

// FromCSVFileTimed creates a timed quiz with the specified duration from the provided CSV file.
// It returns an error in case of a problem with the provided csv file
func FromCSVFileTimed(csvFilename string, duration time.Duration) (*Quiz, error) {
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

		quiz.questions = append(quiz.questions, question{description: row[0], answer: row[1]})
	}

	quiz.duration = duration

	return &quiz, nil
}

// FromCSVFileUntimed creates an untimed quiz from the provided CSV file.
// It returns an error in case of a problem with the provided csv file
func FromCSVFileUntimed(csvFilename string) (*Quiz, error) {
	return FromCSVFileTimed(csvFilename, time.Duration(math.MaxInt64))
}

// QuestionsCnt returns the count of questions in quiz
func (quiz *Quiz) QuestionsCnt() int {
	return len(quiz.questions)
}

type correctQuestion struct {
	question question
	err      error
}

// Run executes an interactive session of the provided quiz,
// by prompting an user with the set of questions on the provided writer.
// In case of a timed quiz, the interactive session will end after the provided duration expires.
// It returns an error in case of a problem with the provided reader or writer.
func (quiz *Quiz) Run(reader io.Reader, writer io.Writer) (cntCorrent int, err error) {
	resultsChan := make(chan correctQuestion)
	go quiz.executeParallel(reader, writer, resultsChan)

	cntCorrect := 0
	quizComplete := false

	timer := time.NewTimer(quiz.duration)
	for !quizComplete {
		select {
		case res, isOpen := <-resultsChan:
			if !isOpen {
				quizComplete = true
				break
			}

			if res.err != nil {
				return cntCorrect, res.err
			}

			cntCorrect++
		case <-timer.C:
			return cntCorrect, nil
		}
	}

	return cntCorrect, nil
}

func (quiz *Quiz) executeParallel(reader io.Reader, writer io.Writer, resultsChan chan correctQuestion) {
	for i, q := range quiz.questions {
		_, err := fmt.Fprintf(writer, "Problem #%d: %s = ", i+1, q.description)
		if err != nil {
			resultsChan <- correctQuestion{question{},
				fmt.Errorf("Caught error while writing to provided writer\n\t %s", err)}
			close(resultsChan)
			return
		}

		var givenAnswer string
		_, err = fmt.Fscanln(reader, &givenAnswer)
		if err != nil {
			resultsChan <- correctQuestion{question{},
				fmt.Errorf("Caught error while reading from provided reader\n\t %s", err)}
			close(resultsChan)
			return
		}

		if givenAnswer == q.answer {
			resultsChan <- correctQuestion{q, nil}
		}
	}

	close(resultsChan)
	return
}
