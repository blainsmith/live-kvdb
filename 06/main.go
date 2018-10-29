package main

import (
	"github.com/blainsmith/live-kvdb/06/db/inmemory"
	"github.com/blainsmith/live-kvdb/06/transport"
	"os"
	"os/signal"
)

func main() {
	db := inmemory.New()

	go transport.NewTCPServer(db, ":1313")
	go transport.NewHTTPServer(db, ":8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
