package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	db := kvdb{}

	go startTCPServer(&db)
	go startHTTPServer(&db)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

type kvdb struct {
	mu    sync.Mutex
	items map[string][]byte
}

func (db *kvdb) get(k string) ([]byte, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if k == "" {
		return nil, errors.New("key is empty")
	}

	return db.items[k], nil
}

func (db *kvdb) set(k string, v []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.items == nil {
		db.items = make(map[string][]byte)
	}

	if k == "" {
		return errors.New("key is empty")
	}

	db.items[k] = v

	return nil
}

func (db *kvdb) del(k string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if k == "" {
		return errors.New("key is empty")
	}

	delete(db.items, k)

	return nil
}

func startTCPServer(db *kvdb) {
	listener, _ := net.Listen("tcp", ":1313")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()

		go func(conn net.Conn) {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				payload := strings.Split(scanner.Text(), " ")

				var resp []byte

				switch strings.ToUpper(payload[0]) {
				case "GET":
					resp, _ = db.get(payload[1])
				case "SET":
					db.set(payload[1], []byte(payload[2]))
				case "DEL":
					db.del(payload[1])
				case "QUIT":
					conn.Close()
					return
				}

				if len(resp) == 0 {
					resp = append(resp, []byte("OK")...)
				}

				resp = append(resp, '\n')
				conn.Write(resp)
			}
		}(conn)
	}
}

func startHTTPServer(db *kvdb) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		key := path[1]

		var resp []byte
		switch r.Method {
		case http.MethodGet:
			resp, _ = db.get(key)

			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		case http.MethodPost:
			value, _ := ioutil.ReadAll(r.Body)
			db.set(key, value)

			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			db.del(key)

			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.ListenAndServe(":8080", nil)
}
