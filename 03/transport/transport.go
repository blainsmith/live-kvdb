package transport

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

type KVDB interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Del(string) error
	Keys() ([]string, error)
}

func StartTCPServer(port string, db KVDB) error {
	tcp, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Println(err)
		}
		addr := conn.LocalAddr()
		log.Printf("Accepted connection: %s", addr.String())

		go func(conn net.Conn) {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				payload := strings.Split(scanner.Text(), " ")

				switch strings.ToUpper(payload[0]) {
				case "GET":
					value, err := db.Get(payload[1])
					if err != nil {
						log.Fatal(err)
						conn.Write([]byte(err.Error()))
					}
					conn.Write(value)
					conn.Write([]byte("\n"))
				case "SET":
					err := db.Set(payload[1], []byte(payload[2]))
					if err != nil {
						log.Fatal(err)
						conn.Write([]byte(err.Error()))
					}
					conn.Write([]byte("OK\n"))
				case "DEL":
					err := db.Del(payload[1])
					if err != nil {
						log.Fatal(err)
						conn.Write([]byte(err.Error()))
					}
					conn.Write([]byte("OK\n"))
				case "KEYS":
					keys, err := db.Keys()
					if err != nil {
						log.Fatal(err)
						conn.Write([]byte(err.Error()))
					}
					for _, key := range keys {
						conn.Write([]byte(key))
						conn.Write([]byte("\n"))
					}
				case "QUIT":
					log.Printf("Closed connection: %s", addr.String())
					conn.Close()
					return
				}
			}
		}(conn)
	}
}

func StartHTTPServer(port string, db KVDB) error {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		segments := strings.Split(req.RequestURI, "/")
		switch req.Method {
		case http.MethodGet:
			if segments[1] != "" {
				value, err := db.Get(segments[1])
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(err.Error()))
					return
				}
				res.WriteHeader(http.StatusOK)
				res.Write(value)
				return
			}

			keys, err := db.Keys()
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
				return
			}
			res.WriteHeader(http.StatusOK)
			for _, key := range keys {
				res.Write([]byte(key))
				res.Write([]byte("\n"))
			}
			return
		case http.MethodPost:
			value, _ := ioutil.ReadAll(req.Body)
			err := db.Set(segments[1], value)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
			}
			res.WriteHeader(http.StatusAccepted)
			res.Write([]byte("OK\n"))
		case http.MethodDelete:
			err := db.Del(segments[1])
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
			}
			res.WriteHeader(http.StatusAccepted)
			res.Write([]byte("OK\n"))
		default:
			res.WriteHeader(http.StatusBadRequest)
		}
	})

	return http.ListenAndServe(port, nil)
}
