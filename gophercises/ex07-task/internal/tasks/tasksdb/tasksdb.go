// Package tasksdb provides primitives for reading and writing tasks into a database.
package tasksdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// TaskDTO encapsulates a simple task with a description as pulled from database
type TaskDTO struct {
	Desc string
}

const dbFilename = "tasks.db"

var (
	// ErrDBNotFound occurs when the database is not found while reading it
	ErrDBNotFound = fmt.Errorf("database not found")
)

// Read the stored tasks from the db.
// Returns a slice of task DTOs.
// Returns an error in case of an issue with the db.
func Read() ([]*TaskDTO, error) {
	b, err := ioutil.ReadFile(dbFilename)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return nil, ErrDBNotFound
		default:
			return nil, fmt.Errorf("failed to open db: %s", err)
		}
	}

	if len(b) == 0 {
		return []*TaskDTO{}, nil
	}

	ts := []*TaskDTO{}
	err = json.Unmarshal(b, &ts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %s", err)
	}

	return ts, nil
}

// Write the provided list of task DTOs in the db.
// It returns an error in case of an issue with the provided data or the db.
func Write(ts []*TaskDTO) error {
	b, err := json.Marshal(ts)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %s", err)
	}

	err = ioutil.WriteFile(dbFilename, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}

	return nil
}

// Create a new empty database on the filesystem
func Create() error {
	f, err := os.Create(dbFilename)
	if err != nil {
		return fmt.Errorf("failed creating new database: %s", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed closing newly created database: %s", err)
	}

	return nil
}
