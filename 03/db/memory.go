package db

import (
	"errors"
	"fmt"
	"sync"
)

type InMemory struct {
	sync.Mutex
	items map[string][]byte
}

func (db *InMemory) init() {
	if db.items == nil {
		db.items = make(map[string][]byte)
	}
}

func (db *InMemory) Get(key string) ([]byte, error) {
	db.Lock()
	defer db.Unlock()

	db.init()

	value, found := db.items[key]
	if !found {
		return nil, fmt.Errorf("%s key not found", key)
	}

	return value, nil
}

func (db *InMemory) Set(key string, value []byte) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}

	db.Lock()
	defer db.Unlock()

	db.init()

	db.items[key] = value

	return nil
}

func (db *InMemory) Del(key string) error {
	db.Lock()
	defer db.Unlock()

	db.init()

	delete(db.items, key)

	return nil
}

func (db *InMemory) Keys() ([]string, error) {
	db.Lock()
	defer db.Unlock()

	db.init()

	var keys []string
	for key := range db.items {
		keys = append(keys, key)
	}

	return keys, nil
}
