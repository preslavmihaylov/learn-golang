package data

import "github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"

type GameInfo struct {
	PlayerHand []decks.Card
	DealerHand []decks.Card
}

type PlayerInterface interface {
	Listen(e GameEvent)
	PlayerTurn(gi GameInfo, actions []Action) Action
}

func EmitEvent(players []Player, e GameEvent) {
	for _, p := range players {
		p.Interface().Listen(e)
	}
}
