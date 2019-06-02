package commands

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists the currently active tasks in your TODO list",
	Args:  cobra.NoArgs,
	Run:   runListCmd,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runListCmd(cmd *cobra.Command, args []string) {
	// Do Stuff Here
}
