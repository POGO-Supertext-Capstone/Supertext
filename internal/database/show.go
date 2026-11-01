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
	"encoding/json"
	"bytes"
	"fmt"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
	bolt "go.etcd.io/bbolt"
)

// uploadShow uploads the slides for a show into the supertext bucket. Uses
// the show's name as the bucket, and holds the data as "text": json slide content
// and "filename": filename of show
func UploadShow(slides []upload.SlideContent, name string, filename string) error {
	dataStore, err := Open()
	if err != nil {
		return fmt.Errorf("UploadShow() failed to open db: %s", err)
	}
	defer dataStore.Close()

	// Serialize slides
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)	// so HTML tags can be read in JSON

	err = encoder.Encode(slides)
	if err != nil {
		return fmt.Errorf("failed to serialize slides: %s", err)
	}
	data := buf.Bytes()

	err = dataStore.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("supertext"))
		if bucket == nil {
			return fmt.Errorf("supertext bucket does not exist")
		}

		// grabs show bucket. If another show with the same name exists, will
		// use that bucket
		showBucket, err := bucket.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("failed to create show bucket %s: %s", name, err)
		}

		if err := showBucket.Put([]byte("text"), data); err != nil {
			return fmt.Errorf("failed to store slides: %s", err)
		}

		if err := showBucket.Put([]byte("filename"), []byte(filename)); err != nil {
			return fmt.Errorf("failed to upload filename: %s", err)
		}

		return nil
	})

	return err
}

// GetShowSlides grabs the slides for a show in the SlideContent struct
func GetShowSlides(name string) ([]upload.SlideContent, error) {
	dataStore, err := Open()
	if err != nil {
		return nil, fmt.Errorf("GetShowSlides() failed to open db: %s", err)
	}
	defer dataStore.Close()

	var slides []upload.SlideContent

	err = dataStore.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("supertext"))
		if bucket == nil {
			return fmt.Errorf("supertext bucket does not exist")
		}

		showBucket := bucket.Bucket([]byte(name))
		if showBucket == nil {
			return fmt.Errorf("show bucket %s does not exist", name)
		}

		data := showBucket.Get([]byte("text"))
		if data == nil {
			return fmt.Errorf("no slide data found for show %s", name)
		}

		// convert back to SlideContent
		if err := json.Unmarshal(data, &slides); err != nil {
			return fmt.Errorf("failed to deserialize slides: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return slides, nil
}

// GetShowJSON retrieves the raw JSON slide data for a show
func GetShowJSON(name string) ([]byte, error) {
	dataStore, err := Open()
	if err != nil {
		return nil, fmt.Errorf("GetShowJSON() failed to open db: %s", err)
	}
	defer dataStore.Close()

	var data []byte

	err = dataStore.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("supertext"))
		if bucket == nil {
			return fmt.Errorf("supertext bucket does not exist")
		}

		showBucket := bucket.Bucket([]byte(name))
		if showBucket == nil {
			return fmt.Errorf("show bucket %s does not exist", name)
		}

		value := showBucket.Get([]byte("text"))
		if value == nil {
			return fmt.Errorf("no slide data found for show %s", name)
		}

		// Copy data
		data = append([]byte(nil), value...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetShowFilename grabs the show's file name from the database
func GetShowFilename(name string) (string, error) {
	dataStore, err := Open()
	if err != nil {
		return "", fmt.Errorf("GetShowSlides() failed to open db: %s", err)
	}
	defer dataStore.Close()

	var filename string

	err = dataStore.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("supertext"))
		if bucket == nil {
			return fmt.Errorf("supertext bucket does not exist")
		}

		showBucket := bucket.Bucket([]byte(name))
		if showBucket == nil {
			return fmt.Errorf("show bucket %s does not exist", name)
		}

		// grab filename
		data := showBucket.Get([]byte("filename"))
		if data == nil {
			return fmt.Errorf("no filename found for show %s", name)
		}
		filename = string(data)

		return nil
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}

func DeleteShow(name string) error {
	dataStore, err := Open()
	if err != nil {
		return fmt.Errorf("DeleteShow() failed to open db: %s", err)
	}
	defer dataStore.Close()

	err = dataStore.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("supertext"))
		if bucket == nil {
			return fmt.Errorf("supertext bucket does not exist")
		}

		// Delete the show bucket
		err := bucket.DeleteBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("failed to delete show %s: %s", name, err)
		}

		return nil
	})

	return err
}
