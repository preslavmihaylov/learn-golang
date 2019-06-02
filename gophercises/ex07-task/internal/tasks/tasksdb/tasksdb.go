// Package tasksdb provides primitives for reading and writing tasks into a database.
package tasksdb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// TaskDTO encapsulates a simple task with a description as pulled from database
type TaskDTO struct {
	Desc string
}

var (
	dbFilename     = "tasks.db"
	tasksBucket    = []byte("tasks")
	activeTasksKey = []byte("active-tasks")
	emptyTasksList = []byte("[]")
	dbPermissions  = os.FileMode(0644)
)

var (
	// ErrDBNotFound occurs when the database is not found while reading it
	ErrDBNotFound = fmt.Errorf("database not found")
)

// Read the stored tasks from the db.
// Returns a slice of task DTOs.
// Returns an error in case of an issue with the db.
func Read() ([]TaskDTO, error) {
	if _, err := os.Stat(dbFilename); err != nil {
		switch {
		case os.IsNotExist(err):
			return nil, ErrDBNotFound
		default:
			return nil, fmt.Errorf("failed to open db: %s", err)
		}
	}

	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close db: %s", err)
		}
	}()

	ts := []TaskDTO{}
	err = db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(tasksBucket)
		if bk == nil {
			return fmt.Errorf("failed to read bucket %s", tasksBucket)
		}

		bs := bk.Get(activeTasksKey)
		if bs == nil {
			return fmt.Errorf("failed to get value for key %s", activeTasksKey)
		}

		err = json.Unmarshal(bs, &ts)
		if err != nil {
			return fmt.Errorf("failed to unmarshal active tasks: %s", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ts, nil
}

// Write the provided list of task DTOs in the db.
// It returns an error in case of an issue with the provided data or the db.
func Write(ts []TaskDTO) error {

	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return fmt.Errorf("failed to open db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close db: %s", err)
		}
	}()

	err = db.Update(func(tx *bolt.Tx) error {
		bs, err := json.Marshal(ts)
		if err != nil {
			return fmt.Errorf("failed to marshal tasks: %s", err)
		}

		bk := tx.Bucket(tasksBucket)
		if bk == nil {
			return fmt.Errorf("failed to read bucket %s", tasksBucket)
		}

		err = bk.Put(activeTasksKey, bs)
		if err != nil {
			return fmt.Errorf("failed persisting tasks: %s", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Create a new empty database on the filesystem
func Create() error {
	_, err := os.Stat(dbFilename)
	switch {
	case os.IsNotExist(err):
		// do nothing
	case err == nil:
		err = os.Remove(dbFilename)
		if err != nil {
			return fmt.Errorf("failed to remove old db: %s", err)
		}
	default:
		return fmt.Errorf("failed to create db: %s", err)
	}

	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return fmt.Errorf("failed creating new database: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed closing newly created database: %s", err)
		}
	}()

	err = db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucket(tasksBucket)
		if err != nil {
			return fmt.Errorf("received error while creating bucket %s: %s", tasksBucket, err)
		}

		err = bk.Put(activeTasksKey, emptyTasksList)
		if err != nil {
			return fmt.Errorf("failed setting default value for key %s: %s", activeTasksKey, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
