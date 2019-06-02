package commands

import (
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a TODO list item as complete",
	Args:  cobra.ExactArgs(1),
	Run:   runDoCmd,
}

func init() {
	rootCmd.AddCommand(doCmd)
}

func runDoCmd(cmd *cobra.Command, args []string) {
	// Do Stuff Here
}
