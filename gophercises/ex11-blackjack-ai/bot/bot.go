package bot

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"
)

func Simulate(roundsCnt, decksCnt, betUnit, minTrueCount int) *BlackjackStats {
	var wg sync.WaitGroup
	var mux sync.Mutex
	finalStats := &BlackjackStats{}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			s := NewStrategy(roundsCnt/runtime.NumCPU(), decksCnt, betUnit, minTrueCount)
			blackjack.Play(decksCnt, 1, s)

			mux.Lock()
			finalStats.Accumulate(&s.Stats)
			mux.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()
	return finalStats
}
