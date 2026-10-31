package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type DataStore struct {
	db   *bolt.DB
	path string
}

// Update this number whenever the database schema changes
const DATABASE_SCHEMA_VERSION = "1"

// Gets the filesystem path where to store the database
// The database is either stored within XDG_CONFIG_HOME or XDG_DATA_HOME ($HOME/.config, $HOME/.local/share) respectively,
// within the Supertext folder.
func ConstructDatabasePath() string {
	dir, dirErr := os.UserConfigDir()
	if dirErr != nil {
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			log.Fatal("Missing both XDG_CONFIG_HOME and XDG_DATA_HOME, config")
		}
	}

	dirErr = os.Mkdir(filepath.Join(dir, "Supertext"), 0750)
	if dirErr != nil && !os.IsExist(dirErr) {
		log.Fatal("Unable to create storage directory")
	}

	configPath := filepath.Join(dir, "Supertext", "database.db")
	return configPath
}

// Opens the Database file in read-write mode, returning the DataStore struct.
// Ensures the default buckets are provided and the database schema version matches. If it doesn't match,
// it is possible to perform migrations as necessary.
//
// Remember to run Close() before the application exits, to ensure transactions are safely written and the database
// isn't corrupted.
func Open() (*DataStore, error) {
	path := ConstructDatabasePath()

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("supertext"))
		if err != nil {
			return fmt.Errorf("Creating default bucket failed: %s", err)
		}

		sv := b.Get([]byte("schema"))
		if sv == nil {
			err = b.Put([]byte("schema"), []byte(DATABASE_SCHEMA_VERSION))
			return err
		}

		if string(sv) != DATABASE_SCHEMA_VERSION {
			log.Print("Database schema out of date")
		}

		return nil
	})

	return &DataStore{
		db,
		path,
	}, nil
}

// Closes the database, releasing the file-lock. Use Open() to re-open the database instance
//
// Do not run a transaction on the closed database.
func (ds *DataStore) Close() {
	ds.db.Close()
}

// Gets the current database schema version.
//
// This function is thread-safe
func (ds *DataStore) GetVersion() (string, error) {
	var schema string

	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("supertext"))
		v := b.Get([]byte("schema"))
		schema = string(v)
		return nil
	})

	return schema, err
}
