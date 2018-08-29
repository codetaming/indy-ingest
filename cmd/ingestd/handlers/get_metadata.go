package handlers

import (
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/codetaming/indy-ingest/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func GetMetadata(w http.ResponseWriter, r *http.Request) {
	p := new(persistence.DynamoPersistence)
	s := new(storage.S3Storage)
	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	metadataId := vars["metadataId"]
	metadataRecord, err := p.GetMetadata(datasetId, metadataId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	metadataContent, err := s.RetrieveMetadata(datasetId + "/" + metadataId)
	w.Header().Set("content-type", "application/json; schema=\""+metadataRecord.DescribedBy+"\"")
	w.Header().Set("date", metadataRecord.Created.Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metadataContent))
}
