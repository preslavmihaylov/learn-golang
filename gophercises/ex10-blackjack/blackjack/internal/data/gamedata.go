package data

import (
	"fmt"
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	bjapi "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type GameData struct {
	Dealer
	api        bjapi.BlackjackAPI
	deck       *decks.Deck
	discarded  []decks.Card
	players    []Player
	playerTurn int
}

func New(decksCnt, playersCnt int, api bjapi.BlackjackAPI) *GameData {
	var err error

	data := GameData{}
	data.deck, err = decks.New(decks.WithDecks(3), decks.Shuffle())
	if err != nil {
		log.Fatalf("failed to initialize deck: %s", err)
	}

	data.api = api
	for i := 0; i < playersCnt; i++ {
		data.players = append(data.players, NewPlayer(fmt.Sprintf("Player %d", i+1)))
	}

	data.Dealer = NewDealer("Dealer")
	data.NewRound()

	return &data
}

func (gd *GameData) Draw() decks.Card {
	if gd == nil {
		log.Fatalf("blackjack data is nil")
	}

	if gd.deck == nil {
		log.Fatalf("deck is nil")
	}

	if len(gd.deck.Cards) <= 0 {
		gd.deck.Cards = append(gd.deck.Cards, gd.discarded...)
		err := gd.deck.Shuffle()
		if err != nil {
			log.Fatalf("failed to shuffle the deck: %s", err)
		}
	}

	return gd.deck.Draw()
}

func (gd *GameData) Discard(cards []decks.Card) {
	gd.discarded = append(gd.discarded, cards...)
}

func (gd *GameData) CurrentPlayer() Player {
	if gd.IsDealersTurn() {
		return gd.Dealer
	}

	return gd.players[gd.playerTurn]
}

func (gd *GameData) Players() []Player {
	return gd.players
}

func (gd *GameData) IsDealersTurn() bool {
	return gd.playerTurn >= len(gd.players)
}

func (gd *GameData) NextPlayersTurn() {
	if gd.IsDealersTurn() {
		gd.playerTurn = 0
	} else {
		gd.playerTurn++
	}
}

func (gd *GameData) NewRound() {
	gd.playerTurn = 0
	for i := range gd.players {
		gd.Discard(gd.players[i].Discard())
	}

	gd.Discard(gd.Dealer.Discard())
	gd.Dealer.Hide()
}

func (gd *GameData) API() bjapi.BlackjackAPI {
	return gd.api
}
