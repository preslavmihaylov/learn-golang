package states

import "time"

func delayedTransition(gs gameState) gameState {
	time.Sleep(time.Second * 2)
	return gs
}

func transition(gs gameState) gameState {
	return gs
}
