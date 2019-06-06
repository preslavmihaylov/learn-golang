package dbconn

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

// DBConnection encapsulates primitives for working with underlying database
type DBConnection struct {
	filename   string
	perms      os.FileMode
	connection *bolt.DB
}

// New creates a new DBConnection with the provided options for the underlying database
func New(filename string, perms os.FileMode) *DBConnection {
	return &DBConnection{filename, perms, nil}
}

// Open a new connection to the underlying database.
// In case of an issue with opening connection, an error is returned.
func (dbc *DBConnection) Open() error {
	connection, err := bolt.Open(dbc.filename, dbc.perms, nil)
	if err != nil {
		return fmt.Errorf("failed to open connection to db: %s", err)
	}

	dbc.connection = connection
	return nil
}

// Close the existing connection to the underlying database.
// In case of no open connection or an issue with closing it, an error is returned.
func (dbc *DBConnection) Close() error {
	if dbc.connection == nil {
		return fmt.Errorf("connection was not established for it to be closed")
	}

	err := dbc.connection.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection to db: %s", err)
	}

	dbc.connection = nil
	return nil
}
