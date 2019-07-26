package data

import (
	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type Dealer interface {
	Player
	Reveal()
	Hide()
	Revealed() bool
}

type dealer struct {
	Player
	revealed bool
}

func NewDealer(name string) Dealer {
	return &dealer{Player: NewPlayer(name)}
}

func (d *dealer) Score() int {
	if d.Revealed() {
		return d.Player.Score()
	}

	return api.CalculateScore(d.Hand()[:1]).Value
}

func (d *dealer) Hand() []decks.Card {
	if len(d.Player.Hand()) == 0 {
		return []decks.Card{}
	}

	if d.Revealed() {
		return d.Player.Hand()
	}

	return d.Player.Hand()[:1]
}

func (d *dealer) Reveal() {
	d.revealed = true
}

func (d *dealer) Hide() {
	d.revealed = false
}

func (d *dealer) Revealed() bool {
	return d.revealed
}
