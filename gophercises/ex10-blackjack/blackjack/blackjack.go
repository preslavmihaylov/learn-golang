package blackjack

import "github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
import "math"

type Player interface {
	Hand() []decks.Card
	Deal(c decks.Card)
	Discard() []decks.Card
	Score() int
	IsSoftScore() bool
	Busted() bool
}

type Dealer interface {
	Player
	Revealed() bool
	Reveal()
	Hide()
}

type BlackjackData struct {
	deck       *decks.Deck
	discarded  []decks.Card
	players    []Player
	dealer     Dealer
	playerTurn int
}

func (bd *BlackjackData) Discard(cards []decks.Card) {
	bd.discarded = append(bd.discarded, cards...)
}

func (bd *BlackjackData) IsDealersTurn() bool {
	return bd.playerTurn >= len(bd.players)
}

func (bd *BlackjackData) NextPlayersTurn() {
	bd.playerTurn++
}

func (bd *BlackjackData) NewRound() {
	bd.playerTurn = 0
	for i := range bd.players {
		bd.Discard(bd.players[i].Discard())
	}

	bd.Discard(bd.dealer.Discard())
	bd.dealer.Hide()
}

type player struct {
	hand []decks.Card
}

func (p *player) Score() int {
	return score(p.hand)
}

func (p *player) IsSoftScore() bool {
	for _, c := range p.hand {
		if c.Rank == decks.Ace {
			return true
		}
	}

	return false
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

type dealer struct {
	player
	revealed bool
}

func (d *dealer) Score() int {
	if len(d.hand) == 0 {
		return 0
	}

	if d.Revealed() {
		return score(d.hand)
	}

	return score(d.hand[:1])
}

func (d *dealer) Hand() []decks.Card {
	if len(d.hand) == 0 {
		return []decks.Card{}
	}

	if d.Revealed() {
		return d.hand
	}

	return d.hand[:1]
}

func (d *dealer) Deal(c decks.Card) {
	d.hand = append(d.hand, c)
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

func score(hand []decks.Card) int {
	scores := []int{0}
	for _, c := range hand {
		if c.Rank >= decks.Two && c.Rank <= decks.Ten {
			for i := range scores {
				scores[i] += int(c.Rank) + 1
			}
		} else if c.Rank == decks.Ace {
			for i := range scores {
				scores = append(scores, scores[i]+11)
				scores[i] += 1
			}
		} else {
			for i := range scores {
				scores[i] += 10
			}
		}
	}

	maxNotBusted := 0
	minBusted := math.MaxInt32
	for _, sc := range scores {
		if sc <= 21 && sc > maxNotBusted {
			maxNotBusted = sc
		}

		if sc > 21 && sc < minBusted {
			minBusted = sc
		}
	}

	if maxNotBusted == 0 {
		return minBusted
	}

	return maxNotBusted
}
