package commands

import (
	"fmt"
	"strconv"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a TODO list item as complete",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires a task id argument")
		}

		if _, err := strconv.Atoi(args[0]); err != nil {
			return fmt.Errorf("requires an integer id argument")
		}

		return nil
	},
	Run: runDoCmd,
}

func init() {
	rootCmd.AddCommand(doCmd)
}

func runDoCmd(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	err := tasks.Do(id)
	if err != nil {
		fmt.Println(err)
	}
}
