// Package commands provides a command line interface for interacting with the application,
// using commands/subcommands with arguments and flags
package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a very simple TODO management CLI application",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists the currently active tasks in your TODO list",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tsks, err := tasks.List()
		if err != nil {
			log.Fatalf("list cmd failed with error: %s", err)
		}

		fmt.Println("You have the following tasks:")
		for i, t := range tsks {
			fmt.Printf("%d. %s\n", i, t.Desc)
		}
	},
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a TODO list item as complete",
	Args:  requiresOneIntegerArg(),
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := strconv.Atoi(args[0])
		err := tasks.Do(id)
		if err != nil {
			log.Fatalf("failed to do task: %s", err)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new TODO list item",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		desc := strings.Join(args, " ")
		task := tasks.New(desc)

		err := tasks.Add(task)
		if err != nil {
			log.Fatalf("failed to add task: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(addCmd)
}

// Execute starts the CLI commands engine
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
