package main

import (
	"github.com/codetaming/indy-ingest/cmd/ingestd/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/validate", handlers.Validate).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
