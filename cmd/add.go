/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:       "add",
	Short:     "add <task_name> command to add a new task",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"Taks name"},
	Run: func(cmd *cobra.Command, args []string) {
		add(args[0], getWriter())
	},
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getLatestIdFromCsv() int {
	// If file doesn't exist yet, return 1 as the first ID
	if !fileExists("todo.csv") {
		return 1
	}

	file, err := os.Open("todo.csv")
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

func add(data string, writer *csv.Writer) {
	id := getLatestIdFromCsv()
	now := time.Now().Format(time.RFC3339)
	writer.Write([]string{strconv.Itoa(id), data, now, "false"})
	writer.Flush()
}

func getWriter() *csv.Writer {
	fileExists := fileExists("todo.csv")

	file, err := os.OpenFile("todo.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
