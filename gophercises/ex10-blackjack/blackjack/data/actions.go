package data

import (
	"fmt"
)

type Action interface {
	fmt.Stringer
	Help() string
	Do(data *GameData)
}

func NewActions(actions ...Action) []Action {
	var res []Action
	for _, a := range actions {
		res = append(res, a)
	}

	res = append(res, ExitAction{})
	res = append(res, HelpAction{res})
	return res
}

type HelpAction struct {
	actions []Action
}

func (ha HelpAction) String() string {
	return "help"
}

func (ha HelpAction) Help() string {
	return "show info about available options"
}

func (ha HelpAction) Do(data *GameData) {
	for _, a := range ha.actions {
		fmt.Printf("\t%s - %s\n", a.String(), a.Help())
	}
}

type HitAction struct{}

func (ha HitAction) String() string {
	return "hit"
}

func (ha HitAction) Help() string {
	return "draw a new card"
}

func (ha HitAction) Do(data *GameData) {
}

type StandAction struct{}

func (ha StandAction) String() string {
	return "stand"
}

func (ha StandAction) Help() string {
	return "end turn and proceed with next player"
}

func (ha StandAction) Do(data *GameData) {
}

type ExitAction struct{}

func (e ExitAction) String() string {
	return "exit"
}

func (e ExitAction) Help() string {
	return "exit the game"
}

func (e ExitAction) Do(data *GameData) {
}
