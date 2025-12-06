package main

import (
	"final_sprint/pkg/db"
	"final_sprint/pkg/server"
	"log"
	"os"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
