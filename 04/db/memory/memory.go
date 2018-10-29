package memory

import (
	"errors"
	"fmt"
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

	db.mutex.Lock()
	defer db.mutex.Unlock()

	value, found := db.items[key]
	if !found {
		return nil, fmt.Errorf("key %s does not exist", key)
	}

	return value, nil
}

func (db *DB) Set(key string, value []byte) error {
	db.init()

	if key == "" {
		return errors.New("key cannot be empty")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.items[key] = value

	return nil
}

func (db *DB) Del(key string) error {
	db.init()

	if key == "" {
		return errors.New("key cannot be empty")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	_, found := db.items[key]
	if !found {
		return fmt.Errorf("key %s does not exist", key)
	}

	delete(db.items, key)

	return nil
}

func (db *DB) Keys() ([]string, error) {
	db.init()

	db.mutex.Lock()
	defer db.mutex.Unlock()

	var keys []string
	for key := range db.items {
		keys = append(keys, key)
	}

	return keys, nil
}
