package handlers

import (
	"encoding/json"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/codetaming/indy-ingest/internal/storage"
	"github.com/codetaming/indy-ingest/internal/utils"
	"github.com/codetaming/indy-ingest/internal/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func AddMetadata(w http.ResponseWriter, r *http.Request) {
	p := new(persistence.DynamoPersistence)
	s := new(storage.S3Storage)
	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	_, err := checkDatasetExists(datasetId, p)

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
	w.Header().Set("content-type", "application/json")

	if result.Valid {
		metadataRecord, metadataId, err := createMetadataRecord(datasetId, schemaUrl, p)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fileLocation, err := createMetadataFile(datasetId, metadataId, string(b[:]), s)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		metadataSuccessMessage := model.MetadataSuccessMessage{
			Info: metadataRecord,
			File: fileLocation,
		}
		baseUrl := os.Getenv("BASE_URL")
		metadataUrl := baseUrl + "/dataset/" + datasetId + "/metadata/" + metadataId

		w.Header().Set("location", metadataUrl)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(metadataSuccessMessage)
	} else {
		validationResultJson, _ := json.Marshal(result)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationResultJson)
	}
}

func checkDatasetExists(datasetId string, p persistence.DatasetExistenceChecker) (bool, error) {
	return p.CheckDatasetIdExists(datasetId)
}

func createMetadataFile(datasetId string, metadataId string, bodyJson string, ms storage.MetadataStorer) (fileLocation string, err error) {
	key := datasetId + "/" + metadataId
	return ms.StoreMetadata(key, bodyJson)
}

func createMetadataRecord(datasetID string, schemaUrl string, mp persistence.MetadataPersister) (metadataRecord model.Metadata, metadataId string, err error) {
	metadataUuid := uuid.Must(uuid.NewUUID()).String()
	m := model.Metadata{
		DatasetId:   datasetID,
		MetadataId:  metadataUuid,
		DescribedBy: schemaUrl,
		Created:     time.Now(),
	}
	persistErr := mp.PersistMetadata(m)
	if persistErr != nil {
		return m, "", persistErr
	}
	return m, metadataUuid, nil
}
