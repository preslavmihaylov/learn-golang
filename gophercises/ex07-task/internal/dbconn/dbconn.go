package dbconn

import (
	"encoding/binary"
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

// Get the value against the given key from the provided bucket.
// It unmarshals the json stored inside into the data provided.
// data should be a pointer type.
// In case of an issue with opening db connection, reading bucket, retrieving value or unmarshaling,
// an error is returned.
// In case of failing to close db connection, a panic occurs.
func Get(bucket []byte, key []byte, data interface{}) error {
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

func Add(bucket []byte, data interface{}) error {
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

		nextID, err := bk.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to get new id from bucket: %s", err)
		}

		err = bk.Put(itob(nextID), bs)
		if err != nil {
			return fmt.Errorf("dbconn.Push failed persisting data: %s", err)
		}

		return nil
	})
}

// Put marshals the provided data into json and sets that as value
// to the provided key in the given bucket.
// In case of an issue with opening db, marshaling data, getting bucket or storing data,
// an error is returned.
// In case of failing to close db connection, a panic occurs.
func Put(bucket []byte, key []byte, data interface{}) error {
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

func ForEach(bucket []byte, callback func(val []byte) error) error {
	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return fmt.Errorf("dbconn.ForEach failed to open db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("dbconn.ForEach failed to close db: %s", err)
		}
	}()

	return db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("dbconn.ForEach failed to get bucket %s: %s", bucket, err)
		}

		cursor := bk.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			err := callback(v)
			if err != nil {
				return fmt.Errorf("dbconn.ForEach received error from callback: %s", err)
			}
		}

		return nil
	})
}

func NextIDForBucket(bucket []byte) ([]byte, error) {
	db, err := bolt.Open(dbFilename, dbPermissions, nil)
	if err != nil {
		return nil, fmt.Errorf("dbconn.NextIDForBucket failed to open db: %s", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("dbconn.NextIDForBucket failed to close db: %s", err)
		}
	}()

	nextIDbs := []byte{}
	err = db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("dbconn.NextIDForBucket failed to get bucket %s", bucket)
		}

		nextIDbs = itob(bk.Sequence() + 1)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nextIDbs, nil
}

func CreateBucket(bucket []byte) error {
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
		_, err := tx.CreateBucket(bucket)
		if err != nil {
			return fmt.Errorf("dbconn.CreateBucket failed to create bucket %s: %s", bucket, err)
		}

		return nil
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
