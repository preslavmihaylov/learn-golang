package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(100)
	var counter int64

	for i := 0; i < 100; i++ {
		go func() {
			atomic.AddInt64(&counter, 1)

			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("counter:", counter)
}
