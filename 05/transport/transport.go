package transport

import (
	"bufio"
	"log"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/blainsmith/live-kvdb/05/db"
)

func NewTCPServer(port string, db db.KVDB) {
	tcp, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Println(err)
		}

		go func(conn net.Conn) {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				payload := strings.Split(scanner.Text(), " ")

				switch strings.ToUpper(payload[0]) {
				case "GET":
					value, err := db.Get(payload[1])
					if err != nil {
						log.Println(err)
						break
					}

					conn.Write(value)
				case "SET":
					err := db.Set(payload[1], []byte(payload[2]))
					if err != nil {
						log.Println(err)
						break
					}
				case "DEL":
					err := db.Del(payload[1])
					if err != nil {
						log.Println(err)
						break
					}
				case "KEYS":
					keys, err := db.Keys()
					if err != nil {
						log.Println(err)
						break
					}

					for _, key := range keys {
						conn.Write([]byte(key))
					}
				case "QUIT":
					conn.Close()
					return
				}
			}
		}(conn)
	}
}

func NewHTTPServer(port string, db db.KVDB) {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		segments := strings.Split(req.RequestURI, "/")
		key := segments[1]

		switch req.Method {
		case http.MethodGet:
			if segments[1] == "" {
				keys, err := db.Keys()
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(err.Error()))
					return
				}

				res.WriteHeader(http.StatusOK)
				for _, key := range keys {
					res.Write([]byte(key + "\n"))
				}
				return
			}

			value, err := db.Get(key)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}

			res.WriteHeader(http.StatusOK)
			res.Write(value)
		case http.MethodPost:
			value, _ := ioutil.ReadAll(req.Body)
			defer req.Body.Close()

			err := db.Set(key, value)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}
		case http.MethodDelete:
			err := db.Del(key)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
			}
		default:
			res.WriteHeader(http.StatusBadRequest)
		}
	})

	http.ListenAndServe(port, nil)
}
