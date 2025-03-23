package cmd

import (
	"fmt"

	"todo-list/core"

	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete [ID]",
	Short: "Mark a todo item as complete",
	Long:  `Complete a todo item by its ID. For example: "complete 123" will mark the todo with ID 123 as done.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		err := core.MarkAsComplete("todo.csv", id)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("Todo item with ID %s marked as complete\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
