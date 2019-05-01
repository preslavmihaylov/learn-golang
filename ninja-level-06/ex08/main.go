package main

import "fmt"

func main() {
	myFunc := getAFunc()
	myFunc()
}

func getAFunc() func() {
	return func() {
		fmt.Println("Here's your func")
	}
}
