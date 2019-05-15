package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(100)
	counter := 0
	for i := 0; i < 100; i++ {
		go func() {
			v := counter
			runtime.Gosched()
			counter = v + 1
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("counter:", counter)
}
