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

func openFile(path string, readable bool, writeable bool) *os.File {
	if readable && writeable {
		tasksFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		checkError(err)
		return tasksFile
	} else if readable {
		tasksFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0666)
		checkError(err)
		return tasksFile
	} else if writeable {
		tasksFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		checkError(err)
		return tasksFile
	}

	return nil
}

func writeLinesToTempThenSwap(lines []string) {
	tempTasksFile := OpenTempTasksFile(false, true)
	for _, line := range lines {
		_, err := tempTasksFile.WriteString(fmt.Sprintf("%s\n", line))
		checkError(err)
	}

	tempTasksFile.Close()

	removeTasksFile()
	os.Rename("../../tasks.tmp.txt", "../../tasks.txt")
}

func OpenTasksFile(readable bool, writeable bool) *os.File {
	return openFile("../../tasks.txt", readable, writeable)
}

func OpenTempTasksFile(readable bool, writeable bool) *os.File {
	return openFile("../../tasks.tmp.txt", readable, writeable)
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
	tasksFile := OpenTasksFile(true, false)

	blackboardAscii := figure.NewFigure("Blackboard", "small", true)
	blackboardAscii.Print()
	fmt.Println("")

	fileLines := GetLinesFromFile(tasksFile)
	if len(fileLines) == 0 {
		fmt.Println("No tasks found, add some!")
		fmt.Println("Syntax: bb add <task name> -p <position>")
	} else {
		i := 1
		for _, task := range fileLines {
			fmt.Printf("%d. %s\n", i, task)
			i++
		}
	}

	fmt.Println("")
}

func Add(name string, position int) {
	tasksFile := OpenTasksFile(false, true)
	_, err := tasksFile.WriteString(fmt.Sprintf("%s\n", name))
	tasksFile.Close()
	checkError(err)

	if position > 0 {
		tasksFile = OpenTasksFile(true, false)
		taskFileLines := GetLinesFromFile(tasksFile)
		tasksFile.Close()
		addedId := len(taskFileLines)
		Move(addedId, position)
	}
}

func Edit(id int, name string) {
	tasksFile := OpenTasksFile(true, false)
	taskFileLines := GetLinesFromFile(tasksFile)
	tasksFile.Close()

	if (id <= 0) || (id > len(taskFileLines)) {
		fmt.Println("ID is out of range")
		return
	}

	taskFileLines[id-1] = name
	writeLinesToTempThenSwap(taskFileLines)
}

func Remove(id int) {
	tasksFile := OpenTasksFile(true, false)

	taskFileLines := GetLinesFromFile(tasksFile)
	tasksFile.Close()
	numOfTasks := len(taskFileLines)
	if numOfTasks > 0 && id <= numOfTasks {
		if id == numOfTasks && numOfTasks == 1 {
			taskFileLines = []string{}
		} else {
			taskFileLines = append(taskFileLines[:id-1], taskFileLines[id:]...)
		}

		writeLinesToTempThenSwap(taskFileLines)
	}
}

func Move(id int, position int) {
	tasksFile := OpenTasksFile(true, false)
	tasksFileLines := GetLinesFromFile(tasksFile)
	tasksFile.Close()
	numOfTasks := len(tasksFileLines)

	if id == position {
		return
	}

	if position < 0 || position > numOfTasks {
		fmt.Println("Position is out of range")
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
}

func Bump(id int) {
	Move(id, 1)
}

func Slump(id int) {
	// Inefficient as has to load the file to get num of lines, then loads it again in Move function
	tasksFile := OpenTasksFile(true, false)
	tasksFileLines := GetLinesFromFile(tasksFile)
	tasksFile.Close()

	Move(id, len(tasksFileLines))
}

func Wipe() {
	removeTasksFile()
}

// https://semver.org/
func main() {
	commando.SetExecutableName("bb").
		SetVersion("v1.0.0").
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
		AddFlag("position,p", "position of the task", commando.Int, 0).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			Add(args["name"].Value, flags["position"].Value.(int))
			List()
		})

	commando.Register("edit").
		SetDescription("Edit a task (by ID)").
		SetShortDescription("Add a task").
		AddArgument("id", "id of the task to edit", "").
		AddArgument("name", "new name for the task", "").
		SetAction(func(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			idAsInt, err := strconv.Atoi(args["id"].Value)
			checkError(err)

			Edit(idAsInt, args["name"].Value)
			List()
		})

	commando.Register("remove").
		SetDescription("Remove a task from the list").
		SetShortDescription("Remove a task").
		AddArgument("id", "id of the task to remove", "").
		SetAction(func(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			id, err := strconv.Atoi(args["id"].Value)
			checkError(err)
			Remove(id)
			List()
		})

	commando.Register("move").
		SetDescription("Move a task (by ID) to the specified position").
		SetShortDescription("Move a task").
		AddArgument("id", "id of the task to move", "").
		AddArgument("position", "position to move the task to", "").
		SetAction(func(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			idAsInt, err := strconv.Atoi(args["id"].Value)
			checkError(err)
			positionAsInt, err := strconv.Atoi(args["position"].Value)
			checkError(err)

			Move(idAsInt, positionAsInt)
			List()
		})

	commando.Register("bump").
		SetDescription("Bump a task (by ID) to the top").
		SetShortDescription("Bump a task").
		AddArgument("id", "id of the task to bump", "").
		SetAction(func(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			idAsInt, err := strconv.Atoi(args["id"].Value)
			checkError(err)

			Bump(idAsInt)
			List()
		})

	commando.Register("slump").
		SetDescription("Slump a task (by ID) to the bottom").
		SetShortDescription("Slump a task").
		AddArgument("id", "id of the task to slump", "").
		SetAction(func(args map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			idAsInt, err := strconv.Atoi(args["id"].Value)
			checkError(err)

			Slump(idAsInt)
			List()
		})

	commando.Register("wipe").
		SetDescription("Wipe all tasks from the list").
		SetShortDescription("Wipe all tasks").
		SetAction(func(_ map[string]commando.ArgValue, _ map[string]commando.FlagValue) {
			Wipe()
			List()
		})

	commando.Parse(nil)
}
