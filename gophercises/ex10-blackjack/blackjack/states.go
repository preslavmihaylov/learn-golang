package blackjack

import (
	"fmt"
	"log"
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type GameState func(data *BlackjackData) GameState

func Transition(gs GameState) GameState {
	time.Sleep(time.Second * 2)
	return gs
}

func InitState(data *BlackjackData) GameState {
	var err error
	data.deck, err = decks.New(decks.WithDecks(3), decks.Shuffle())
	if err != nil {
		log.Fatalf("failed to initialize deck: %s", err)
	}

	data.players = []Player{&player{}}
	data.dealer = &dealer{}

	return Transition(DealState)
}

func DealState(data *BlackjackData) GameState {
	fmt.Println("--- Dealing cards...")
	handSize := 2
	for i := 0; i < handSize; i++ {
		for i := range data.players {
			data.players[i].Deal(data.deck.Draw())
		}

		data.dealer.Deal(data.deck.Draw())
	}

	data.playerTurn = 0
	return Transition(PlayerTurnState)
}

func PlayerTurnState(data *BlackjackData) GameState {
	if data.IsDealersTurn() {
		return DealerTurnState
	}

	fmt.Printf("--- Player %d's turn\n", data.playerTurn+1)
	printTurnInfo(data)

	var a Action
	var nextState GameState
	actions := NewActions(HitAction{}, StandAction{})
	for a == nil || nextState == nil {
		a = Prompt(data, actions)
		if a != nil {
			nextState = a.Do(data)
		}
	}

	return Transition(nextState)
}

func DealerTurnState(data *BlackjackData) GameState {
	fmt.Println("--- Dealer's turn")
	fmt.Println("Players still in game:")
	for i, p := range data.players {
		if !p.Busted() {
			printPlayerInfo(data, i)
		}
	}

	if !data.dealer.Revealed() {
		fmt.Println("Dealer reveals hand")
		data.dealer.Reveal()
	}

	printDealerInfo(data)
	if data.dealer.Score() <= 16 || (data.dealer.Score() == 17 && data.dealer.IsSoftScore()) {
		return Transition(HitState)
	}

	return Transition(ResolveState)
}

func ResolveState(data *BlackjackData) GameState {
	fmt.Println("--- Resolution")
	for i, p := range data.players {
		if p.Busted() {
			fmt.Printf("Player %d Busted!\n", i+1)
			continue
		}

		fmt.Printf("Player %d's Score: %d, Dealer's Score: %d\n",
			i+1, p.Score(), data.dealer.Score())
		if data.dealer.Score() < p.Score() {
			fmt.Printf("Player %d won!\n", i+1)
		} else if data.dealer.Score() > p.Score() {
			fmt.Printf("Player %d lost!\n", i+1)
		} else {
			fmt.Printf("Player %d tied!\n", i+1)
		}
	}

	return Transition(RoundEndsState)
}

func RoundEndsState(data *BlackjackData) GameState {
	fmt.Println("--- Turn Ends")
	data.NewRound()

	return Transition(DealState)
}

func HitState(data *BlackjackData) GameState {
	var player Player
	var playerStr string

	if !data.IsDealersTurn() {
		player = data.players[data.playerTurn]
		playerStr = fmt.Sprintf("Player %d", data.playerTurn+1)
	} else {
		player = data.dealer
		playerStr = "Dealer"
	}

	fmt.Printf("--- %s hits!\n", playerStr)
	c := data.deck.Draw()
	player.Deal(c)

	fmt.Printf("Got %s\n", c)
	if player.Busted() {
		fmt.Printf("%s Busted!\n", playerStr)
		data.NextPlayersTurn()
	}

	return Transition(PlayerTurnState)
}

func StandState(data *BlackjackData) GameState {
	fmt.Printf("--- Player %d stands.\n", data.playerTurn+1)
	data.NextPlayersTurn()

	return Transition(PlayerTurnState)
}

func printTurnInfo(data *BlackjackData) {
	printPlayerInfo(data, data.playerTurn)
	printDealerInfo(data)
}

func printPlayerInfo(data *BlackjackData, playerIndex int) {
	player := data.players[playerIndex]
	fmt.Printf("Player %d's Hand:\n", playerIndex+1)
	for _, c := range player.Hand() {
		fmt.Printf("\t%s\n", c)
	}

	fmt.Printf("\nPlayer %d's Score: %d\n", playerIndex+1, player.Score())
	fmt.Println()
}

func printDealerInfo(data *BlackjackData) {
	fmt.Println("Dealer's Hand:")
	for _, c := range data.dealer.Hand() {
		fmt.Printf("\t%s\n", c)
	}

	fmt.Printf("\nDealer's Score: %d\n", data.dealer.Score())
	fmt.Println()
}
