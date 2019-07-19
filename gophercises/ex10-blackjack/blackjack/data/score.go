package data

import (
	"math"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type Score struct {
	Value  int
	IsSoft bool
}

func CalculateScore(hand []decks.Card) Score {
	scores := []Score{Score{}}
	for _, c := range hand {
		if c.Rank >= decks.Two && c.Rank <= decks.Ten {
			for i := range scores {
				scores[i].Value += int(c.Rank) + 1
			}
		} else if c.Rank == decks.Ace {
			for i := range scores {
				scores = append(scores, Score{scores[i].Value + 11, true})
				scores[i].Value += 1
			}
		} else {
			for i := range scores {
				scores[i].Value += 10
			}
		}
	}

	maxNotBusted := Score{0, false}
	minBusted := Score{math.MaxInt32, false}
	for _, sc := range scores {
		val := sc.Value
		if val <= 21 && val > maxNotBusted.Value {
			maxNotBusted = sc
		}

		if val > 21 && val < minBusted.Value {
			minBusted = sc
		}
	}

	if maxNotBusted.Value == 0 {
		return minBusted
	}

	return maxNotBusted
}
