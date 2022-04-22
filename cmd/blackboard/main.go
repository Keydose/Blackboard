package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/common-nighthawk/go-figure"
	"github.com/thatisuday/commando"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func removeTasksFile() {
	err := os.Remove("../../tasks.txt")
	checkError(err)
}

func OpenTempFile(writeable bool) *os.File {
	if writeable {
		temp, err := os.OpenFile("../../tasks.tmp.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		checkError(err)
		return temp
	} else {
		temp, err := os.OpenFile("../../tasks.tmp.txt", os.O_CREATE|os.O_RDONLY, 0666)
		checkError(err)
		return temp
	}
}

func OpenTasksFile(writeable bool) *os.File {
	if writeable {
		tasks, err := os.OpenFile("../../tasks.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		checkError(err)
		return tasks
	} else {
		tasks, err := os.OpenFile("../../tasks.txt", os.O_CREATE|os.O_RDONLY, 0666)
		checkError(err)
		return tasks
	}
}

func GetLinesFromFile(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func List() {
	tasks := OpenTasksFile(false)
	defer tasks.Close()

	myFigure := figure.NewFigure("Blackboard", "small", true)
	myFigure.Print()
	fmt.Println("")

	tasksScanner := bufio.NewScanner(tasks)
	i := 1
	for _, task := range GetLinesFromFile(tasks) {
		fmt.Printf("[%d] %s\n", i, task)
		i++
	}

	if i == 1 {
		fmt.Println("No tasks found, add some!")
		fmt.Println("Syntax: bb add <task name> -p <position>")
	}

	checkError(tasksScanner.Err())
}

func Wipe() {
	removeTasksFile()
	List()
}

func Add(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	tasks := OpenTasksFile(true)
	defer tasks.Close()

	_, err := tasks.WriteString(fmt.Sprintf("%s\n", args["name"].Value))
	checkError(err)

	List()
}

func Remove(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	id, err := strconv.Atoi(args["id"].Value)
	checkError(err)

	tasks := OpenTasksFile(false)

	lines := GetLinesFromFile(tasks)
	tasks.Close()
	lines = append(lines[:id-1], lines[id:]...)

	tempFile := OpenTempFile(true)
	for _, line := range lines {
		_, err := tempFile.WriteString(fmt.Sprintf("%s\n", line))
		checkError(err)
	}

	tempFile.Close()

	removeTasksFile()
	os.Rename("../../tasks.tmp.txt", "../../tasks.txt")

	List()
}

// https://semver.org/
func main() {
	commando.SetExecutableName("bb").
		SetVersion("v0.1.0").
		SetDescription("A minimalistic CLI task list app - just move it to the top if it's more urgent!")

	commando.Register("list").
		SetDescription("List all tasks").
		SetShortDescription("List all tasks").
		SetAction(func(_ map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			List()
		})

	commando.Register("add").
		SetDescription("Add a task to the list of tasks").
		SetShortDescription("Add a task").
		AddArgument("name", "name of the task to create", "").
		AddFlag("position,p", "position of the task", commando.Int, 1).
		SetAction(Add)

	commando.Register("remove").
		SetDescription("Remove a task from the list").
		SetShortDescription("Remove a task").
		AddArgument("id", "id of the task to remove", "").
		SetAction(Remove)

	commando.Parse(nil)
}
