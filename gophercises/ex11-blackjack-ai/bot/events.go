package bot

import (
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

func (bb *BlackjackStrategy) Listen(e api.GameEvent) {
	switch ev := e.(type) {
	case api.DeckShuffledEvent:
		bb.counting.reset()

		log.Println("[DeckShuffledEvent]: Deck got shuffled. Resetting card counting")
	case api.DealCardsEvent:
		bb.counting.runCount(ev.DealerHand...)
		for _, hand := range ev.Hands {
			bb.counting.runCount(hand...)
		}
	case api.StartBetEvent:
		bb.Stats.Balance = ev.Balance

		log.Printf("[StartBetEvent]: Balance=%d", ev.Balance)
	case api.PlayerTurnEvent:
		bb.playerHand = ev.PlayerHand
		bb.dealerHand = ev.DealerHand

		log.Printf("[PlayerTurnEvent]: Received %s %s", bb.playerHand, bb.dealerHand)
	case api.DealerTurnEvent:
		if ev.DealerRevealed {
			bb.counting.runCount(ev.DealerHand[len(ev.DealerHand)-1])
		}

		log.Printf("[DealerTurnEvent] Dealer's hand: %s", ev.DealerHand)
	case api.DoubleDownEvent:
		bb.counting.runCount(ev.Card)

		log.Printf("[DoubleDownEvent]: %s Received %s", ev.PlayerName, ev.Card)
	case api.HitEvent:
		bb.counting.runCount(ev.Card)

		log.Printf("[HitEvent]: %s hits. Received %s", ev.PlayerName, ev.Card)
		if ev.Busted {
			log.Printf("[HitEvent]: %s busted!", ev.PlayerName)
		}
	case api.StandEvent:
		log.Printf("[StandEvent] %s stands", ev.PlayerName)
	case api.BlackjackEvent:
		log.Printf("[PlayerBlackjack] %s Got Blackjack!", ev.PlayerName)
	case api.ResolveEvent:
		for _, res := range ev.Results {
			if res.Outcome == api.PlayerBlackjack ||
				res.Outcome == api.DealerBusted ||
				res.Outcome == api.Won {

				bb.Stats.HandsWon++
			} else if res.Outcome == api.Tied {
				bb.Stats.HandsTied++
			} else {
				bb.Stats.HandsLost++
			}

			log.Printf("[ResolveEvent]: Outcome=%s", res.Outcome)
		}
	case api.RoundEndsEvent:
		bb.roundsCnt--
		bb.hasSplit = false

		log.Printf("[RoundEndsEvent]: %d rounds left", bb.roundsCnt)
	}
}
