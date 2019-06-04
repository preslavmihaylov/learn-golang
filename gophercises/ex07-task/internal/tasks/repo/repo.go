// Package repo provides primitives for reading and writing tasks into a database.
package repo

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/dbconn"
)

// TaskDTO encapsulates a simple task with a description as pulled from database
type TaskDTO struct {
	Desc string
}

var (
	tasksBucket    = []byte("tasks")
	activeTasksKey = []byte("active-tasks")
)

var (
	// ErrDBNotFound occurs when the database is not found while reading it
	ErrDBNotFound = fmt.Errorf("database not found")
)

// Read the stored tasks from the db.
// Returns a slice of task DTOs.
// Returns an error in case of an issue with the db.
func Read() ([]TaskDTO, error) {
	if !dbconn.DBExists() {
		return nil, ErrDBNotFound
	}

	ts := []TaskDTO{}
	err := dbconn.Pull(tasksBucket, activeTasksKey, &ts)
	if err != nil {
		return nil, fmt.Errorf("repo.Read failed to get data from db: %s", err)
	}

	return ts, nil
}

// Write the provided list of task DTOs in the db.
// It returns an error in case of an issue with the provided data or the db.
func Write(ts []TaskDTO) error {
	if !dbconn.DBExists() {
		return ErrDBNotFound
	}

	err := dbconn.Push(tasksBucket, activeTasksKey, ts)
	if err != nil {
		return fmt.Errorf("write failed to push to db: %s", err)
	}

	return nil
}

// Create a new empty tasks database on the filesystem
func Create() error {
	err := dbconn.RemoveDB()
	if err != nil {
		return fmt.Errorf("repo.Create failed to remove db: %s", err)
	}

	err = dbconn.Push(tasksBucket, activeTasksKey, []TaskDTO{})
	if err != nil {
		return fmt.Errorf("repo.Create failed to push empty tasks list to db: %s", err)
	}

	return nil
}
