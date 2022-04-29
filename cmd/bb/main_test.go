package main

import (
	"fmt"
	"reflect"
	"testing"
)

// Wipes tasks.txt, adds a task and then tests lines from file
func TestAdd(t *testing.T) {
	Wipe()
	Add("test", 0)
	file := OpenTasksFile(false)
	got := GetLinesFromFile(file)
	fmt.Printf("%v", got)
	want := []string{"test"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
