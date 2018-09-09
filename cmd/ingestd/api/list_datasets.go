package api

import (
	"encoding/json"
	"net/http"
)

func (api *API) ListDatasets(w http.ResponseWriter, r *http.Request) {
	datasets, err := api.dataStore.ListDatasets()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(datasets)
}
