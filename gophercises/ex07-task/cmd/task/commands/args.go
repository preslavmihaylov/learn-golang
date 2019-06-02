package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func requiresOneIntegerArg() cobra.PositionalArgs {
	f := func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires a task id argument")
		}

		if _, err := strconv.Atoi(args[0]); err != nil {
			return fmt.Errorf("requires an integer id argument")
		}

		return nil
	}

	return cobra.PositionalArgs(f)
}
