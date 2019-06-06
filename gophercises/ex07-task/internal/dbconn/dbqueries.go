package dbconn

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
)

// Get the value against the given key from the provided bucket.
// In case of no open connection or an issue with reading bucket or retrieving value, error is returned.
func (dbc *DBConnection) Get(bucket []byte, key []byte) ([]byte, error) {
	if dbc.connection == nil {
		return nil, fmt.Errorf("connection not open")
	}

	outBs := []byte{}
	err := dbc.connection.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("failed to read bucket %s", bucket)
		}

		bs := bk.Get(key)
		if bs == nil {
			return fmt.Errorf("failed to get value for key %s", key)
		}

		outBs = bs
		return nil
	})
	if err != nil {
		return nil, err
	}

	return outBs, nil
}

// Add creates a new record in the bucket by adding the provided data.
// In case of no open connection or an issue with reading bucket or adding value, an error is returned.
func (dbc *DBConnection) Add(bucket []byte, data []byte) error {
	if dbc.connection == nil {
		return fmt.Errorf("connection not open")
	}

	return dbc.connection.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("failed to get or create bucket %s: %s", bucket, err)
		}

		nextID, err := bk.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to get new id from bucket: %s", err)
		}

		err = bk.Put(itob(nextID), data)
		if err != nil {
			return fmt.Errorf("failed persisting data: %s", err)
		}

		return nil
	})
}

// Put modifies the existing value from the bucket[key] with the provided data.
// In case of no open connection or an issue with reading bucket or putting value, an error is returned.
func (dbc *DBConnection) Put(bucket []byte, key []byte, data []byte) error {
	if dbc.connection == nil {
		return fmt.Errorf("connection not open")
	}

	return dbc.connection.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("failed to get or create bucket %s: %s", bucket, err)
		}

		err = bk.Put(key, data)
		if err != nil {
			return fmt.Errorf("failed persisting data: %s", err)
		}

		return nil
	})
}

// Delete creates a new record in the bucket by adding the provided data.
// In case of no open connection or an issue with reading bucket or adding value, an error is returned.
func (dbc *DBConnection) Delete(bucket []byte, key []byte) error {
	if dbc.connection == nil {
		return fmt.Errorf("connection not open")
	}

	return dbc.connection.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("failed to read bucket %s", bucket)
		}

		err := bk.Delete(key)
		if err != nil {
			return fmt.Errorf("failed to delete value for key %s: %s", key, err)
		}

		return nil
	})
}

// ForEach iterates over all records in provided bucket, and invokes callback(record).
// In case of no open connection or an issue with reading bucket or values, an error is returned.
func (dbc *DBConnection) ForEach(bucket []byte, callback func(val []byte) error) error {
	if dbc.connection == nil {
		return fmt.Errorf("connection not open")
	}

	return dbc.connection.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("failed to get bucket %s", bucket)
		}

		cursor := bk.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			err := callback(v)
			if err != nil {
				return fmt.Errorf("received error from callback: %s", err)
			}
		}

		return nil
	})
}

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

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
