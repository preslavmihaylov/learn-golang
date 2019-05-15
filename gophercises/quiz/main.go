package main

import (
	"fmt"
	"log"
	"os"

	"github.com/preslavmihaylov/learn-golang/gophercises/quiz/quiz"
)

func main() {
	quiz, err := quiz.New("problems.csv")
	if err != nil {
		log.Fatalf("Received error while reading CSV file: %s\n", err)
	}

	cntCorrect, err := quiz.Run(os.Stdin)
	if err != nil {
		log.Fatalf("Received error while running quiz: %s\n", err)
	}

	fmt.Printf("You scored %d out of %d\n", cntCorrect, len(quiz.Questions))
}
