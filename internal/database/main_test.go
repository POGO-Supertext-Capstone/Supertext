// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
