package api

type BlackjackAPI interface {
	Listen(e GameEvent)
	PlayerTurn(actions []Action) Action
}
