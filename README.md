![Blackboard Header Image](assets/github_banner.png)

# Blackboard

Using text files under the hood, Blackboard aims to be a minimalistic task management app that focuses on what feels natural.
If you find yourself needing a place to quickly jot down things that need to be done, as well as keep track of what order you should do them in, this is probably the tool for you.

Blackboard is opinionated and doesn't try to be something that it isn't:
- No task priorities
- No due dates
- No urgency flags

When used in conjunction with a more elaborate task management system such as Jira, Trello, Asana or ClickUp (to name a few), this satisfies that gap in between: where tasks aren't quite within scope of those boards, but still need to get done.

## Highlights

- Add tasks to your list, at any position
- Bump tasks to the top of your list
- Slump tasks to the bottom of your list
- Move tasks to specific positions
- Remove tasks from the list when completed
- Wipe all tasks

## Prerequisites

- [Go 1.11.1 or greater](https://go.dev/doc/install)

## Installation

Execute the following command in your terminal, this will install Blackboard globally.

```
GO111MODULE=off go get -u github.com/Keydose/Blackboard/cmd/bb
```

Alternatively you can download the bb.exe binary from the [releases](https://github.com/Keydose/Blackboard/releases), however you will only be able to use this within the present directory, unless you add it to your PATH.

## Usage/Examples

### Add a task

```
# Add a task to the bottom of the list
bb add "Get back to John re. our upcoming API integration"

# Add a task in position 4
bb add "Schedule a meeting that should just be an email" -p 4
```

### Edit a task

```
# Edit task 3
bb edit 3 "Actually, I need to reply to their message"
```

### Bump a task to the top

```
# Bump task 3 to the top
bb bump 3
```

### Slump a task to the bottom

```
# Slump task 5 to the bottom
bb slump 5
```

### Move a task to a certain position

```
# Move task 7 to position 2
bb move 7 2
```

### Remove a task

```
# Remove task 9
bb remove 9
```

### Wipe all tasks

```
# Wipe tasks by deleting tasks.txt
bb wipe
```

## Authors

- [@keydose](https://www.github.com/keydose)

## Acknowledgements

- [Inspired by Taskbook](https://github.com/klaussinani/taskbook)

## Highlights

- Add tasks to your list, at any position
- Bump tasks to the top of your list
- Slump tasks to the bottom of your list
- Move tasks to specific positions
- Remove tasks from the list when completed

## Installation

TBD - Likely to be on npm
  
## Authors

- [@keydose](https://www.github.com/keydose)


## Acknowledgements

 - [Inspired by Taskbook](https://github.com/klaussinani/taskbook)

## Contributing

Contributions are always welcome!

As this is just a hobby project for learning Go, I am almost certain that there will be plenty of ways to make efficiency improvements, or tidy things up. 

