package data

import (
	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type Player interface {
	Name() string
	Hand() []decks.Card
	Deal(c decks.Card)
	Discard() []decks.Card
	Score() int
	IsSoftScore() bool
	Busted() bool
}

type player struct {
	name string
	hand []decks.Card
}

func NewPlayer(name string) Player {
	return &player{name, []decks.Card{}}
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Score() int {
	return calculateScore(p.hand).value
}

func (p *player) IsSoftScore() bool {
	return calculateScore(p.hand).isSoft
}

func (p *player) Hand() []decks.Card {
	return p.hand
}

func (p *player) Discard() []decks.Card {
	cards := p.hand
	p.hand = nil

	return cards
}

func (p *player) Deal(c decks.Card) {
	p.hand = append(p.hand, c)
}

func (p *player) Busted() bool {
	return p.Score() > 21
}