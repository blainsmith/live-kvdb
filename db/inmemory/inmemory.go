package inmemory

import (
	"sync"
)

type DB struct {
	mutex sync.Mutex
	item  map[string][]byte
}

func New() *DB {
	return &DB{
		item: make(map[string][]byte),
	}
}

func (db *DB) Set(key string, value []byte) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.item[key] = value
}

func (db *DB) Get(key string) []byte {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.item[key]
}

func (db *DB) Del(key string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.item, key)
}

func (db *DB) Keys() []string {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var keys []string
	for key := range db.item {
		keys = append(keys, key)
	}

	return keys
}
