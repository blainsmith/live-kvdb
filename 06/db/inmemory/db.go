package inmemory

import "sync"

type KVDB struct {
	mutex sync.Mutex
	items map[string][]byte
}

func New() *KVDB {
	return &KVDB{
		items: make(map[string][]byte),
	}
}

func (db *KVDB) Set(key string, value []byte) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.items[key] = value
}

func (db *KVDB) Get(key string) []byte {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.items[key]
}

func (db *KVDB) Del(key string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.items, key)
}

func (db *KVDB) Keys() []string {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var keys []string

	for key := range db.items {
		keys = append(keys, key)
	}

	return keys
}
