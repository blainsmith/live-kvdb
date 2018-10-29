package transport

import (
	"bufio"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/blainsmith/live-kvdb/db"
)

func NewTCPServer(db db.KVDB, port string) {
	tcp, _ := net.Listen("tcp", port)

	for {
		conn, _ := tcp.Accept()

		go func(conn net.Conn) {
			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				payload := strings.Split(scanner.Text(), " ")

				switch strings.ToUpper(payload[0]) {
				case "GET":
					value := db.Get(payload[1])
					conn.Write(value)
					conn.Write([]byte("\n"))
				case "SET":
					db.Set(payload[1], []byte(payload[2]))
				case "DEL":
					db.Del(payload[1])
				case "KEYS":
					keys := db.Keys()

					for _, key := range keys {
						conn.Write([]byte(key))
						conn.Write([]byte("\n"))
					}
				case "QUIT":
					conn.Write([]byte("Goodbye\n"))
					conn.Close()
				}
			}
		}(conn)
	}
}

func NewHTTPServer(db db.KVDB, port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.RequestURI, "/")
		key := segments[1]

		switch r.Method {
		case http.MethodGet:
			if key == "" {
				keys := db.Keys()

				for _, key := range keys {
					w.Write([]byte(key))
					w.Write([]byte("\n"))
				}
				w.WriteHeader(http.StatusOK)
				return
			}

			value := db.Get(key)

			w.WriteHeader(http.StatusOK)
			w.Write(value)
		case http.MethodPost:
			value, _ := ioutil.ReadAll(r.Body)

			db.Set(key, value)

			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			db.Del(key)
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.ListenAndServe(port, nil)
}
