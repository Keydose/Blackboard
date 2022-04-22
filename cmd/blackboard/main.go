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

func removeFile(path string) {
	err := os.Remove(path)
	checkError(err)
}

func removeTasksFile() {
	removeFile("../../tasks.txt")
}

func openFile(path string, writeable bool) *os.File {
	if writeable {
		tasksFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		checkError(err)
		return tasksFile
	} else {
		tasksFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0666)
		checkError(err)
		return tasksFile
	}
}

func writeLinesToTempThenSwap(lines []string) {
	tempTasksFile := OpenTempTasksFile(true)
	for _, line := range lines {
		_, err := tempTasksFile.WriteString(fmt.Sprintf("%s\n", line))
		checkError(err)
	}

	tempTasksFile.Close()

	removeTasksFile()
	os.Rename("../../tasks.tmp.txt", "../../tasks.txt")
}

func OpenTasksFile(writeable bool) *os.File {
	return openFile("../../tasks.txt", writeable)
}

func OpenTempTasksFile(writeable bool) *os.File {
	return openFile("../../tasks.tmp.txt", writeable)
}

func GetLinesFromFile(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}

	checkError(scanner.Err())

	return fileLines
}

func List() {
	tasksFile := OpenTasksFile(false)

	blackboardAscii := figure.NewFigure("Blackboard", "small", true)
	blackboardAscii.Print()
	fmt.Println("")

	fileLines := GetLinesFromFile(tasksFile)
	if len(fileLines) == 1 {
		fmt.Println("No tasks found, add some!")
		fmt.Println("Syntax: bb add <task name> -p <position>")
	} else {
		i := 1
		for _, task := range fileLines {
			fmt.Printf("(%d) %s\n", i, task)
			i++
		}
	}
}

func Add(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	tasks := OpenTasksFile(true)
	defer tasks.Close()

	_, err := tasks.WriteString(fmt.Sprintf("%s\n", args["name"].Value))
	checkError(err)

	List()
}

func Move(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
	position, err := strconv.Atoi(args["position"].Value)
	checkError(err)

	tasksFile := OpenTasksFile(false)
	tasksFileLines := GetLinesFromFile(tasksFile)
	tasksFile.Close()
	numOfTasks := len(tasksFileLines)

	if position < 1 || position > numOfTasks {
		fmt.Println("Position is out of range")
		return
	}

	id, err := strconv.Atoi(args["id"].Value)
	checkError(err)

	if id == position {
		List()
		return
	}

	task := tasksFileLines[id-1]

	if id == numOfTasks && position == 1 {
		// Move from bottom to top
		tasksFileLines = append([]string{task}, tasksFileLines[0:numOfTasks-1]...)
	} else if id == 1 && position == numOfTasks {
		// Move from top to bottom
		tasksFileLines = append(tasksFileLines[1:numOfTasks], task)
	} else {
		// Everything before task being moved, then everything after it
		tasksFileLines = append(tasksFileLines[:id-1], tasksFileLines[id:]...)
		// Buffer slice that covers the start up to position
		bufferSlice := make([]string, position)
		// Copy everything up to (and including) position into buffer
		copy(bufferSlice, tasksFileLines[:position-1])
		// Set position to task (task is now moved)
		bufferSlice[position-1] = task

		// Join start -> position, with everything after it
		tasksFileLines = append(bufferSlice, tasksFileLines[position-1:]...)
	}

	writeLinesToTempThenSwap(tasksFileLines)
	List()
}

func Remove(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	lineNumberToRemove, err := strconv.Atoi(args["id"].Value)
	checkError(err)

	tasksFile := OpenTasksFile(false)

	taskFileLines := GetLinesFromFile(tasksFile)
	tasksFile.Close()
	numOfTasks := len(taskFileLines)
	if numOfTasks > 1 && lineNumberToRemove >= numOfTasks {
		if lineNumberToRemove == numOfTasks {
			taskFileLines = []string{}
		} else {
			taskFileLines = append(taskFileLines[:lineNumberToRemove-1], taskFileLines[lineNumberToRemove:]...)
		}

		writeLinesToTempThenSwap(taskFileLines)
	}

	List()
}

func Wipe() {
	removeTasksFile()
	List()
}

// https://semver.org/
func main() {
	commando.SetExecutableName("bb").
		SetVersion("v0.2.1").
		SetDescription("Using text files under the hood, Blackboard aims to be a minimalistic task management app that focuses on what feels natural.")

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

	commando.Register("move").
		SetDescription("Move a task (by ID) to the specified position").
		SetShortDescription("Move a task").
		AddArgument("id", "id of the task to move", "").
		AddArgument("position", "position to move the task to", "").
		SetAction(Move)

	commando.Register("remove").
		SetDescription("Remove a task from the list").
		SetShortDescription("Remove a task").
		AddArgument("id", "id of the task to remove", "").
		SetAction(Remove)

	commando.Register("wipe").
		SetDescription("Wipe all tasks from the list").
		SetShortDescription("Wipe all tasks").
		SetAction(func(_ map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			Wipe()
		})

	commando.Parse(nil)
}
