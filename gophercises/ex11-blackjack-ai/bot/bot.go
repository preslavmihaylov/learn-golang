package bot

import "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"

func Simulate(roundsCnt, decksCnt int) BlackjackStats {
	s := NewStrategy(roundsCnt, decksCnt, 25, 4)
	blackjack.Play(decksCnt, 1, s)

	return s.Stats
}
