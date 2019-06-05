// Package repo provides primitives for reading and writing tasks into a database.
package repo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/dbconn"
)

// TaskDTO encapsulates a simple task with a description as pulled from database
type TaskDTO struct {
	ID       []byte
	Desc     string
	Complete bool
}

var (
	tasksBucket = []byte("tasks")
)

var (
	// ErrDBNotFound occurs when the database is not found while reading it
	ErrDBNotFound = fmt.Errorf("database not found")
)

func init() {
	if !dbconn.DBExists() {
		err := dbconn.CreateDB()
		if err != nil {
			log.Fatalf("failed to create empty db: %s", err)
		}

		err = dbconn.CreateBucket(tasksBucket)
		if err != nil {
			log.Fatalf("failed to create tasks bucket in db: %s", err)
		}
	}
}

// GetAllIncomplete the stored tasks from the db.
// Returns a slice of task DTOs.
// Returns an error in case of an issue with the db.
func GetAllIncomplete() ([]TaskDTO, error) {
	if !dbconn.DBExists() {
		return nil, ErrDBNotFound
	}

	tsks := []TaskDTO{}
	err := dbconn.ForEach(tasksBucket, func(val []byte) error {
		tsk := TaskDTO{}
		err := json.Unmarshal(val, &tsk)
		if err != nil {
			return fmt.Errorf("failed to unmarshal data: %s", err)
		}

		if !tsk.Complete {
			tsks = append(tsks, tsk)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repo.Read failed to get data from db: %s", err)
	}

	return tsks, nil
}

// Add the provided list of task DTOs in the db.
// It returns an error in case of an issue with the provided data or the db.
func Add(tsk TaskDTO) error {
	if !dbconn.DBExists() {
		return ErrDBNotFound
	}

	nextIDbs, err := dbconn.NextIDForBucket(tasksBucket)
	if err != nil {
		return fmt.Errorf("failed to get next id for bucket %s: %s", tasksBucket, err)
	}

	tsk.ID = nextIDbs
	err = dbconn.Add(tasksBucket, tsk)
	if err != nil {
		return fmt.Errorf("write failed to push to db: %s", err)
	}

	return nil
}

func Put(tsk TaskDTO) error {
	if !dbconn.DBExists() {
		return ErrDBNotFound
	}

	err := dbconn.Put(tasksBucket, tsk.ID, tsk)
	if err != nil {
		return fmt.Errorf("failed to put value in db: %s", err)
	}

	return nil
}
