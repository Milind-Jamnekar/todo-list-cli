package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
	"testing"
	"todo-list/core"
)

func TestAddCommand(t *testing.T) {
	// Setup: Create a temporary file for testing
	tempFile := "test_todo.csv"
	defer os.Remove(tempFile) // Ensure the file is deleted after the test

	taskName := "Buy milk"
	core.Add(taskName, tempFile, core.GetWriter(tempFile))

	// Verify: Check if the task was added to the file
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// Read the file content
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	content := buf.String()

	if !bytes.Contains([]byte(content), []byte(taskName)) {
		t.Fatalf("Task '%s' was not added to the file", taskName)
	}
}

func TestListCommand(t *testing.T) {
	// Setup: Create a temporary file for testing
	tempFile := "test_todo.csv"
	defer os.Remove(tempFile) // Ensure the file is deleted after the test

	// Add some tasks to the file
	core.Add("Task 1", tempFile, core.GetWriter(tempFile))
	core.Add("Task 2", tempFile, core.GetWriter(tempFile))

	// Redirect stdout for testing
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the List command
	core.List(tempFile, false)

	// Capture the output
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Restore stdout
	os.Stdout = old

	// Verify: Check if the tasks are listed
	if !bytes.Contains([]byte(output), []byte("Task 1")) || !bytes.Contains([]byte(output), []byte("Task 2")) {
		t.Log(output)
		t.Fatalf("Tasks were not listed")
	}
}

func TestDoneCommand(t *testing.T) {
	// Setup: Create a temporary file for testing
	tempFile := "test_todo.csv"
	defer os.Remove(tempFile) // Ensure the file is deleted after the test

	// Get writer for the temporary file
	writer := core.GetWriter(tempFile)

	// Add some tasks to the file
	core.Add("Task 1", tempFile, writer)
	core.Add("Task 2", tempFile, writer)

	// Call the Done command
	err := core.MarkAsComplete(tempFile, "1")
	if err != nil {
		t.Fatalf("Failed to mark task as complete: %v", err)
	}

	// Verify: Check if the task is marked as done in the file
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV content: %v", err)
	}

	// Find the header row to identify column positions
	if len(records) == 0 {
		t.Fatalf("CSV file is empty")
	}

	// Use the existing helper function to find the Done column
	header := records[0]
	doneColIndex := core.FindDoneColumnIndex(header)
	if doneColIndex == -1 {
		t.Fatalf("Done column not found in CSV header")
	}

	// Find the ID column index
	idColIndex := -1
	for i, colName := range header {
		if strings.ToLower(colName) == "id" {
			idColIndex = i
			break
		}
	}
	if idColIndex == -1 {
		t.Fatalf("ID column not found in CSV header")
	}

	// Find the task with ID = 1 and check if it's marked as done
	taskFound := false
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) <= idColIndex || len(row) <= doneColIndex {
			continue // Skip malformed rows
		}

		if row[idColIndex] == "1" {
			taskFound = true
			if row[doneColIndex] != "true" {
				t.Errorf("Task was not marked as done. Expected 'true' in Done column, got '%s'", row[doneColIndex])
			}
			break
		}
	}

	if !taskFound {
		t.Fatalf("Task with ID 1 not found in the CSV")
	}

	// Additional test: Verify that task with ID = 2 is still not done
	task2Found := false
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) <= idColIndex || len(row) <= doneColIndex {
			continue
		}

		if row[idColIndex] == "2" {
			task2Found = true
			if row[doneColIndex] == "true" {
				t.Errorf("Task 2 was incorrectly marked as done")
			}
			break
		}
	}

	if !task2Found {
		t.Fatalf("Task with ID 2 not found in the CSV")
	}
}

func TestDoneCommandWithNonExistentId(t *testing.T) {
	// Setup test file
	tempFile := "test_todo.csv"
	defer os.Remove(tempFile)

	writer := core.GetWriter(tempFile)
	core.Add("Task 1", tempFile, writer)

	// Try to mark a non-existent task as done
	err := core.MarkAsComplete(tempFile, "999")

	// Should return an error
	if err == nil {
		t.Fatal("Expected error when marking non-existent task as done, but got nil")
	}
}
