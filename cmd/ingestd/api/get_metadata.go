package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (api *API) GetMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	metadataId := vars["metadataId"]
	metadataRecord, err := api.dataStore.GetMetadata(datasetId, metadataId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	metadataContent, err := api.fileStore.RetrieveMetadata(datasetId + "/" + metadataId)
	w.Header().Set("content-type", "application/json; schema=\""+metadataRecord.DescribedBy+"\"")
	w.Header().Set("date", metadataRecord.Created.Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metadataContent))
}
