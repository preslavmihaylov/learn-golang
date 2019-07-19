package states

import (
	"time"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/data"
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

	var ev data.DealCardsEvent
	for i := 0; i < handSize; i++ {
		for i := range gd.Players() {
			c := gd.Draw()
			gd.Players()[i].Deal(c)
			ev.PlayerHand = append(ev.PlayerHand, c)
		}

		c := gd.Draw()
		gd.Dealer.Deal(c)
		ev.DealerHand = append(ev.DealerHand, c)
	}

	data.EmitEvent(gd.Players(), ev)

	return delayedTransition(playerTurnState)
}

func playerTurnState(gd *data.GameData) gameState {
	if gd.IsDealersTurn() {
		return transition(dealerTurnState)
	}

	player := gd.CurrentPlayer()
	data.EmitEvent(gd.Players(), data.PlayerTurnEvent{
		PlayerName: player.Name(),
		PlayerHand: player.Hand(),
		DealerHand: gd.Dealer.Hand(),
	})

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
			return delayedTransition(hitState)
		case data.StandAction:
			return delayedTransition(standState)
		case data.ExitAction:
			return transition(exitState)
		default:
			// continue
		}
	}
}

func dealerTurnState(gd *data.GameData) gameState {
	var ev data.DealerTurnEvent
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

	data.EmitEvent(gd.Players(), ev)
	if dealerShouldPlay(gd) {
		return delayedTransition(hitState)
	} else {
		return delayedTransition(standState)
	}
}

func resolveState(gd *data.GameData) gameState {
	var ev data.ResolveEvent
	ev.Results = make(map[string]data.Result)

	for _, p := range gd.Players() {
		res := data.Result{
			PlayerScore: p.Score(),
			DealerScore: gd.Dealer.Score(),
		}

		if p.Busted() {
			res.Outcome = data.Busted
		} else if gd.Dealer.Busted() {
			res.Outcome = data.DealerBusted
		} else if gd.Dealer.Score() < p.Score() {
			res.Outcome = data.Won
		} else if gd.Dealer.Score() > p.Score() {
			res.Outcome = data.Lost
		} else {
			res.Outcome = data.Tied
		}

		ev.Results[p.Name()] = res
	}

	data.EmitEvent(gd.Players(), ev)
	return delayedTransition(roundEndsState)
}

func roundEndsState(gd *data.GameData) gameState {
	data.EmitEvent(gd.Players(), data.RoundEndsEvent{})
	gd.NewRound()

	return delayedTransition(dealState)
}

func hitState(gd *data.GameData) gameState {
	var ev data.HitEvent

	player := gd.CurrentPlayer()
	ev.PlayerName = player.Name()

	c := gd.Draw()
	player.Deal(c)
	ev.Card = c

	if player.Busted() {
		ev.Busted = true
		gd.NextPlayersTurn()
	}

	data.EmitEvent(gd.Players(), ev)

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
	var ev data.StandEvent

	player := gd.CurrentPlayer()
	ev.PlayerName = player.Name()

	data.EmitEvent(gd.Players(), ev)
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
