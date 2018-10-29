package transport

import (
	"bufio"
	"github.com/blainsmith/live-kvdb/06/db"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
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
				case "SET":
					db.Set(payload[1], []byte(payload[2]))
				case "GET":
					value := db.Get(payload[1])
					conn.Write(value)
					conn.Write([]byte("\n"))
				case "DEL":
					db.Del(payload[1])
				case "KEYS":
					keys := db.Keys()
					for _, key := range keys {
						conn.Write([]byte(key))
						conn.Write([]byte("\n"))
					}
				case "QUIT":
					conn.Write([]byte("Goodbye!\n"))
					conn.Close()
				}
			}
		}(conn)
	}
}

func NewHTTPServer(db db.KVDB, port string) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		segments := strings.Split(request.RequestURI, "/")
		key := segments[1]

		switch request.Method {
		case http.MethodGet:
			if key == "" {
				keys := db.Keys()
				writer.WriteHeader(http.StatusOK)
				for _, key := range keys {
					writer.Write([]byte(key))
					writer.Write([]byte("\n"))
				}
				return
			}

			value := db.Get(key)
			writer.WriteHeader(http.StatusOK)
			writer.Write(value)
		case http.MethodPost:
			value, _ := ioutil.ReadAll(request.Body)
			db.Set(key, value)
			writer.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			db.Del(key)
			writer.WriteHeader(http.StatusNoContent)
		}
	})

	http.ListenAndServe(port, nil)
}