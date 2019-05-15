package main

import "fmt"

func main() {
	defer foo()
	fmt.Println("Hello from the other side")
}

func foo() {
	fmt.Println("This should be deferred")
}
