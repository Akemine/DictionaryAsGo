package dictionary

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"time"

	"github.com/dgraph-io/badger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Add adds a new entry to the Dictionary with the specified word, definition, and tag.
func (d *Dictionary) Add(word, definition string, tag Tag) error {
	// Create a title case converter for English language
	titleCaser := cases.Title(language.English)

	// Create an Entry struct with title-cased word, provided definition, tag, and current timestamp
	entry := Entry{
		Word:       titleCaser.String(word),
		Definition: definition,
		Tag:        tag,
		CreatedAt:  time.Now(),
	}

	// Encode the entry into a buffer using gob encoding
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(entry)
	if err != nil {
		return err
	}

	// Update the Badger database with the serialized entry
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(word), buffer.Bytes())
	})
}

// Get retrieves the Entry associated with the specified word from the Dictionary.
func (d *Dictionary) Get(word string) (Entry, error) {
	var entry Entry
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(word))
		if err != nil {
			return err
		}
		// Decode the Badger item value into an Entry object
		entry, err = getEntry(item)
		return err
	})

	return entry, err
}

// List retrieves words from the Dictionary.
// It returns:
// - []string: An alphabetically sorted array of words.
// - map[string]Entry: A map of words and their definitions.
// - error: Any error encountered during retrieval.
func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	entries := make(map[string]Entry)
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		iterator := txn.NewIterator(opts)
		defer iterator.Close()
		// On mets à 0, si c'est valid, on continue
		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			item := iterator.Item()
			entry, err := getEntry(item)
			if err != nil {
				return err
			}
			entries[entry.Word] = entry
		}
		return nil
	})
	return sortedKeys(entries), entries, err
}

// ListTag retrieves words and their corresponding entries for a given tag from the Dictionary.
// It returns:
// - []string: An alphabetically sorted array of words.
// - map[string]Entry: A map of words and their definitions.
// - error: Any error encountered during retrieval.
func (d *Dictionary) ListTag(t Tag) ([]string, map[string]Entry, error) {
	entries := make(map[string]Entry) // Initialize a map to store entries
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		iterator := txn.NewIterator(opts) // Create a new iterator for the transaction
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			item := iterator.Item()                      // Get the current item from the iterator
			entry, err := getEntry(item)                 // Decode the item value into an Entry struct
			if IsValidTag(entry.Tag) && entry.Tag == t { // Check if the entry's tag matches the provided tag
				entries[entry.Word] = entry // Add the entry to the map with word as key
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	return sortedKeys(entries), entries, err
}

// Remove deletes the specified word from the Badger database.
func (d *Dictionary) Remove(word string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(word))
	})
}

// sortedKeys returns the keys of the entries map in sorted order.
func sortedKeys(entries map[string]Entry) []string {
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// getEntry decodes a Badger item value into an Entry object.
func getEntry(item *badger.Item) (Entry, error) {
	var entry Entry
	var buffer bytes.Buffer // Create a buffer to hold item value
	err := item.Value(func(val []byte) error {
		_, err := buffer.Write(val) // Write item value to buffer
		return err
	})
	if err != nil {
		return entry, fmt.Errorf("error writing item value to buffer: %w", err)
	}
	decode := gob.NewDecoder(&buffer) // Create a new decoder for the buffer
	err = decode.Decode(&entry)       // Decode buffer contents into Entry struct
	return entry, err                 // Return decoded Entry struct
}
