package dbconn

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var (
	dbFilename    = "tasks.db"
	dbPermissions = os.FileMode(0644)
)

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

// Pull gets the value against the given key from the provided bucket.
// It unmarshals the json stored inside into the data provided.
// data should be a pointer type.
// In case of an issue with opening db connection, reading bucket, retrieving value or unmarshaling,
// an error is returned.
// In case of failing to close db connection, a panic occurs.
func Pull(bucket []byte, key []byte, data interface{}) error {
	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return fmt.Errorf("dbconn.Pull failed opening db connection: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("dbconn.Pull failed closing database: %s", err)
		}
	}()

	return db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("dbconn.Pull failed to read bucket %s", bucket)
		}

		bs := bk.Get(key)
		if bs == nil {
			return fmt.Errorf("dbconn.Pull failed to get value for key %s", key)
		}

		err = json.Unmarshal(bs, data)
		if err != nil {
			return fmt.Errorf("dbconn.Pull failed to unmarshal data %v: %s", data, err)
		}

		return nil
	})
}

// Push marshals the provided data into json and pushes that as value
// to the provided key in the given bucket.
// In case of an issue with opening db, marshaling data, getting bucket or storing data,
// an error is returned.
// In case of failing to close db connection, a panic occurs.
func Push(bucket []byte, key []byte, data interface{}) error {
	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return fmt.Errorf("dbconn.Push failed to open db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("dbconn.Push failed to close db: %s", err)
		}
	}()

	return db.Update(func(tx *bolt.Tx) error {
		bs, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("dbconn.Push failed to marshal data: %s", err)
		}

		bk, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("dbconn.Push failed to get or create bucket %s: %s", bucket, err)
		}

		err = bk.Put(key, bs)
		if err != nil {
			return fmt.Errorf("dbconn.Push failed persisting data: %s", err)
		}

		return nil
	})
}

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
