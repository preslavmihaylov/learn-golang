package data

import (
	"fmt"
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	bjapi "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

const (
	HandSize = 2
)

type GameData struct {
	Dealer
	api        bjapi.BlackjackAPI
	deck       *decks.Deck
	decksCnt   int
	discarded  []decks.Card
	players    []Player
	playerTurn int
}

func New(decksCnt, playersCnt int, api bjapi.BlackjackAPI) *GameData {
	var err error

	data := GameData{}
	data.deck, err = decks.New(decks.WithDecks(decksCnt), decks.Shuffle())
	if err != nil {
		log.Fatalf("failed to initialize deck: %s", err)
	}

	data.decksCnt = decksCnt
	data.api = api
	for i := 0; i < playersCnt; i++ {
		data.players = append(data.players, NewPlayer(fmt.Sprintf("Player %d", i+1)))
	}

	data.Dealer = NewDealer("Dealer")
	data.NewRound()

	return &data
}

func (gd *GameData) Draw() decks.Card {
	if gd.deck == nil {
		log.Fatalf("deck is nil")
	}

	return gd.deck.Draw()
}

func (gd *GameData) ShouldShuffle() bool {
	return float64(len(gd.deck.Cards)) <= float64(decks.DeckSize)*1.5
}

func (gd *GameData) Shuffle() {
	gd.deck.InsertBottom(gd.discarded)
	gd.discarded = []decks.Card{}
	err := gd.deck.Shuffle()
	if err != nil {
		log.Fatalf("failed to shuffle the deck: %s", err)
	}
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

func (gd *GameData) SplitCurrentPlayer() {
	pl := gd.players[gd.playerTurn]
	p1, p2 := pl.Split()

	gd.players[gd.playerTurn] = p1

	// insert second split player
	gd.players = append(gd.players[:gd.playerTurn+1],
		append([]Player{p2}, gd.players[gd.playerTurn+1:len(gd.players)]...)...)
}

func (gd *GameData) NewRound() {
	gd.playerTurn = 0
	for i := range gd.players {
		gd.Discard(gd.players[i].Discard())
	}

	for i := 0; i < len(gd.players); i++ {
		if gd.players[i].IsSplit() {
			if gd.players[i].IsSplit() && i+1 >= len(gd.players) {
				log.Fatalf("Found single split player at the end of players slice")
			} else if !gd.players[i+1].IsSplit() {
				log.Fatalf("Found split player without its other half")
			}

			gd.players[i].Unsplit(gd.players[i+1])
			gd.players = append(gd.players[:i+1], gd.players[i+2:]...)
		}
	}

	gd.Discard(gd.Dealer.Discard())
	gd.Dealer.Hide()
}

func (gd *GameData) API() bjapi.BlackjackAPI {
	return gd.api
}
