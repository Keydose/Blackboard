package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/thatisuday/commando"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func openTasksFile(writeable bool) *os.File {
	if writeable {
		tasks, err := os.OpenFile("../../tasks.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		checkError(err)
		return tasks
	} else {
		tasks, err := os.Open("../../tasks.txt")
		checkError(err)
		return tasks
	}
}

func list() {
	tasks := openTasksFile(false)
	defer tasks.Close()

	myFigure := figure.NewFigure("Blackboard", "small", true)
	myFigure.Print()
	fmt.Println("")

	tasksScanner := bufio.NewScanner(tasks)
	i := 1
	for tasksScanner.Scan() {
		fmt.Printf("[%d] %s\n", i, tasksScanner.Text())
		i++
	}

	checkError(tasksScanner.Err())
}

func add(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	tasks := openTasksFile(true)
	defer tasks.Close()

	// TODO: Add support for adding tasks at a specific position

	_, err := tasks.WriteString(fmt.Sprintf("%s\n", args["name"].Value))
	checkError(err)

	list()
}

// TODO: For move function, maybe we can use this lib: https://github.com/manifoldco/promptui
//       - You'd navigate the current blackboard with arrow keys, select one, then move to where you want it to be and press enter

func main() {
	commando.SetExecutableName("bb").
		SetVersion("v1.0.0").
		SetDescription("A minimalistic CLI task list app - just move it to the top if it's more urgent!")

	commando.Register("add").
		SetDescription("Add a task to the list of tasks").
		SetShortDescription("Add a task").
		AddArgument("name", "name of the task to create", "").
		AddFlag("position,p", "position of the task", commando.Int, 1).
		SetAction(add)

	commando.Register("list").
		SetDescription("List all tasks").
		SetShortDescription("List all tasks").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			list()
		})

	commando.Register("remove").
		SetDescription("Remove a task from the list").
		SetShortDescription("Remove a task")

	commando.Parse(nil)
}
