package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go foo()
	go bar()

	wg.Wait()
	fmt.Println("main")
}

func foo() {
	fmt.Println("foo")
	wg.Done()
}

func bar() {
	fmt.Println("bar")
	wg.Done()
}
