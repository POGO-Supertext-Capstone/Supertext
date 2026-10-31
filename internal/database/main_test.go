package database

import (
	"fmt"
	"sync"
	"testing"
)

// TestGetDatabaseSchema tests that we can fetch the database
// schema version. This acts as a simple sanity check that the database,
// getting, and putting data are working correctly
func TestGetDatabaseSchema(t *testing.T) {
	t.Run("schema exists", func(t *testing.T) {
		db, err := Open()
		if err != nil {
			t.Errorf("open(): %s", err)
		}

		var version string
		version, err = db.GetVersion()
		if err != nil {
			t.Errorf("GetVersion(): %s", err)
		}

		if version != DATABASE_SCHEMA_VERSION {
			t.Errorf("GetVersion() = %+v, want %+v", version, DATABASE_SCHEMA_VERSION)
		}

		db.Close()
	})
}

// Test that the database transaction supports concurrency. This doesn't
// necessarily stress-test the database, just that the database engine
// can support goroutines.
func TestTransactionConcurrency(t *testing.T) {
	t.Run("schema concurrency", func(t *testing.T) {
		var wg sync.WaitGroup
		count := 10

		db, dbErr := Open()
		if dbErr != nil {
			t.Errorf("open(): %s", dbErr)
		}

		for i := range count {
			wg.Go(func() {
				fmt.Printf("Running instance %d\n", i)
				version, err := db.GetVersion()
				if err != nil {
					t.Errorf("GetVersion(): %s", err)
				}

				if version != DATABASE_SCHEMA_VERSION {
					t.Errorf("GetVersion() = %+v, want %+v", version, DATABASE_SCHEMA_VERSION)
				}
			})
		}

		wg.Wait()

		db.Close()
	})
}
