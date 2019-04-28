package main

import "fmt"

func main() {
	for i := 'A'; i <= 'Z'; i++ {
		for j := 0; j < 3; j++ {
			fmt.Printf("%c ", i)
		}

		fmt.Println()
	}
}
