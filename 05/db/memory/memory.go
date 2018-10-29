package memory

import (
	"errors"
	"sync"
)

type DB struct {
	mutex sync.Mutex
	items map[string][]byte
}

func (db *DB) init() {
	if db.items == nil {
		db.items = make(map[string][]byte)
	}
}

func (db *DB) Get(key string) ([]byte, error) {
	db.init()

	value, found := db.items[key]
	if !found {
		return nil, errors.New("not found")
	}

	return value, nil
}

func (db *DB) Set(key string, value []byte) error {
	db.init()

	if key == "" {
		return errors.New("key must not be empty")
	}

	db.items[key] = value

	return nil
}

func (db *DB) Del(key string) error {
	delete(db.items, key)

	return nil
}

func (db *DB) Keys() ([]string, error) {
	var keys []string

	for key := range db.items {
		keys = append(keys, key)
	}

	return keys, nil
}
