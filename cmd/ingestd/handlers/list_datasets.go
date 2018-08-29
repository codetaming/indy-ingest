package handlers

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"net/http"
)

func ListDatasets(w http.ResponseWriter, r *http.Request) {
	p := new(persistence.DynamoPersistence)
	datasets, err := p.ListDatasets()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(datasets)
}
