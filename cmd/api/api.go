package api

import (
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/codetaming/indy-ingest/internal/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type API struct {
	logger    *log.Logger
	dataStore persistence.DataStore
	fileStore persistence.FileStore
	validator validator.Validator
}

func (api *API) SetupRoutes(r *mux.Router) {
	r.Handle("/", http.FileServer(http.Dir("./ui")))
	r.HandleFunc("/validate", api.Logger(api.Validate)).Methods("POST")
	r.HandleFunc("/dataset", api.Logger(api.CreateDataset)).Methods("POST")
	r.HandleFunc("/dataset", api.Logger(api.ListDatasets)).Methods("GET")
	r.HandleFunc("/dataset/{datasetId}", api.Logger(api.GetDataset)).Methods("GET")
	r.HandleFunc("/dataset/{datasetId}/metadata", api.Logger(api.AddMetadata)).Methods("POST")
	r.HandleFunc("/dataset/{datasetId}/metadata", api.Logger(api.ListMetadata)).Methods("GET")
	r.HandleFunc("/dataset/{datasetId}/metadata/{metadataId}", api.Logger(api.GetMetadata)).Methods("GET")
}

func (api *API) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer api.logger.Printf("request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewAPI(logger *log.Logger, dataStore persistence.DataStore, fileStore persistence.FileStore, validator validator.Validator) *API {
	return &API{
		logger:    logger,
		dataStore: dataStore,
		fileStore: fileStore,
		validator: validator,
	}
}
