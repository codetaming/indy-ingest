package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) Dummy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := vars
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}
