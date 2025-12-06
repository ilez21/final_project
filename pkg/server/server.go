package server

import (
	"final_sprint/pkg/api"
	"log"
	"net/http"
	"os"
)

func Start() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}
