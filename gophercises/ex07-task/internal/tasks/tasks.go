// Package tasks deals with managing adding, doing and listing tasks
package tasks

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks/repo"
)

// Task encapsulates a simple task with a description
type Task repo.TaskDTO

// New returns a new Task with the provided description
func New(desc string) *Task {
	return &Task{ID: []byte{}, Desc: desc, Complete: false}
}

// Add adds the provided task to the tasks list.
// It returns an error in case of a problem with the db.
func Add(task *Task) error {
	err := repo.Add(repo.TaskDTO(*task))
	if err != nil {
		return fmt.Errorf("failed to write to db: %s", err)
	}

	return nil
}

// List returns the currently active tasks
func List() ([]Task, error) {
	tskDTOs, err := repo.GetAllIncomplete()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks from repo: %s", err)
	}

	return toTasks(tskDTOs), nil
}

// Do marks a given task complete, given its id.
// It returns an error in case of an invalid task id or a problem with the db.
func Do(id int) error {
	tskDTOs, err := repo.GetAllIncomplete()
	if err != nil {
		return fmt.Errorf("failed to get all from repo: %s", err)
	}

	if id >= len(tskDTOs) {
		return fmt.Errorf("invalid task id")
	}

	tskDTOs[id].Complete = true
	err = repo.Put(tskDTOs[id])
	if err != nil {
		return fmt.Errorf("failed to put task in repo: %s", err)
	}

	return nil
}

func toTasks(tsDTOs []repo.TaskDTO) []Task {
	ts := []Task{}
	for _, tDTO := range tsDTOs {
		t := Task(tDTO)
		ts = append(ts, t)
	}

	return ts
}
