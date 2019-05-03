package main

import "fmt"

func main() {
	c := make(chan int)
	go func(c chan int) {
		for i := 0; i < 100; i++ {
			c <- i
		}

		close(c)
	}(c)

	for v := range c {
		fmt.Println(v)
	}
}
