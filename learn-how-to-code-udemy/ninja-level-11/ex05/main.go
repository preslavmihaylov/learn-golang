package main

import (
	"errors"
	"log"
)

func main() {
	_, err := Sqrt(-10.23)
	if err != nil {
		log.Println(err)
	}
}

// Sqrt returns the square root of the provided number.
// In case of a negative number, an error is returned.
func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("Received a negative number")
	}
	return 42, nil
}
