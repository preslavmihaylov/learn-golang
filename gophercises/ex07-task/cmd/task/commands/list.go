package commands

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks"
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
	ts := tasks.List()
	fmt.Println("You have the following tasks:")
	for i, t := range ts {
		fmt.Printf("%d. %s\n", i, t.Desc)
	}
}
