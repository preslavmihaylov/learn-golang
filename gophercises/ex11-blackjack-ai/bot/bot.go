package bot

import "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"

func Simulate(roundsCnt int) BlackjackStats {
	s := NewStrategy(roundsCnt)
	blackjack.Play(1, s)

	return s.Stats
}
