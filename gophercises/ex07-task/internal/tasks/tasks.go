package tasks

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks/tasksdb"
)

type Task tasksdb.TaskDTO

var ts []*Task

func init() {
	tsDTOs, err := tasksdb.Read()
	if err != nil {
		fmt.Printf("failed to load tasks db. Received err: %s\n", err)
	}

	ts = toTasks(tsDTOs)
}

func New(desc string) *Task {
	return &Task{Desc: desc}
}

func Add(task *Task) error {
	ts = append(ts, task)

	err := tasksdb.Write(toTaskDTOs(ts))
	if err != nil {
		return fmt.Errorf("failed to write to db: %s", err)
	}

	return nil
}

func List() []*Task {
	return ts
}

func Do(id int) error {
	if id >= len(ts) {
		return fmt.Errorf("invalid task id")
	}

	ts = append(ts[:id], ts[id+1:]...)
	err := tasksdb.Write(toTaskDTOs(ts))
	if err != nil {
		return fmt.Errorf("failed to write to db: %s", err)
	}

	return nil
}

func toTaskDTOs(ts []*Task) []*tasksdb.TaskDTO {
	tsDTOs := []*tasksdb.TaskDTO{}
	for _, t := range ts {
		tDTO := tasksdb.TaskDTO(*t)
		tsDTOs = append(tsDTOs, &tDTO)
	}

	return tsDTOs
}

func toTasks(tsDTOs []*tasksdb.TaskDTO) []*Task {
	ts := []*Task{}
	for _, tDTO := range tsDTOs {
		t := Task(*tDTO)
		ts = append(ts, &t)
	}

	return ts
}
