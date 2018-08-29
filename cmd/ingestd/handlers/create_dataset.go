package handlers

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

func CreateDataset(w http.ResponseWriter, _ *http.Request) {
	p := new(persistence.DynamoPersistence)
	d := model.Dataset{
		Owner:     model.DefaultOwner,
		DatasetId: uuid.Must(uuid.NewUUID()).String(),
		Created:   time.Now(),
	}
	err := p.PersistDataset(d)
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
