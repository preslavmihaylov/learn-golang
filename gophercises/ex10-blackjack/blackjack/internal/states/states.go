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
	return transition(dealState)
}

func dealState(gd *data.GameData) gameState {
	handSize := 2

	var ev api.DealCardsEvent
	ev.Hands = make(map[string][]decks.Card)
	for i := 0; i < handSize; i++ {
		for i, p := range gd.Players() {
			c := gd.Draw()
			gd.Players()[i].Deal(c)
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

	player := gd.CurrentPlayer()
	gd.API().Listen(api.PlayerTurnEvent{
		PlayerName: player.Name(),
		PlayerHand: player.Hand(),
		DealerHand: gd.Dealer.Hand(),
	})

	var a api.Action
	as := api.NewActions(api.HitAction{}, api.StandAction{})
	for {
		pi := gd.API()
		a = pi.PlayerTurn(as)
		if a == nil {
			continue
		}

		switch a.(type) {
		case api.HitAction:
			return delayedTransition(hitState)
		case api.StandAction:
			return delayedTransition(standState)
		case api.ExitAction:
			return transition(exitState)
		case api.HelpAction:
			for _, a := range as {
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
		} else if gd.Dealer.Score() < p.Score() {
			res.Outcome = api.Won
		} else if gd.Dealer.Score() > p.Score() {
			res.Outcome = api.Lost
		} else {
			res.Outcome = api.Tied
		}

		ev.Results[p.Name()] = res
	}

	gd.API().Listen(ev)
	return delayedTransition(roundEndsState)
}

func roundEndsState(gd *data.GameData) gameState {
	gd.API().Listen(api.RoundEndsEvent{})
	gd.NewRound()

	return delayedTransition(dealState)
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
		gd.NextPlayersTurn()
	}

	gd.API().Listen(ev)

	var nextState gameState
	if !gd.IsDealersTurn() {
		nextState = delayedTransition(playerTurnState)
		if player.Busted() {
			nextState = delayedTransition(dealerTurnState)
		}
	} else {
		nextState = delayedTransition(dealerTurnState)
		if player.Busted() {
			nextState = delayedTransition(resolveState)
		}
	}

	return nextState
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

	gd.NextPlayersTurn()
	return nextState
}

func exitState(gd *data.GameData) gameState {
	return nil
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
