package dictionary

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

// Dictionary represents a dictionary backed by BadgerDB.
type Dictionary struct {
	db *badger.DB
}

// Entry represents an entry in the dictionary.
type Entry struct {
	Word       string    // The word being defined
	Definition string    // The definition of the word
	Tag        Tag       // The tag associated with the word
	CreatedAt  time.Time // The timestamp when the entry was created
}

// String returns a formatted string representation of the Entry.
func (e Entry) String() string {
	created := e.CreatedAt.Format(time.Stamp) // Format the timestamp
	return fmt.Sprintf("%-10v\t%-50v%-6v\t%v", e.Word, e.Definition, e.Tag, created)
}

// New creates and initializes a new Dictionary instance with a BadgerDB database.
func New(dir string) (*Dictionary, error) {
	opts := badger.DefaultOptions(dir) // Set BadgerDB options with directory path
	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err // Return error if database opening fails
	}

	dict := &Dictionary{
		db: db, // Initialize Dictionary with the opened database
	}

	return dict, nil // Return the initialized Dictionary instance
}

// Close closes the BadgerDB database associated with the Dictionary.
func (d *Dictionary) Close() {
	d.db.Close()
}
