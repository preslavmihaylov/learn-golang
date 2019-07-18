package data

import "github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"

type Dealer interface {
	Player
	Revealed() bool
	Reveal()
	Hide()
}

type dealer struct {
	Player
	revealed bool
}

func NewDealer(name string) Dealer {
	return &dealer{Player: NewPlayer(nil, name)}
}

func (d *dealer) Score() int {
	if len(d.Hand()) == 0 {
		return 0
	}

	if d.Revealed() {
		return calculateScore(d.Hand()).value
	}

	return calculateScore(d.Hand()[:1]).value
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

func (d *dealer) Deal(c decks.Card) {
	d.Player.Deal(c)
}

func (d *dealer) Revealed() bool {
	return d.revealed
}

func (d *dealer) Reveal() {
	d.revealed = true
}

func (d *dealer) Hide() {
	d.revealed = false
}
