package states

import (
	"fmt"
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/actions"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/data"
)

type GameState func(data *data.GameData) GameState

func DelayedTransition(gs GameState) GameState {
	time.Sleep(time.Second * 2)
	return gs
}

func Transition(gs GameState) GameState {
	return gs
}

func InitState(data *data.GameData) GameState {
	return Transition(DealState)
}

func DealState(data *data.GameData) GameState {
	fmt.Println("--- Dealing cards...")
	handSize := 2
	for i := 0; i < handSize; i++ {
		for i := range data.Players() {
			data.Players()[i].Deal(data.Draw())
		}

		data.Dealer.Deal(data.Draw())
	}

	return DelayedTransition(PlayerTurnState)
}

func PlayerTurnState(data *data.GameData) GameState {
	if data.IsDealersTurn() {
		return Transition(DealerTurnState)
	}

	player := data.CurrentPlayer()
	fmt.Printf("--- %s's turn\n", player.Name())
	printTurnInfo(data)

	var a actions.Action
	as := actions.NewActions(actions.HitAction{}, actions.StandAction{})
	for {
		a = actions.Prompt(data, as)
		if a == nil {
			continue
		}

		a.Do(data)
		switch a.(type) {
		case actions.HitAction:
			return DelayedTransition(HitState)
		case actions.StandAction:
			return DelayedTransition(StandState)
		case actions.ExitAction:
			return Transition(ExitState)
		default:
			// continue
		}
	}
}

func DealerTurnState(data *data.GameData) GameState {
	fmt.Printf("--- %s's turn\n", data.Dealer.Name())
	fmt.Println("Players still in game:")
	for _, p := range data.Players() {
		if !p.Busted() {
			printPlayerInfo(data, p)
		}
	}

	if !data.Dealer.Revealed() {
		fmt.Printf("%s reveals hand\n", data.Dealer.Name())
		data.Dealer.Reveal()
	}

	printPlayerInfo(data, data.Dealer)
	if dealerShouldPlay(data) {
		return DelayedTransition(HitState)
	} else {
		return DelayedTransition(StandState)
	}
}

func ResolveState(data *data.GameData) GameState {
	fmt.Println("--- Resolution")
	for _, p := range data.Players() {
		if p.Busted() {
			fmt.Printf("%s Busted!\n", p.Name())
			continue
		} else if data.Dealer.Busted() {
			fmt.Printf("%s Busted. %s wins!\n", data.Dealer.Name(), p.Name())
			continue
		}

		fmt.Printf("%s's Score: %d, %s's Score: %d\n",
			p.Name(), p.Score(), data.Dealer.Name(), data.Dealer.Score())

		fmt.Printf("%s ", p.Name())
		if data.Dealer.Score() < p.Score() {
			fmt.Println("won!")
		} else if data.Dealer.Score() > p.Score() {
			fmt.Println("lost!")
		} else {
			fmt.Println("tied!")
		}
	}

	return DelayedTransition(RoundEndsState)
}

func RoundEndsState(data *data.GameData) GameState {
	fmt.Println("--- Turn Ends")
	data.NewRound()

	return DelayedTransition(DealState)
}

func HitState(data *data.GameData) GameState {
	player := data.CurrentPlayer()
	fmt.Printf("--- %s hits!\n", player.Name())
	c := data.Draw()
	player.Deal(c)

	fmt.Printf("Got %s\n", c)
	if player.Busted() {
		fmt.Printf("%s Busted!\n", player.Name())
		data.NextPlayersTurn()
	}

	if !data.IsDealersTurn() {
		return DelayedTransition(PlayerTurnState)
	} else {
		return DelayedTransition(DealerTurnState)
	}
}

func StandState(data *data.GameData) GameState {
	player := data.CurrentPlayer()
	fmt.Printf("--- %s stands.\n", player.Name())
	data.NextPlayersTurn()

	if !data.IsDealersTurn() {
		return DelayedTransition(PlayerTurnState)
	} else {
		return DelayedTransition(ResolveState)
	}
}

func ExitState(data *data.GameData) GameState {
	return nil
}

func printTurnInfo(data *data.GameData) {
	if data.IsDealersTurn() {
		return
	}

	p := data.CurrentPlayer()
	printPlayerInfo(data, p)
	printPlayerInfo(data, data.Dealer)
}

func printPlayerInfo(data *data.GameData, player data.Player) {
	fmt.Printf("%s's Hand:\n", player.Name())
	for _, c := range player.Hand() {
		fmt.Printf("\t%s\n", c)
	}

	fmt.Printf("\n%s's Score: %d\n", player.Name(), player.Score())
	fmt.Println()
}

func dealerShouldPlay(data *data.GameData) bool {
	d := data.Dealer
	for _, p := range data.Players() {
		if !p.Busted() && d.Score() <= p.Score() &&
			(d.Score() <= 16 || (d.Score() == 17 && d.IsSoftScore())) {
			return true
		}
	}

	return false
}
