package store

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

// Store is a wrapper around Bolt's database
type Store struct {
	db *bolt.DB
}

// Field contains data to interface with the database
type Field struct {
	Tag        string `json:"tag,omitempty"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// Init will return an initialised database store
func Init() (*Store, error) {
	db, err := bolt.Open("./store.bolt", 0666, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("store"))
		return nil
	})

	return &Store{
		db: db,
	}, nil
}

// GetAll will iterate over all objects in the store and return them
func (s *Store) GetAll() ([]*Field, error) {
	var all []*Field
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))

		// Iterate over items in sorted key order.
		if err := b.ForEach(func(k, v []byte) error {
			fmt.Printf("A %s is %s.\n", k, v)
			all = append(all, &Field{
				Identifier: string(k),
				Password:   string(v),
			})
			return nil
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get will retrieve a value from the store
func (s *Store) Get(key string) (*Field, error) {
	var value string

	bKey := []byte(key)
	if bKey == nil {
		return nil, ErrEmptyKey
	}

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("store"))
		// HACK(mnzt): expensive on memory
		v := b.Get(bKey)
		if v == nil {
			return ErrNotFound
		}
		value = string(v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Field{
		Identifier: key,
		Password:   value,
	}, nil
}

// Put will place a field into the database
func (s *Store) Put(in *Field) error {
	if in == nil {
		return ErrEmptyKey
	}
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		if err := b.Put([]byte(in.Identifier), []byte(in.Password)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

var (
	// ErrEmptyKey is the error for when a key is empty
	ErrEmptyKey = errors.New("error empty key provided")
	// ErrNotFound is the error when a value is not found
	ErrNotFound = errors.New("error key/value not found")
)
