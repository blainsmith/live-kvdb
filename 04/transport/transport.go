package transport

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/blainsmith/live-kvdb/04/db"
)

func NewTCPServer(port string, db db.KeyValue) {
	tcp, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Println(err)
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			payload := strings.Split(scanner.Text(), " ")

			switch strings.ToUpper(payload[0]) {
			case "GET":
				value, err := db.Get(payload[1])
				if err != nil {
					log.Println(err)
					conn.Write([]byte(err.Error()))
					break
				}
				conn.Write(value)
			case "SET":
				err := db.Set(payload[1], []byte(payload[2]))
				if err != nil {
					log.Println(err)
					conn.Write([]byte(err.Error()))
					break
				}
				conn.Write([]byte("OK"))
			case "DEL":
				err := db.Del(payload[1])
				if err != nil {
					log.Println(err)
					conn.Write([]byte(err.Error()))
					break
				}
				conn.Write([]byte("OK"))
			case "KEYS":
				keys, err := db.Keys()
				if err != nil {
					log.Println(err)
					conn.Write([]byte(err.Error()))
					break
				}

				for _, key := range keys {
					conn.Write([]byte(key + "\n"))
				}
			case "QUIT":
				conn.Write([]byte("BYE!"))
				conn.Close()
				break
			}
		}
	}
}

func NewHTTPServer(port string, db db.KeyValue) {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		segments := strings.Split(req.RequestURI, "/")

		switch req.Method {
		case http.MethodGet:
			if segments[1] == "" {
				keys, err := db.Keys()
				if err != nil {
					log.Println(err)
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(err.Error()))
					return
				}

				for _, key := range keys {
					res.Write([]byte(key + "\n"))
				}
				res.WriteHeader(http.StatusOK)
				return
			}

			value, err := db.Get(segments[1])
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}

			res.WriteHeader(http.StatusOK)
			res.Write(value)
		case http.MethodPost:
			value, _ := ioutil.ReadAll(req.Body)
			err := db.Set(segments[1], value)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}
			res.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			err := db.Del(segments[1])
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}
			res.WriteHeader(http.StatusOK)
		default:
			res.WriteHeader(http.StatusBadRequest)
		}
	})

	http.ListenAndServe(port, nil)
}
