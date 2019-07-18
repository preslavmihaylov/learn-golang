package data

import "github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"

type GameInfo struct {
	PlayerHand []decks.Card
	DealerHand []decks.Card
}

type PlayerInterface interface {
	PlayerTurn(gi GameInfo, actions []Action) Action
}
