//go:generate stringer -type=Outcome

package api

import "github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"

type GameEvent interface{}

type DeckShuffledEvent struct{}

type StartBetEvent struct {
	PlayerName string
	Balance    int
}

type BetEvent struct {
	PlayerName string
	Bet        int
}

type DealCardsEvent struct {
	Hands      map[string][]decks.Card
	DealerHand []decks.Card
}

type PlayerTurnEvent struct {
	PlayerName string
	PlayerHand []decks.Card
	DealerHand []decks.Card
}

type DealerTurnEvent struct {
	PlayersInGame  map[string][]decks.Card
	DealerHand     []decks.Card
	DealerRevealed bool
}

type HitEvent struct {
	PlayerName string
	Card       decks.Card
	Busted     bool
}

type StandEvent struct {
	PlayerName string
}

type BlackjackEvent struct {
	PlayerName string
}

type DoubleDownEvent struct {
	PlayerName string
	Card       decks.Card
}

type SplitEvent struct {
	PlayerName string
}

type RoundEndsEvent struct{}

type Outcome uint8

const (
	Lost Outcome = iota
	Tied
	Won
	PlayerBlackjack
	Busted
	DealerBusted
)

type Result struct {
	Outcome     Outcome
	PlayerScore int
	DealerScore int
}

type ResolveEvent struct {
	Results map[string]Result
}
