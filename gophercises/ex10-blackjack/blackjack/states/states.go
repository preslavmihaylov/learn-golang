package states

import (
	"fmt"
	"time"

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

func InitState(gd *data.GameData) GameState {
	return Transition(DealState)
}

func DealState(gd *data.GameData) GameState {
	fmt.Println("--- Dealing cards...")
	handSize := 2
	for i := 0; i < handSize; i++ {
		for i := range gd.Players() {
			gd.Players()[i].Deal(gd.Draw())
		}

		gd.Dealer.Deal(gd.Draw())
	}

	return DelayedTransition(PlayerTurnState)
}

func PlayerTurnState(gd *data.GameData) GameState {
	if gd.IsDealersTurn() {
		return Transition(DealerTurnState)
	}

	player := gd.CurrentPlayer()
	fmt.Printf("--- %s's turn\n", player.Name())
	printTurnInfo(gd)

	var a data.Action
	as := data.NewActions(data.HitAction{}, data.StandAction{})
	for {
		pi := player.Interface()
		a = pi.PlayerTurn(data.GameInfo{
			PlayerHand: player.Hand(),
			DealerHand: gd.Dealer.Hand(),
		}, as)
		if a == nil {
			continue
		}

		a.Do(gd)
		switch a.(type) {
		case data.HitAction:
			return DelayedTransition(HitState)
		case data.StandAction:
			return DelayedTransition(StandState)
		case data.ExitAction:
			return Transition(ExitState)
		default:
			// continue
		}
	}
}

func DealerTurnState(gd *data.GameData) GameState {
	fmt.Printf("--- %s's turn\n", gd.Dealer.Name())
	fmt.Println("Players still in game:")
	for _, p := range gd.Players() {
		if !p.Busted() {
			printPlayerInfo(gd, p)
		}
	}

	if !gd.Dealer.Revealed() {
		fmt.Printf("%s reveals hand\n", gd.Dealer.Name())
		gd.Dealer.Reveal()
	}

	printPlayerInfo(gd, gd.Dealer)
	if dealerShouldPlay(gd) {
		return DelayedTransition(HitState)
	} else {
		return DelayedTransition(StandState)
	}
}

func ResolveState(gd *data.GameData) GameState {
	fmt.Println("--- Resolution")
	for _, p := range gd.Players() {
		if p.Busted() {
			fmt.Printf("%s Busted!\n", p.Name())
			continue
		} else if gd.Dealer.Busted() {
			fmt.Printf("%s Busted. %s wins!\n", gd.Dealer.Name(), p.Name())
			continue
		}

		fmt.Printf("%s's Score: %d, %s's Score: %d\n",
			p.Name(), p.Score(), gd.Dealer.Name(), gd.Dealer.Score())

		fmt.Printf("%s ", p.Name())
		if gd.Dealer.Score() < p.Score() {
			fmt.Println("won!")
		} else if gd.Dealer.Score() > p.Score() {
			fmt.Println("lost!")
		} else {
			fmt.Println("tied!")
		}
	}

	return DelayedTransition(RoundEndsState)
}

func RoundEndsState(gd *data.GameData) GameState {
	fmt.Println("--- Turn Ends")
	gd.NewRound()

	return DelayedTransition(DealState)
}

func HitState(gd *data.GameData) GameState {
	player := gd.CurrentPlayer()
	fmt.Printf("--- %s hits!\n", player.Name())
	c := gd.Draw()
	player.Deal(c)

	fmt.Printf("Got %s\n", c)
	if player.Busted() {
		fmt.Printf("%s Busted!\n", player.Name())
		gd.NextPlayersTurn()
	}

	var nextState GameState
	if !gd.IsDealersTurn() {
		nextState = DelayedTransition(PlayerTurnState)
		if player.Busted() {
			nextState = DelayedTransition(DealerTurnState)
		}
	} else {
		nextState = DelayedTransition(DealerTurnState)
		if player.Busted() {
			nextState = DelayedTransition(ResolveState)
		}
	}

	return nextState
}

func StandState(gd *data.GameData) GameState {
	player := gd.CurrentPlayer()
	fmt.Printf("--- %s stands.\n", player.Name())

	var nextState GameState
	if !gd.IsDealersTurn() {
		nextState = DelayedTransition(PlayerTurnState)
	} else {
		nextState = DelayedTransition(ResolveState)
	}

	gd.NextPlayersTurn()
	return nextState
}

func ExitState(gd *data.GameData) GameState {
	return nil
}

func printTurnInfo(gd *data.GameData) {
	if gd.IsDealersTurn() {
		return
	}

	p := gd.CurrentPlayer()
	printPlayerInfo(gd, p)
	printPlayerInfo(gd, gd.Dealer)
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
