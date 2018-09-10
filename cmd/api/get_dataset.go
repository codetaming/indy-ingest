package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetDataset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	dataset, err := api.dataStore.GetDataset(datasetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(dataset)
}
