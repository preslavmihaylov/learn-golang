package main

import "fmt"

func main() {
	for i := 0; i < 100; i++ {
		if (i % 3) == 0 {
			fmt.Println("Fizz")
		}

		if (i % 7) == 0 {
			fmt.Println("Buzz")
		}

		if (i % 21) == 0 {
			fmt.Println("FizzBuzz")
		}
	}
}
