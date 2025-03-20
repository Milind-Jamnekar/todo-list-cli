package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todo items",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the all flag value
		all, _ := cmd.Flags().GetBool("all")
		List(all)
	},
}

func List(showAll bool) {
	file, err := os.Open("todo.csv")
	if err != nil {
		rootCmd.Println("There is no todo csv file in your system. Start adding todos first.")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		rootCmd.Println("There is a problem reading csv file. File might be corrupted.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

	// Display header if exists
	if len(data) > 0 {
		for j, value := range data[0] {
			if j == 2 { // Date column
				fmt.Fprintf(w, "Time Since\t")
			} else {
				fmt.Fprintf(w, "%s\t", value)
			}
		}
		fmt.Fprintln(w) // End the header row
	}

	// Skip header row if exists
	startRow := 1
	if len(data) <= 1 {
		startRow = 0 // No header or empty file
	}

	for i := startRow; i < len(data); i++ {
		row := data[i]

		// Check if we should display this row
		// Assuming "Done" column exists and contains "true" or "false"
		// You may need to adjust the doneColIndex based on your actual CSV structure
		doneColIndex := findDoneColumnIndex(data[0])

		if !showAll && doneColIndex >= 0 && doneColIndex < len(row) {
			// Skip completed items when showAll is false
			if strings.ToLower(row[doneColIndex]) == "true" {
				continue
			}
		}

		for j, value := range row {
			if j == 2 { // Assuming column 3 is the date
				parsedTime, err := time.Parse(time.RFC3339, value)
				if err != nil {
					fmt.Fprintf(w, "Invalid date\t")
				} else {
					timeDiff := timediff.TimeDiff(parsedTime)
					fmt.Fprintf(w, "%s\t", timeDiff)
				}
			} else {
				fmt.Fprintf(w, "%s\t", value)
			}
		}
		fmt.Fprintln(w) // End the row
	}

	w.Flush() // Flush the tabwriter to actually display the table
}

// Helper function to find the index of the "Done" column
func findDoneColumnIndex(header []string) int {
	for i, colName := range header {
		if strings.ToLower(colName) == "done" {
			return i
		}
	}
	return -1 // Not found
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add the all flag
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks (including completed ones)")
}
