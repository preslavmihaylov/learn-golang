package dbconn

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// NextIDForBucket gets the next autoincremented ID for given bucket, without actually incrementing it.
// In case of no open connection or an issue with reading from bucket, an error is returned.
func (dbc *DBConnection) NextIDForBucket(bucket []byte) ([]byte, error) {
	if dbc.connection == nil {
		return nil, fmt.Errorf("connection not open")
	}

	nextIDbs := []byte{}
	err := dbc.connection.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("failed to get bucket %s", bucket)
		}

		nextIDbs = itob(bk.Sequence() + 1)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nextIDbs, nil
}

// CreateBucket will add a new bucket with the provided identifier in the db.
// In case of no open connection or an issue with creating bucket, an error is returned.
func (dbc *DBConnection) CreateBucket(bucket []byte) error {
	if dbc.connection == nil {
		return fmt.Errorf("connection not open")
	}

	return dbc.connection.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucket)
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %s", bucket, err)
		}

		return nil
	})
}
