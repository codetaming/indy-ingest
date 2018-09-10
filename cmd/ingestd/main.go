package main

import (
	"github.com/codetaming/indy-ingest/cmd/ingestd/api"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	serverPort     = os.Getenv("SERVER_PORT")
	region         = os.Getenv("AWS_REGION")
	datasetTable   = os.Getenv("DATASET_TABLE")
	metadataTable  = os.Getenv("METADATA_TABLE")
	metadataBucket = os.Getenv("METADATA_BUCKET")
)

func init() {
	if serverPort == "" {
		log.Fatal("$SERVER_PORT not set")
	}
	if region == "" {
		log.Fatal("$AWS_REGION not set")
	}
	if datasetTable == "" {
		log.Fatal("$DATASET_TABLE not set")
	}
	if metadataTable == "" {
		log.Fatal("$METADATA_TABLE not set")
	}
	if metadataBucket == "" {
		log.Fatal("$METADATA_BUCKET not set")
	}
}

func main() {
	router := mux.NewRouter()

	logger := log.New(os.Stdout, "ingest ", log.LstdFlags|log.Lshortfile)

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

	logger.Printf("server starting on port %s", serverPort)
	err = http.ListenAndServe(":"+serverPort, router)
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}
