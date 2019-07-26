package api

import (
	"fmt"
	"strconv"
)

type Action interface {
	fmt.Stringer
	Help() string
	ArgsCnt() int
	SetArgs(args ...string) error
}

func NewActions(actions ...Action) []Action {
	var res []Action
	for _, a := range actions {
		res = append(res, a)
	}

	res = append(res, &ExitAction{})
	res = append(res, &HelpAction{NoArgsAction{}, res})
	return res
}

type NoArgsAction struct{}

func (*NoArgsAction) ArgsCnt() int {
	return 0
}

func (*NoArgsAction) SetArgs(args ...string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many args passed")
	}

	return nil
}

type HelpAction struct {
	NoArgsAction
	actions []Action
}

func (*HelpAction) String() string {
	return "help"
}

func (*HelpAction) Help() string {
	return "show info about available options"
}

type BetAction struct {
	Bet int
}

func (*BetAction) String() string {
	return "bet"
}

func (*BetAction) Help() string {
	return "place a bet"
}

func (*BetAction) ArgsCnt() int {
	return 1
}

func (ba *BetAction) SetArgs(args ...string) error {
	if len(args) > 1 {
		return fmt.Errorf("Too many args passed")
	}

	var err error
	ba.Bet, err = strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse bet: %s", err)
	}

	return nil
}

type HitAction struct {
	NoArgsAction
}

func (*HitAction) String() string {
	return "hit"
}

func (*HitAction) Help() string {
	return "draw a new card"
}

type StandAction struct {
	NoArgsAction
}

func (*StandAction) String() string {
	return "stand"
}

func (*StandAction) Help() string {
	return "end turn and proceed with next player"
}

type DoubleAction struct {
	NoArgsAction
}

func (*DoubleAction) String() string {
	return "double"
}

func (*DoubleAction) Help() string {
	return "double the bet and play only one card more"
}

type SplitAction struct {
	NoArgsAction
}

func (*SplitAction) String() string {
	return "split"
}

func (*SplitAction) Help() string {
	return "split the current hand"
}

type ExitAction struct {
	NoArgsAction
}

func (*ExitAction) String() string {
	return "exit"
}

func (*ExitAction) Help() string {
	return "exit the game"
}
