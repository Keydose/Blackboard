package main

import (
	"fmt"
	"reflect"
	"testing"
)

// TODO: Write more elaborate and more comprehensive tests that actually help - the below is more of a proof of concept

// Wipes tasks.txt, adds a task and then tests lines from file
func TestAdd(t *testing.T) {
	Wipe()
	Add("test", 0)
	file := OpenTasksFile(true, false)
	got := GetLinesFromFile(file)
	fmt.Printf("%v", got)
	want := []string{"test"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
