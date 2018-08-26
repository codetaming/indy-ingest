package main

import (
	"fmt"
	"github.com/codetaming/indy-ingest/cmd/ingestd/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func init() {
	fmt.Println("AWS_REGION:", os.Getenv("AWS_REGION"))
	fmt.Println("DATASET_TABLE:", os.Getenv("DATASET_TABLE"))
	fmt.Println("METADATA_TABLE:", os.Getenv("METADATA_TABLE"))
	fmt.Println("METADATA_BUCKET:", os.Getenv("METADATA_BUCKET"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/validate", handlers.Validate).Methods("POST")
	router.HandleFunc("/dataset", handlers.CreateDataset).Methods("POST")
	router.HandleFunc("/dataset", handlers.ListDatasets).Methods("GET")
	router.HandleFunc("/dataset/{datasetId}", handlers.GetDataset).Methods("GET")
	router.HandleFunc("/dataset/{datasetId}/metadata", handlers.AddMetadata).Methods("POST")
	router.HandleFunc("/dataset/{datasetId}/metadata", handlers.ListMetadata).Methods("GET")
	router.HandleFunc("/dataset/{datasetId}/metadata/{metadataId}", handlers.GetMetadata).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}
