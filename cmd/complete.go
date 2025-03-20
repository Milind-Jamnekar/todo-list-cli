package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

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
		err := markAsComplete(id)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("Todo item with ID %s marked as complete\n", id)
		}
	},
}

func markAsComplete(id string) error {
	// Open the CSV file
	file, err := os.Open("todo.csv")
	if err != nil {
		return fmt.Errorf("no todo file found. Start adding todos first")
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading todo file: %v", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("todo file is empty")
	}

	// Find the ID column index
	header := data[0]
	idColIndex := -1
	doneColIndex := -1

	for i, colName := range header {
		if strings.ToLower(colName) == "id" {
			idColIndex = i
		}
		if strings.ToLower(colName) == "done" {
			doneColIndex = i
		}
	}

	if idColIndex == -1 {
		return fmt.Errorf("ID column not found in todo file")
	}

	if doneColIndex == -1 {
		return fmt.Errorf("Done column not found in todo file")
	}

	// Find the row with the matching ID and update the Done column
	found := false
	for i := 1; i < len(data); i++ {
		if i < len(data) && idColIndex < len(data[i]) {
			if strings.TrimSpace(data[i][idColIndex]) == strings.TrimSpace(id) {
				data[i][doneColIndex] = "true"
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("todo item with ID %s not found", id)
	}

	// Write the updated data back to the file
	file.Close() // Close the file before writing

	outputFile, err := os.Create("todo.csv")
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	err = writer.WriteAll(data)
	if err != nil {
		return fmt.Errorf("error writing to todo file: %v", err)
	}
	writer.Flush()

	return nil
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
