package states

import (
	"fmt"
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/internal/data"
)

type gameState func(data *data.GameData) gameState

func delayedTransition(gs gameState) gameState {
	time.Sleep(time.Second * 2)
	return gs
}

func transition(gs gameState) gameState {
	return gs
}

func InitState(gd *data.GameData) gameState {
	return transition(betState)
}

func betState(gd *data.GameData) gameState {
	if gd.IsDealersTurn() {
		gd.NextPlayersTurn()
		return transition(dealState)
	}

	pl := gd.CurrentPlayer()
	gd.API().Listen(api.StartBetEvent{
		PlayerName: pl.Name(),
		Balance:    pl.Balance(),
	})

	var a api.Action
	as := api.NewActions(&api.BetAction{})
	for {
		gameAPI := gd.API()
		a = gameAPI.BetTurn(as)
		if a == nil {
			continue
		}

		switch act := a.(type) {
		case *api.BetAction:
			pl.Bet(act.Bet)
			gd.API().Listen(api.BetEvent{
				PlayerName: pl.Name(),
				Bet:        act.Bet,
			})

			gd.NextPlayersTurn()
			return delayedTransition(betState)
		case *api.HelpAction:
			for _, a := range as {
				fmt.Printf("\t%s - %s\n", a.String(), a.Help())
			}
		default:
			// continue
		}
	}
}

func dealState(gd *data.GameData) gameState {
	handSize := 2

	var ev api.DealCardsEvent
	ev.Hands = make(map[string][]decks.Card)
	for i := 0; i < handSize; i++ {
		for j, p := range gd.Players() {
			c := gd.Draw()
			gd.Players()[j].Deal(c)
			ev.Hands[p.Name()] = append(ev.Hands[p.Name()], c)
		}

		c := gd.Draw()
		gd.Dealer.Deal(c)
		ev.DealerHand = append(ev.DealerHand, c)
	}

	gd.API().Listen(ev)

	return delayedTransition(playerTurnState)
}

func playerTurnState(gd *data.GameData) gameState {
	if gd.IsDealersTurn() {
		return transition(dealerTurnState)
	}

	pl := gd.CurrentPlayer()
	gd.API().Listen(api.PlayerTurnEvent{
		PlayerName: pl.Name(),
		PlayerHand: pl.Hand(),
		DealerHand: gd.Dealer.Hand(),
	})

	if pl.HasBlackjack() {
		gd.API().Listen(api.BlackjackEvent{
			PlayerName: pl.Name(),
		})

		gd.NextPlayersTurn()
		return delayedTransition(playerTurnState)
	}

	var a api.Action
	acts := api.NewActions(&api.HitAction{}, &api.StandAction{})
	if pl.CanDoubleDown() {
		acts = append([]api.Action{&api.DoubleAction{}}, acts...)
	}

	if pl.CanSplit() {
		acts = append([]api.Action{&api.SplitAction{}}, acts...)
	}

	for {
		pi := gd.API()
		a = pi.PlayerTurn(acts)
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
			gd.NextPlayersTurn()

			gd.API().Listen(api.DoubleDownEvent{
				PlayerName: pl.Name(),
				Card:       c,
			})

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
			for _, a := range acts {
				fmt.Printf("\t%s - %s\n", a.String(), a.Help())
			}
		default:
			// continue
		}
	}
}

func dealerTurnState(gd *data.GameData) gameState {
	var ev api.DealerTurnEvent
	ev.PlayersInGame = make(map[string][]decks.Card)

	for _, p := range gd.Players() {
		if !p.Busted() {
			ev.PlayersInGame[p.Name()] = p.Hand()
		}
	}

	if !gd.Dealer.Revealed() {
		ev.DealerRevealed = true
		gd.Dealer.Reveal()
	}

	ev.DealerHand = gd.Dealer.Hand()

	gd.API().Listen(ev)
	if dealerShouldPlay(gd) {
		return delayedTransition(hitState)
	} else {
		return delayedTransition(standState)
	}
}

func resolveState(gd *data.GameData) gameState {
	var ev api.ResolveEvent
	ev.Results = make(map[string]api.Result)

	for _, p := range gd.Players() {
		res := api.Result{
			PlayerScore: p.Score(),
			DealerScore: gd.Dealer.Score(),
		}

		if p.Busted() {
			res.Outcome = api.Busted
		} else if gd.Dealer.Busted() {
			res.Outcome = api.DealerBusted
		} else if p.HasBlackjack() {
			res.Outcome = api.PlayerBlackjack
		} else if gd.Dealer.Score() < p.Score() {
			res.Outcome = api.Won
		} else if gd.Dealer.Score() > p.Score() {
			res.Outcome = api.Lost
		} else {
			res.Outcome = api.Tied
		}

		if res.Outcome == api.Won || res.Outcome == api.DealerBusted {
			p.Payout(1)
		} else if res.Outcome == api.PlayerBlackjack {
			p.Payout(1.5)
		} else if res.Outcome == api.Tied {
			p.Payout(0)
		} else {
			p.LoseBet()
		}

		ev.Results[p.Name()] = res
	}

	gd.API().Listen(ev)
	return delayedTransition(roundEndsState)
}

func roundEndsState(gd *data.GameData) gameState {
	gd.API().Listen(api.RoundEndsEvent{})
	gd.NewRound()

	return delayedTransition(betState)
}

func hitState(gd *data.GameData) gameState {
	var ev api.HitEvent

	player := gd.CurrentPlayer()
	ev.PlayerName = player.Name()

	c := gd.Draw()
	player.Deal(c)
	ev.Card = c

	if player.Busted() {
		ev.Busted = true
	}

	gd.API().Listen(ev)

	var nextState gameState
	if !gd.IsDealersTurn() {
		nextState = playerTurnState
		if player.Busted() {
			gd.NextPlayersTurn()
		}
	} else {
		nextState = dealerTurnState
		if player.Busted() {
			nextState = resolveState
		}
	}

	return delayedTransition(nextState)
}

func standState(gd *data.GameData) gameState {
	var ev api.StandEvent

	player := gd.CurrentPlayer()
	ev.PlayerName = player.Name()

	gd.API().Listen(ev)
	var nextState gameState
	if !gd.IsDealersTurn() {
		nextState = delayedTransition(playerTurnState)
	} else {
		nextState = delayedTransition(resolveState)
	}

	if !gd.IsDealersTurn() {
		gd.NextPlayersTurn()
	}

	return nextState
}

func exitState(gd *data.GameData) gameState {
	return nil
}

func dealerShouldPlay(data *data.GameData) bool {
	d := data.Dealer
	for _, p := range data.Players() {
		if !p.HasBlackjack() && !p.Busted() && d.Score() <= p.Score() &&
			(d.Score() <= 16 || (d.Score() == 17 && d.HasSoftScore())) {
			return true
		}
	}

	return false
}
