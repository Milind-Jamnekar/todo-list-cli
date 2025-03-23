package cmd

import (
	"todo-list/core"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todo items",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the all flag value
		isShowAll, _ := cmd.Flags().GetBool("all")
		core.List("todo.csv", isShowAll)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add the all flag
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks (including completed ones)")
}
