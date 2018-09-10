package api

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/utils"
	"github.com/codetaming/indy-ingest/internal/validator"
	"io/ioutil"
	"log"
	"net/http"
)

func (api *API) Validate(w http.ResponseWriter, r *http.Request) {
	schemaUrl, err := utils.ExtractSchemaUrlArray(r.Header)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	result, err := validator.Validate(schemaUrl, string(b[:]))
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}
