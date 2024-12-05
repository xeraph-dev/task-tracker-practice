// Copyright 2024 xeraph. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTaskAlreadyExists       = errors.New("task already exists")
	ErrTaskDoesNotExist        = errors.New("task does not exist")
	ErrOnlyOneArgumentAllowed  = errors.New("only one argument is allowed")
	ErrOnlyTwoArgumentsAllowed = errors.New("only two arguments are allowed")
	ErrInvalidTaskStatus       = errors.New("invalid task status")
)

type TaskStatus uint8

const (
	_ TaskStatus = iota
	TaskStatusTodo
	TaskStatusInProgress
	TaskStatusDone
)

var taskStatusMapFromString = map[string]TaskStatus{
	"todo":        TaskStatusTodo,
	"in-progress": TaskStatusInProgress,
	"done":        TaskStatusDone,
}

var taskStatusMapToString = map[TaskStatus]string{
	TaskStatusTodo:       "todo",
	TaskStatusInProgress: "in-progress",
	TaskStatusDone:       "done",
}

func NewTaskStatus(str string) TaskStatus {
	return taskStatusMapFromString[str]
}

func (status TaskStatus) String() string {
	return taskStatusMapToString[status]
}

func (status TaskStatus) Valid() bool {
	_, ok := taskStatusMapToString[status]
	return ok
}

type TaskId uint64

