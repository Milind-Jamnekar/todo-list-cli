/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"todo-list/core"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:       "add",
	Short:     "add <task_name> command to add a new task",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"Taks name"},
	Run: func(cmd *cobra.Command, args []string) {
		fileName := "todo.csv"
		core.Add(args[0], fileName, core.GetWriter(fileName))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// writer := getWriter()
	// add(data, writer)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
