package main

import (
	"os"
	"os/signal"

	"github.com/blainsmith/live-kvdb/04/db"
	"github.com/blainsmith/live-kvdb/04/db/memory"
	"github.com/blainsmith/live-kvdb/04/transport"
)

func main() {
	var db db.KeyValue

	db = &memory.DB{}

	go transport.NewTCPServer(":9090", db)
	go transport.NewHTTPServer(":8080", db)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
