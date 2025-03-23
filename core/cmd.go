package core

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

// Add a new task to the CSV file
// Example: core.Add("Buy milk", core.GetWriter())
func Add(data string, fileName string, writer *csv.Writer) {
	id := GetLatestIdFromCsv(fileName)
	now := time.Now().Format(time.RFC3339)
	data = strings.TrimSpace(data)
	writer.Write([]string{strconv.Itoa(id), data, now, "false"})
	writer.Flush()
}

func List(fileName string, showAll bool) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("There is no todo csv file in your system. Start adding todos first.")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("There is a problem reading csv file. File might be corrupted.")
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
		doneColIndex := FindDoneColumnIndex(data[0])

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
