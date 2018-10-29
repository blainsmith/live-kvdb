package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	db := NewKVDB()

	go tcpServer(db)
	go httpServer(db)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func tcpServer(db *kvdb) {
	tcp, _ := net.Listen("tcp", ":1337")

	for {
		conn, _ := tcp.Accept()

		go func(conn net.Conn) {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				// SET key value
				payload := strings.Split(scanner.Text(), " ")

				resp := []byte("OK")
				switch strings.ToUpper(payload[0]) {
				case "GET":
					resp = db.Get(payload[1])
				case "SET":
					db.Set(payload[1], []byte(payload[2]))
				case "DEL":
					db.Del(payload[1])
				case "KEYS":
					keys := db.Keys()

					fmt.Fprint(conn, keys)
					conn.Write([]byte("\n"))
				case "QUIT":
					conn.Write([]byte("kthxbye!\n"))
					conn.Close()
					return
				default:
					conn.Write([]byte("unrecoginzed command\n"))
				}

				resp = append(resp, '\n')
				conn.Write(resp)
			}
		}(conn)
	}
}

func httpServer(db *kvdb) {
	// POST /key
	// value
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")

		if len(path) == 1 {
			keys := db.Keys()

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, keys)
			return
		}

		key := path[1]
		switch r.Method {
		case http.MethodGet:
			value := db.Get(key)

			w.WriteHeader(http.StatusOK)
			w.Write(value)
		case http.MethodPost:
			value, _ := ioutil.ReadAll(r.Body)
			db.Set(key, value)

			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			db.Del(key)

			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No"))
		}
	})

	http.ListenAndServe(":1338", nil)
}

type kvdb struct {
	mu    sync.Mutex
	items map[string][]byte
}

func NewKVDB() *kvdb {
	return &kvdb{
		items: make(map[string][]byte),
	}
}

func (db *kvdb) Get(key string) []byte {
	db.mu.Lock()
	defer db.mu.Unlock()

	hash := sha256.Sum256([]byte(key))

	return db.items[fmt.Sprintf("%x", hash)]
}

func (db *kvdb) Set(key string, value []byte) {
	db.mu.Lock()
	defer db.mu.Unlock()

	hash := sha256.Sum256([]byte(key))

	db.items[fmt.Sprintf("%x", hash)] = value
}

func (db *kvdb) Del(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	hash := sha256.Sum256([]byte(key))

	delete(db.items, fmt.Sprintf("%x", hash))
}

func (db *kvdb) Keys() []string {
	keys := make([]string, len(db.items))

	for key := range db.items {
		keys = append(keys, key)
	}

	return keys
}