type Task struct {
	Id          TaskId     `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStoreMeta struct {
	CurrentId uint64 `json:"current_id"`
}

type TaskStore struct {
	dbPath string
	Meta   TaskStoreMeta `json:"meta"`
	Tasks  []Task        `json:"tasks"`
}

func NewTaskStore() (store *TaskStore, err error) {
	store = new(TaskStore)
	store.Meta.CurrentId = 1
	store.Tasks = make([]Task, 0)

	var dir string
	if dir, err = os.UserConfigDir(); err != nil {
		return
	}
	store.dbPath = path.Join(dir, "task", "task.json")

	return
}

func (store *TaskStore) Load() (err error) {
	if err = os.MkdirAll(path.Dir(store.dbPath), os.ModePerm); err != nil {
		return
	}

	if _, err = os.Stat(store.dbPath); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}

	var data []byte
	if data, err = os.ReadFile(store.dbPath); err != nil {
		return
	}

	if err = json.Unmarshal(data, store); err != nil {
		return
	}

	return
}

func (store *TaskStore) Save() (err error) {
	var data []byte
	if data, err = json.Marshal(store); err != nil {
		return
	}

	if err = os.WriteFile(store.dbPath, data, os.ModePerm); err != nil {
		return
	}

	return
}

func (store *TaskStore) Create(task Task) (newTask Task, err error) {
	if store.Exists(task) {
		err = ErrTaskAlreadyExists
		return
	}

	task.Id = TaskId(store.Meta.CurrentId)
	store.Meta.CurrentId++
	task.Status = TaskStatusTodo
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt

	store.Tasks = append(store.Tasks, task)
	return task, store.Save()
}

func (store *TaskStore) Update(task Task) (err error) {
	if store.Exists(task) {
		err = ErrTaskAlreadyExists
		return
	}

	task.UpdatedAt = time.Now()

	store.Tasks[store.Index(task.Id)] = task
	return store.Save()
}

func (store *TaskStore) Delete(task Task) (err error) {
	index := store.Index(task.Id)
	store.Tasks = slices.Delete(store.Tasks, index, index+1)
	return store.Save()
}

func (store *TaskStore) Exists(task Task) bool {
	return slices.ContainsFunc(store.Tasks, func(v Task) bool {
		return v.Id != task.Id && v.Description == task.Description
	})
}

func (store *TaskStore) Index(id TaskId) int {
	return slices.IndexFunc(store.Tasks, func(task Task) bool {
		return task.Id == id
	})
}

func (store *TaskStore) GetById(id TaskId) (task Task, err error) {
	index := store.Index(id)
	if index == -1 {
		err = ErrTaskDoesNotExist
		return
	}

	task = store.Tasks[index]
	return
}

func (store *TaskStore) GetByStatus(status TaskStatus) (tasks []Task) {
	for _, task := range store.Tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return
}

type CommandState struct {
	TaskStore *TaskStore
	Args      []string
}

func NewCommandState(args []string) (state *CommandState, err error) {
	state = new(CommandState)
	state.Args = args
	state.TaskStore, err = NewTaskStore()
	return
}

func helpCommand(*CommandState) (err error) {
	fmt.Print(`USAGE: task [command] [args]

COMMANDS:
	help       show this message
	add        add a new task
	update     update a task
	delete     delete a task
	mark       change a task status
	list       list all tasks

EXAMPLES:
	task help

	task-cli add "Buy groceries"
	task-cli update 1 "Buy groceries and cook dinner"
	task-cli delete 1

	task-cli mark 1 done
	task-cli mark 1 todo
	task-cli mark 1 in-progress

	task-cli list
	task-cli list done
	task-cli list todo
	task-cli list in-progress
`)
	return
}

func addCommand(state *CommandState) (err error) {
	if len(state.Args) != 1 {
		err = ErrOnlyOneArgumentAllowed
		return
	}

	var task Task
	task.Description = state.Args[0]

	if task, err = state.TaskStore.Create(task); err != nil {
		return
	}

	fmt.Printf("Task added successfully: (ID: %d)\n", task.Id)

	return
}

func updateCommand(state *CommandState) (err error) {
	if len(state.Args) != 2 {
		err = ErrOnlyTwoArgumentsAllowed
		return
	}

	var id uint64
	if id, err = strconv.ParseUint(state.Args[0], 10, 64); err != nil {
		return
	}

	var task Task
	if task, err = state.TaskStore.GetById(TaskId(id)); err != nil {
		return
	}

	task.Description = state.Args[1]

	if err = state.TaskStore.Update(task); err != nil {
		return
	}

	fmt.Println("Task updated successfully")
	return
}

func deleteCommand(state *CommandState) (err error) {
	if len(state.Args) != 1 {
		err = ErrOnlyOneArgumentAllowed
		return
	}

	var id uint64
	if id, err = strconv.ParseUint(state.Args[0], 10, 64); err != nil {
		return
	}

	var task Task
	if task, err = state.TaskStore.GetById(TaskId(id)); err != nil {
		return
	}

	if err = state.TaskStore.Delete(task); err != nil {
		return
	}

	fmt.Println("Task deleted successfully")
	return
}

func markCommand(state *CommandState) (err error) {
	if len(state.Args) != 2 {
		err = ErrOnlyTwoArgumentsAllowed
		return
	}

	var id uint64
	if id, err = strconv.ParseUint(state.Args[0], 10, 64); err != nil {
		return
	}

	var status TaskStatus
	if status = NewTaskStatus(state.Args[1]); !status.Valid() {
		err = ErrInvalidTaskStatus
		return
	}

	var task Task
	if task, err = state.TaskStore.GetById(TaskId(id)); err != nil {
		return
	}

	task.Status = status

	if err = state.TaskStore.Update(task); err != nil {
		return
	}

	fmt.Println("Task status updated to", task.Status.String())
	return
}

func listCommand(state *CommandState) (err error) {
	if len(state.Args) > 1 {
		err = ErrOnlyOneArgumentAllowed
		return
	}

	var tasks []Task

	if len(state.Args) == 0 {
		tasks = state.TaskStore.Tasks
	} else {
		var status TaskStatus
		if status = NewTaskStatus(state.Args[0]); !status.Valid() {
			err = ErrInvalidTaskStatus
			return
		}

		tasks = state.TaskStore.GetByStatus(status)
	}

	maxStatusLen := max(len(TaskStatusTodo.String()), len(TaskStatusInProgress.String()), len(TaskStatusDone.String()))
	dateLen := len(time.Now().Format(time.DateTime))

	{
		header := strings.Builder{}
		header.WriteString("id")
		header.WriteString("    ")
		header.WriteString("status")
		header.WriteString("    " + strings.Repeat(" ", maxStatusLen-len("status")))
		header.WriteString("created at")
		header.WriteString("    " + strings.Repeat(" ", dateLen-len("created at")))
		header.WriteString("updated at")
		header.WriteString("    " + strings.Repeat(" ", dateLen-len("updated at")))
		header.WriteString("description")
		fmt.Println(header.String())
	}

	currentId := strconv.FormatUint(state.TaskStore.Meta.CurrentId, 10)

	for _, task := range tasks {
		id := strconv.FormatUint(uint64(task.Id), 10)
		status := task.Status.String()

		body := strings.Builder{}
		body.WriteString(id)
		idLen := len(currentId) - len(id)
		if idLen == 0 {
			idLen = 1
		}
		body.WriteString("    " + strings.Repeat(" ", idLen))
		body.WriteString(status)
		body.WriteString("    " + strings.Repeat(" ", maxStatusLen-len(status)))
		body.WriteString(task.CreatedAt.Format(time.DateTime))
		body.WriteString("    ")
		body.WriteString(task.UpdatedAt.Format(time.DateTime))
		body.WriteString("    ")
		body.WriteString(task.Description)
		fmt.Println(body.String())
	}

	return
}

var commandsMap = map[string]func(*CommandState) error{
	"help":   helpCommand,
	"add":    addCommand,
	"update": updateCommand,
	"delete": deleteCommand,
	"mark":   markCommand,
	"list":   listCommand,
}

func main() {
	if len(os.Args) <= 1 {
		commandsMap["help"](nil)
	}

	command := os.Args[1]
	if commandFn, ok := commandsMap[command]; !ok {
		log.Fatal("invalid command: ", command)
	} else if state, err := NewCommandState(os.Args[2:]); err != nil {
		log.Fatal(err)
	} else if err = state.TaskStore.Load(); err != nil {
		log.Fatal(err)
	} else if err = commandFn(state); err != nil {
		log.Fatal(err)
	}
}
