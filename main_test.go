package main

import {
	"testing"
	"todo-list/cmd"
}

func TestHello(t *testing.T) {

	want := "Hello Golang"

	got := List()

	if want != got {
		t.Fatalf("want %s, got %s\n", want, got)
	}
}
