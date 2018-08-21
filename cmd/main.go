package main

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/utils"
	"github.com/codetaming/indy-ingest/internal/validator"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/validate", Validate).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func Validate(writer http.ResponseWriter, request *http.Request) {
	schemaUrl, err := utils.ExtractSchemaUrlArray(request.Header)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	b, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	result, err := validator.Validate(schemaUrl, string(b[:]))
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, err.Error(), 500)
		return
	}

	writer.Header().Set("content-type", "application/json")
	json.NewEncoder(writer).Encode(result)
}
