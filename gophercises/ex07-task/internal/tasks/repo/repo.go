// Package repo provides primitives for reading and writing tasks into a database.
package repo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex07-task/internal/dbconn"
)

// TaskDTO encapsulates a simple task with a description and a completed status
type TaskDTO struct {
	ID       []byte
	Desc     string
	Complete bool
}

var (
	tasksBucket   = []byte("tasks")
	dbFilename    = "tasks.db"
	dbPermissions = os.FileMode(0644)
)

var dbConnection *dbconn.DBConnection

func init() {
	dbConnection = dbconn.New(dbFilename, dbPermissions)
	if !dbConnection.DBExists() {
		err := dbConnection.CreateDB()
		if err != nil {
			log.Fatalf("failed to create empty db: %s", err)
		}

		err = dbConnection.Open()
		if err != nil {
			log.Fatalf("failed to open connection to db: %s", err)
		}
		defer func() {
			err := dbConnection.Close()
			if err != nil {
				log.Fatalf("failed to close connection to db: %s", err)
			}
		}()

		err = dbConnection.CreateBucket(tasksBucket)
		if err != nil {
			log.Fatalf("failed to create tasks bucket in db: %s", err)
		}
	}
}

// GetAllIncomplete gets all incomplete tasks from the db.
// Returns an error in case of an issue with the db.
func GetAllIncomplete() ([]TaskDTO, error) {
	tsks, err := getAllWithStatus(false)
	if err != nil {
		return nil, err
	}

	return tsks, nil
}

// GetAllComplete gets all complete tasks from the db.
// Returns an error in case of an issue with the db.
func GetAllComplete() ([]TaskDTO, error) {
	tsks, err := getAllWithStatus(true)
	if err != nil {
		return nil, err
	}

	return tsks, nil
}

// Add the provided TaskDTO in the db.
// It returns an error in case of an issue with the provided data or the db.
func Add(tsk TaskDTO) error {
	if !dbConnection.DBExists() {
		return fmt.Errorf("database not found")
	}

	err := dbConnection.Open()
	if err != nil {
		return fmt.Errorf("failed to open connection to db: %s", err)
	}
	defer func() {
		err := dbConnection.Close()
		if err != nil {
			log.Fatalf("failed to close connection to db: %s", err)
		}
	}()

	nextIDbs, err := dbConnection.NextIDForBucket(tasksBucket)
	if err != nil {
		return fmt.Errorf("failed to get next id for bucket %s: %s", tasksBucket, err)
	}

	tsk.ID = nextIDbs
	bs, err := json.Marshal(tsk)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %s", err)
	}

	err = dbConnection.Add(tasksBucket, bs)
	if err != nil {
		return fmt.Errorf("write failed to push to db: %s", err)
	}

	return nil
}

// Put the provided TaskDTO in the db.
// It returns an error in case of an issue with provided data or db.
func Put(tsk TaskDTO) error {
	if !dbConnection.DBExists() {
		return fmt.Errorf("database not found")
	}

	bs, err := json.Marshal(tsk)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %s", err)
	}

	err = dbConnection.Open()
	if err != nil {
		return fmt.Errorf("failed to open connection to db: %s", err)
	}
	defer func() {
		err := dbConnection.Close()
		if err != nil {
			log.Fatalf("failed to close connection to db: %s", err)
		}
	}()

	err = dbConnection.Put(tasksBucket, tsk.ID, bs)
	if err != nil {
		return fmt.Errorf("failed to put value in db: %s", err)
	}

	return nil
}

func getAllWithStatus(completionStatus bool) ([]TaskDTO, error) {
	if !dbConnection.DBExists() {
		return nil, fmt.Errorf("database not found")
	}

	tsks := []TaskDTO{}

	err := dbConnection.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to db: %s", err)
	}
	defer func() {
		err := dbConnection.Close()
		if err != nil {
			log.Fatalf("failed to close connection to db: %s", err)
		}
	}()

	err = dbConnection.ForEach(tasksBucket, func(val []byte) error {
		tsk := TaskDTO{}
		err := json.Unmarshal(val, &tsk)
		if err != nil {
			return fmt.Errorf("failed to unmarshal data: %s", err)
		}

		if tsk.Complete == completionStatus {
			tsks = append(tsks, tsk)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repo.Read failed to get data from db: %s", err)
	}

	return tsks, nil
}
