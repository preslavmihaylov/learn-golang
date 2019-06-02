package commands

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new TODO list item",
	Args:  cobra.MinimumNArgs(1),
	Run:   runAddCmd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAddCmd(cmd *cobra.Command, args []string) {
	// Do Stuff Here
}
