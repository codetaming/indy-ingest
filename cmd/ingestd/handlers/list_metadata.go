package handlers

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

func ListMetadata(w http.ResponseWriter, r *http.Request) {
	p := new(persistence.DynamoPersistence)
	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	metadata, err := p.ListMetadata(datasetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}
