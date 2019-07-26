package states

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/internal/data"
)

func betActions(gd *data.GameData) GameState {
	acts := api.NewActions(&api.BetAction{})
	for {
		a := gd.API().BetTurn(acts)
		if a == nil {
			continue
		}

		switch act := a.(type) {
		case *api.BetAction:
			pl := gd.CurrentPlayer()
			pl.Bet(act.Bet)
			gd.API().Listen(api.BetEvent{
				PlayerName: pl.Name(),
				Bet:        act.Bet,
			})

			gd.NextPlayersTurn()
			return delayedTransition(betState)
		case *api.HelpAction:
			printHelp(acts)
		default:
			// continue
		}
	}
}

func playerActions(gd *data.GameData) GameState {
	acts := api.NewActions(&api.HitAction{}, &api.StandAction{})

	pl := gd.CurrentPlayer()
	if pl.CanDoubleDown() {
		acts = append([]api.Action{&api.DoubleAction{}}, acts...)
	}

	if pl.CanSplit() {
		acts = append([]api.Action{&api.SplitAction{}}, acts...)
	}

	for {
		a := gd.API().PlayerTurn(acts)
		if a == nil {
			continue
		}

		switch a.(type) {
		case *api.HitAction:
			return delayedTransition(hitState)
		case *api.StandAction:
			return delayedTransition(standState)
		case *api.DoubleAction:
			c := gd.Draw()
			pl.Deal(c)
			pl.DoubleDown()
			gd.API().Listen(api.DoubleDownEvent{
				PlayerName: pl.Name(),
				Card:       c,
			})

			gd.NextPlayersTurn()
			return delayedTransition(playerTurnState)
		case *api.SplitAction:
			gd.SplitCurrentPlayer()
			gd.API().Listen(api.SplitEvent{
				PlayerName: pl.Name(),
			})

			return delayedTransition(playerTurnState)
		case *api.ExitAction:
			return transition(exitState)
		case *api.HelpAction:
			printHelp(acts)
		default:
			// continue
		}
	}
}

func dealerShouldPlay(data *data.GameData) bool {
	d := data.Dealer
	for _, p := range data.Players() {
		if !p.HasBlackjack() && !p.IsBusted() && d.Score() <= p.Score() &&
			(d.Score() <= 16 || (d.Score() == 17 && d.HasSoftScore())) {
			return true
		}
	}

	return false
}

func printHelp(acts []api.Action) {
	for _, a := range acts {
		fmt.Printf("\t%s - %s\n", a.String(), a.Help())
	}
}
