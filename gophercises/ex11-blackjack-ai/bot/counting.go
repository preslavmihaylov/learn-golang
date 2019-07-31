package bot

import (
	"math"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type BlackjackCounting struct {
	runningCnt int
	discarded  int
	totalDecks int
	minTrueCnt int
}

func (bc *BlackjackCounting) reset() {
	bc.runningCnt = 0
	bc.discarded = 0
}

func (bc *BlackjackCounting) runCount(cards ...decks.Card) {
	bc.discarded += len(cards)
	for _, c := range cards {
		bc.runningCnt += countValue(c)
	}
}

func (bc *BlackjackCounting) unitsToBet() int {
	discardedDecks := bc.discarded / decks.DeckSize
	decksLeft := float64(bc.totalDecks - discardedDecks)
	if bc.discarded-(discardedDecks*decks.DeckSize) >= decks.DeckSize/2 {
		decksLeft -= 0.5
	}

	trueCnt := int(math.Round(float64(bc.runningCnt) / decksLeft))

	if trueCnt < bc.minTrueCnt {
		return 0
	}

	return trueCnt
}

func countValue(c decks.Card) int {
	if c.Rank >= decks.Ten || c.Rank == decks.Ace {
		return -1
	} else if c.Rank >= decks.Seven && c.Rank <= decks.Nine {
		return 0
	} else {
		return 1
	}
}
