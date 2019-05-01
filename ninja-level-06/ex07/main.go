package main

import "fmt"

func main() {
	x := foo
	x()
}

func foo() {
	fmt.Println("foo here")
}
