package main

import "fmt"

type customErr struct{}

func (e customErr) Error() string {
	return "This is a custom error"
}

func main() {
	err := foo()
	fmt.Println(err)
}

func foo() error {
	return customErr{}
}
