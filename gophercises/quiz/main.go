package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/quiz/quiz"
)

const invalidTime = -1

func main() {
	var csvFilename string
	var quizTimeSecs int
	flag.StringVar(&csvFilename, "csv", "problems.csv",
		"a csv file in the format 'question,answer'. Default is 'problems.csv'")
	flag.IntVar(&quizTimeSecs, "limit", invalidTime,
		"the time limit for taking the quiz in seconds. Default is forever.")
	flag.Parse()

	var q *quiz.Q
	var err error
	if quizTimeSecs != invalidTime {
		q, err = quiz.FromCSVFileTimed(csvFilename, time.Duration(quizTimeSecs)*time.Second)
	} else {
		q, err = quiz.FromCSVFileUntimed(csvFilename)
	}

	if err != nil {
		log.Fatalf("\nReceived error while constructing new quiz\n\t %s\n", err)
	}

	cntCorrect, err := q.Run(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatalf("\nReceived error while running quiz\n\t %s\n", err)
	}

	fmt.Printf("You scored %d out of %d\n", cntCorrect, q.QuestionsCnt())
}
