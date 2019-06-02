package commands

import (
	"log"
	"strings"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks"
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
	desc := strings.Join(args, " ")
	task := tasks.New(desc)

	err := tasks.Add(task)
	if err != nil {
		log.Fatalf("failed to add task: %s", err)
	}
}
