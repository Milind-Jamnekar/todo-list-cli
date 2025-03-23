package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GetLatestIdFromCsv(fileName string) int {
	// If file doesn't exist yet, return 1 as the first ID
	if !FileExists(fileName) {
		return 1
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var lineCount int = 0

	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lineCount++
	}

	// Subtract 1 for header and add 1 for new ID
	return lineCount
}

// Helper function to find the index of the "Done" column
func FindDoneColumnIndex(header []string) int {
	for i, colName := range header {
		if strings.ToLower(colName) == "done" {
			return i
		}
	}
	return -1 // Not found
}

func GetWriter(fileName string) *csv.Writer {
	fileExists := FileExists(fileName)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)

	// Write header if it's a new file
	if !fileExists {
		writer.Write([]string{"Id", "Task", "Created", "Done"})
		writer.Flush()
	}

	return writer
}

func MarkAsComplete(fileName string, id string) error {
	// Open the CSV file
	file, err := os.Open(fileName)
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

	outputFile, err := os.Create(fileName)
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
