package data

import (
	"math"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type Score struct {
	value  int
	isSoft bool
}

func calculateScore(hand []decks.Card) Score {
	scores := []Score{Score{}}
	for _, c := range hand {
		if c.Rank >= decks.Two && c.Rank <= decks.Ten {
			for i := range scores {
				scores[i].value += int(c.Rank) + 1
			}
		} else if c.Rank == decks.Ace {
			for i := range scores {
				scores = append(scores, Score{scores[i].value + 11, true})
				scores[i].value += 1
			}
		} else {
			for i := range scores {
				scores[i].value += 10
			}
		}
	}

	maxNotBusted := Score{0, false}
	minBusted := Score{math.MaxInt32, false}
	for _, sc := range scores {
		val := sc.value
		if val <= 21 && val > maxNotBusted.value {
			maxNotBusted = sc
		}

		if val > 21 && val < minBusted.value {
			minBusted = sc
		}
	}

	if maxNotBusted.value == 0 {
		return minBusted
	}

	return maxNotBusted
}
