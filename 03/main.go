package main

import (
	"os"
	"os/signal"

	"github.com/blainsmith/live-kvdb/three/db"
	"github.com/blainsmith/live-kvdb/three/transport"
)

func main() {
	kvdb := db.InMemory{}

	go transport.StartTCPServer(":1313", &kvdb)
	go transport.StartHTTPServer(":8080", &kvdb)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
