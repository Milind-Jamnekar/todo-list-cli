/*
Copyright © 2025 Milind Jamnekar HERE milind.jamnekar@yahoo.com
*/
package main

import (
	"todo-list/cmd"
)

type Task struct {
	Id   int
	Task int
	Date string
	Done bool
}

func main() {
	cmd.Execute()
}
