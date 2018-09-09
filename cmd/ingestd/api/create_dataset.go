package api

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

func (api *API) CreateDataset(w http.ResponseWriter, _ *http.Request) {
	d := model.Dataset{
		Owner:     model.DefaultOwner,
		DatasetId: uuid.Must(uuid.NewUUID()).String(),
		Created:   time.Now(),
	}
	err := api.dataStore.PersistDataset(d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	baseUrl := os.Getenv("BASE_URL")
	w.Header().Set("location", baseUrl+"/dataset/"+d.DatasetId)
	json.NewEncoder(w).Encode(d)
}
