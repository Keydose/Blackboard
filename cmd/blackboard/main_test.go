package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/thatisuday/commando"
)

func TestAdd(t *testing.T) {
	Add(map[string]commando.ArgValue{"name": {Value: "test"}}, map[string]commando.FlagValue{})
	file := OpenTasksFile(false)
	got := GetLinesFromFile(file)
	fmt.Printf("%v", got)
	want := []string{"test"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
