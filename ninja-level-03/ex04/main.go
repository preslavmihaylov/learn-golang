package main

import "fmt"

func main() {
	year := 1985
	for {
		if year > 2019 {
			break
		}

		fmt.Println(year)
		year++
	}
}
