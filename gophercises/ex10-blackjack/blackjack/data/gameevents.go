package data

import "github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"

type GameEvent interface{}

type DealCardsEvent struct {
	PlayerHand []decks.Card
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

type RoundEndsEvent struct{}

type Outcome uint8

const (
	Lost Outcome = iota
	Tied
	Won
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
