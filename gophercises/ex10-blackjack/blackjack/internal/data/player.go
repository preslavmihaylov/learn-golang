package data

import (
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	bjapi "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type Player interface {
	Name() string
	Hand() []decks.Card
	Deal(c decks.Card)
	Split() (Player, Player)
	Unsplit(p2 Player)
	Balance() int
	Bet(amount int)
	DoubleDown()
	Payout(coef float64)
	LoseBet()
	Discard() []decks.Card
	Score() int
	IsSplit() bool
	HasSoftScore() bool
	HasBlackjack() bool
	CanDoubleDown() bool
	CanSplit() bool
	Busted() bool
}

type player struct {
	name    string
	isSplit bool
	balance int
	bet     int
	hand    []decks.Card
}

func NewPlayer(name string) Player {
	return &player{name, false, 0, 0, []decks.Card{}}
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

func (p *player) IsSplit() bool {
	return p.isSplit
}

func (p *player) Split() (Player, Player) {
	if len(p.hand) != 2 {
		log.Fatalf("Cannot split players. Invalid hand length")
	}

	var p1, p2 player
	p1 = *p
	p1.name += " (split 1)"
	p1.isSplit = true
	p1.hand = []decks.Card{p.hand[0]}

	p2 = *p
	p2.name += " (split 2)"
	p2.isSplit = true
	p2.hand = []decks.Card{p.hand[1]}
	p2.balance = -p.bet

	return &p1, &p2
}

func (p1 *player) Unsplit(p2 Player) {
	suffixLen := len(" (split 1)")
	p1.name = p1.name[:len(p1.name)-suffixLen]
	p1.isSplit = false
	p1.balance += p2.Balance()
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

func (p *player) CanDoubleDown() bool {
	if len(p.hand) != 2 {
		return false
	}

	score := p.Score()
	if p.HasSoftScore() {
		score -= 10
	}

	return score >= 9 && score <= 11
}

func (p *player) CanSplit() bool {
	return len(p.hand) == 2 &&
		p.hand[0].Rank == p.hand[1].Rank &&
		!p.IsSplit()
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
