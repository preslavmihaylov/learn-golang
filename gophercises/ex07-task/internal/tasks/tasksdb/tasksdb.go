// Package tasksdb provides primitives for reading and writing tasks into a database.
package tasksdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// TaskDTO encapsulates a simple task with a description as pulled from database
type TaskDTO struct {
	Desc string
}

const dbFilename = "tasks.db"

// Read reads the stored tasks in the db and returns a slice of task DTOs.
// It returns an error in case of an issue with the db.
func Read() ([]*TaskDTO, error) {
	b, err := ioutil.ReadFile(dbFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	ts := []*TaskDTO{}
	err = json.Unmarshal(b, &ts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %s", err)
	}

	return ts, nil
}

// Write persists the provided list of task DTOs in the db.
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
