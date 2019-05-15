package main

import (
	"fmt"
)

func main() {
	c := make(chan int)

	go func(c chan int) {
		c <- 42
		close(c)
	}(c)

	v, ok := <-c
	fmt.Println(v, ok)

	v, ok = <-c
	fmt.Println(v, ok)
}
