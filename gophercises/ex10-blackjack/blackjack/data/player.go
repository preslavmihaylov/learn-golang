package data

import (
	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type Player interface {
	Interface() PlayerInterface
	Name() string
	Hand() []decks.Card
	Deal(c decks.Card)
	Discard() []decks.Card
	Score() int
	IsSoftScore() bool
	Busted() bool
}

type player struct {
	playerInterface PlayerInterface
	name            string
	hand            []decks.Card
}

func NewPlayer(pi PlayerInterface, name string) Player {
	return &player{pi, name, []decks.Card{}}
}

func (p *player) Interface() PlayerInterface {
	return p.playerInterface
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Score() int {
	return CalculateScore(p.hand).Value
}

func (p *player) IsSoftScore() bool {
	return CalculateScore(p.hand).IsSoft
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
