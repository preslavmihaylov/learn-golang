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
	var mux sync.Mutex

	for i := 0; i < 100; i++ {
		go func() {
			mux.Lock()
			v := counter
			runtime.Gosched()
			counter = v + 1
			mux.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("counter:", counter)
}
