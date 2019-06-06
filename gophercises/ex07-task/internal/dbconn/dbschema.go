package dbconn

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// DBExists checks if the database exists.
func (dbc *DBConnection) DBExists() bool {
	if _, err := os.Stat(dbc.filename); err != nil {
		switch {
		case os.IsNotExist(err):
			return false
		default:
			log.Fatalf("received unexpected error when reading db: %s", err)
		}
	}

	return true
}

// CreateDB creates a new empty database.
// In case of an issue with removing old database or creating a new one, an error is returned.
// In case of failing to close db connection, a panic occurs.
func (dbc *DBConnection) CreateDB() error {
	err := dbc.RemoveDB()
	if err != nil {
		return fmt.Errorf("failed to remove existing db: %s", err)
	}

	db, err := bolt.Open(dbc.filename, dbc.perms, nil)
	if err != nil {
		return fmt.Errorf("failed to create new db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed closing database: %s", err)
		}
	}()

	return nil
}

// RemoveDB removes the existing database or does nothing if db doesn't exist.
// In case of an issue with removing old database, an error is returned.
func (dbc *DBConnection) RemoveDB() error {
	if dbc.DBExists() {
		err := os.Remove(dbc.filename)
		if err != nil {
			return fmt.Errorf("failed to remove old db: %s", err)
		}
	}

	return nil
}
