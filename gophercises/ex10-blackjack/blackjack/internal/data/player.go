package data

import (
	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	bjapi "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type Player interface {
	Name() string
	Hand() []decks.Card
	Deal(c decks.Card)
	Balance() int
	Bet(amount int)
	DoubleDown()
	Payout(coef float64)
	LoseBet()
	Discard() []decks.Card
	Score() int
	HasSoftScore() bool
	HasBlackjack() bool
	CanDouble() bool
	Busted() bool
}

type player struct {
	name    string
	balance int
	bet     int
	hand    []decks.Card
}

func NewPlayer(name string) Player {
	return &player{name, 0, 0, []decks.Card{}}
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Balance() int {
	return p.balance
}

func (p *player) Bet(amount int) {
	p.balance -= amount
	p.bet = amount
}

func (p *player) DoubleDown() {
	p.balance -= p.bet
	p.bet += p.bet
}

func (p *player) LoseBet() {
	p.bet = 0
}

func (p *player) Payout(coef float64) {
	winning := float64(p.bet) * coef
	p.balance += p.bet + int(winning)
	p.bet = 0
}

func (p *player) Score() int {
	return bjapi.CalculateScore(p.hand).Value
}

func (p *player) HasSoftScore() bool {
	return bjapi.CalculateScore(p.hand).IsSoft
}

func (p *player) HasBlackjack() bool {
	if len(p.hand) != 2 {
		return false
	}

	hasTen := p.hand[0].Rank >= decks.Ten || p.hand[1].Rank >= decks.Ten
	hasAce := p.hand[0].Rank == decks.Ace || p.hand[1].Rank == decks.Ace
	return hasTen && hasAce
}

func (p *player) CanDouble() bool {
	if len(p.hand) != 2 {
		return false
	}

	score := p.Score()
	if p.HasSoftScore() {
		score -= 10
	}

	return score >= 9 && score <= 11
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
