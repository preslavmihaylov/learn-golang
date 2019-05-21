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

// Q stands for 'Quiz' and encapsulates a set of questions
type Q struct {
	questions []question
	duration  time.Duration
}

// question represents a single question, consisting of a description and an answer
type question struct {
	desc string
	ans  string
}

// FromCSVFileTimed creates a timed quiz with the specified duration from the provided CSV file.
// It returns an error in case of a problem with the provided csv file
func FromCSVFileTimed(csvFilename string, duration time.Duration) (*Q, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("Caught error while opening file\n\t %s", err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Caught error while reading CSV file\n\t %s", err)
	}

	qz := Q{}
	for rIndx, row := range records {
		if len(row) != 2 {
			return nil, fmt.Errorf(
				"Invalid CSV format. Expected 2 columns per row, Got %d on row %d", len(row), rIndx)
		}

		qz.questions = append(qz.questions, question{desc: row[0], ans: row[1]})
	}

	qz.duration = duration

	return &qz, nil
}

// FromCSVFileUntimed creates an untimed quiz from the provided CSV file.
// It returns an error in case of a problem with the provided csv file
func FromCSVFileUntimed(csvFilename string) (*Q, error) {
	return FromCSVFileTimed(csvFilename, time.Duration(math.MaxInt64))
}

// QuestionsCnt returns the count of questions in quiz
func (qz *Q) QuestionsCnt() int {
	return len(qz.questions)
}

type correctQuestion struct {
	q   question
	err error
}

// Run executes an interactive session of the provided quiz,
// by prompting an user with the set of questions on the provided writer.
// In case of a timed quiz, the interactive session will end after the provided duration expires.
// It returns an error in case of a problem with the provided reader or writer.
func (qz *Q) Run(reader io.Reader, writer io.Writer) (cntCorrent int, err error) {
	correctQuestionsChan := make(chan correctQuestion)
	go qz.executeParallel(reader, writer, correctQuestionsChan)

	cntCorrect := 0
	qzComplete := false

	timer := time.NewTimer(qz.duration)
	for !qzComplete {
		select {
		case res, isOpen := <-correctQuestionsChan:
			if !isOpen {
				qzComplete = true
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

func (qz *Q) executeParallel(reader io.Reader, writer io.Writer, resultsChan chan correctQuestion) {
	for i, qstion := range qz.questions {
		_, err := fmt.Fprintf(writer, "Problem #%d: %s = ", i+1, qstion.desc)
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

		if givenAnswer == qstion.ans {
			resultsChan <- correctQuestion{qstion, nil}
		}
	}

	close(resultsChan)
	return
}
