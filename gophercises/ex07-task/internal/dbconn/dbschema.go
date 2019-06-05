package dbconn

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// DBExists checks if the database exists.
// It returns an error in case there is an issue with reading the database info from the filesystem.
func DBExists() bool {
	if _, err := os.Stat(dbFilename); err != nil {
		switch {
		case os.IsNotExist(err):
			return false
		default:
			log.Fatalf("dbconn.DBExists received unexpected error when reading db: %s", err)
		}
	}

	return true
}

// CreateDB creates a new empty database.
// In case of an issue with removing old database or creating a new one, an error is returned.
// In case of failing to close db connection, a panic occurs.
func CreateDB() error {
	err := RemoveDB()
	if err != nil {
		return fmt.Errorf("dbconn.CreateDB failed to remove existing db: %s", err)
	}

	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return fmt.Errorf("dbconn.CreateDB failed to create new db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("dbconn.CreateDB failed closing database: %s", err)
		}
	}()

	return nil
}

// RemoveDB removes the existing database.
// If no database exists, it returns nil.
// In case of an error with removing old database, an error is returned.
func RemoveDB() error {
	if DBExists() {
		err := os.Remove(dbFilename)
		if err != nil {
			return fmt.Errorf("dbconn.RemoveDB failed to remove old db: %s", err)
		}
	}

	return nil
}
