package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a very simple TODO management CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute runs the entire cobra commands configuration
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
