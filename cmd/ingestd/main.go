package main

import (
	"github.com/codetaming/indy-ingest/cmd/ingestd/api"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	logger := log.New(os.Stdout, "ingest ", log.LstdFlags|log.Lshortfile)

	region := os.Getenv("AWS_REGION")
	datasetTable := os.Getenv("DATASET_TABLE")
	metadataTable := os.Getenv("METADATA_TABLE")
	metadataBucket := os.Getenv("METADATA_BUCKET")

	dataStore, err := persistence.NewDynamoPersistence(logger, region, datasetTable, metadataTable)
	if err != nil {
		logger.Fatalf("failed to create data store: %v", err)
	}

	fileStore, err := persistence.NewS3FileStore(logger, region, metadataBucket)
	if err != nil {
		logger.Fatalf("failed to create file store: %v", err)
	}

	a := api.NewAPI(logger, dataStore, fileStore)
	a.SetupRoutes(router)

	logger.Println("server starting")
	err = http.ListenAndServe(":9000", router)
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}
