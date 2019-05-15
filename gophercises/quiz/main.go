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
		log.Fatalf("\nReceived error while constructing new quiz\n\t %s\n", err)
	}

	cntCorrect, err := quiz.Run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("\nReceived error while running quiz\n\t %s\n", err)
	}

	fmt.Printf("You scored %d out of %d\n", cntCorrect, len(quiz.Questions))
}
