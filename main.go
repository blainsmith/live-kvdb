package main

import (
	"os"
	"os/signal"

	"github.com/blainsmith/live-kvdb/db/inmemory"
	"github.com/blainsmith/live-kvdb/transport"
)

func main() {
	db := inmemory.New()

	go transport.NewTCPServer(db, ":9090")
	go transport.NewHTTPServer(db, ":8080")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
