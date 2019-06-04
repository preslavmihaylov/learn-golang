// Package tasks deals with managing adding, doing and listing tasks
package tasks

import (
	"fmt"
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/tasks/repo"
)

// Task encapsulates a simple task with a description
type Task repo.TaskDTO

var ts []Task

func init() {
	tsDTOs, err := repo.Read()
	if err != nil {
		switch err {
		case repo.ErrDBNotFound:
			fmt.Println("Database not found. Creating a new empty DB...")

			err = repo.Create()
			if err != nil {
				log.Fatalf("failed creating new database: %s", err)
			}

			ts = []Task{}
		default:
			log.Fatalf("error while reading database: %s", err)
		}

		return
	}

	ts = toTasks(tsDTOs)
}

// New returns a new Task with the provided description
func New(desc string) *Task {
	return &Task{Desc: desc}
}

// Add adds the provided task to the tasks list.
// It returns an error in case of a problem with the db.
func Add(task *Task) error {
	ts = append(ts, *task)

	err := repo.Write(toTaskDTOs(ts))
	if err != nil {
		return fmt.Errorf("failed to write to db: %s", err)
	}

	return nil
}

// List returns the currently active tasks
func List() []Task {
	return ts
}

// Do marks a given task complete, given its id.
// It returns an error in case of an invalid task id or a problem with the db.
func Do(id int) error {
	if id >= len(ts) {
		return fmt.Errorf("invalid task id")
	}

	ts = append(ts[:id], ts[id+1:]...)
	err := repo.Write(toTaskDTOs(ts))
	if err != nil {
		return fmt.Errorf("failed to write to db: %s", err)
	}

	return nil
}

func toTaskDTOs(ts []Task) []repo.TaskDTO {
	tsDTOs := []repo.TaskDTO{}
	for _, t := range ts {
		tDTO := repo.TaskDTO(t)
		tsDTOs = append(tsDTOs, tDTO)
	}

	return tsDTOs
}

func toTasks(tsDTOs []repo.TaskDTO) []Task {
	ts := []Task{}
	for _, tDTO := range tsDTOs {
		t := Task(tDTO)
		ts = append(ts, t)
	}

	return ts
}
