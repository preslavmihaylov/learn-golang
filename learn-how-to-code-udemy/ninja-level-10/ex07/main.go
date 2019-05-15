package main

import "fmt"
import "sync"

func main() {
	c := launch()

	for v := range c {
		fmt.Println(v)
	}
}

func launch() chan int {
	c := make(chan int)

	go monitorGofuncs(c)

	return c
}

func monitorGofuncs(c chan int) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(c chan int, offset int) {
			for j := offset; j < offset+10; j++ {
				c <- j
			}

			wg.Done()
		}(c, i*10)
	}

	wg.Wait()
	close(c)
}
