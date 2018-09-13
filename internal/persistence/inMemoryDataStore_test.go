package persistence

import (
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestInMemoryStore(t *testing.T) {
	logger := log.New(os.Stdout, "ingest-test ", log.LstdFlags|log.Lshortfile)
	dataStore, err := NewInMemoryDataStore(logger)
	assert.Nil(t, err)
	datasetId := uuid.Must(uuid.NewUUID()).String()
	dataset := model.Dataset{
		Owner:     "owner",
		DatasetId: datasetId,
		Created:   time.Now(),
	}
	dataStore.PersistDataset(dataset)
	retrievedDataset, err := dataStore.GetDataset(datasetId)
	assert.Nil(t, err)
	assert.Equal(t, datasetId, dataset.DatasetId)
	assert.Equal(t, dataset, retrievedDataset)
	datasets, err := dataStore.ListDatasets()
	assert.Nil(t, err)
	assert.Equal(t, dataset, datasets[0])
	exists, err := dataStore.CheckDatasetIdExists(datasetId)
	assert.True(t, exists)
	metadataUuid := uuid.Must(uuid.NewUUID()).String()
	metadata := model.Metadata{
		DatasetId:   datasetId,
		MetadataId:  metadataUuid,
		DescribedBy: "",
		Created:     time.Now(),
	}
	dataStore.PersistMetadata(metadata)
	retrievedMetadata, err := dataStore.GetMetadata(datasetId, metadataUuid)
	assert.Nil(t, err)
	assert.Equal(t, metadataUuid, metadata.MetadataId)
	assert.Equal(t, metadata, retrievedMetadata)
	metadatas, err := dataStore.ListMetadata(datasetId)
	assert.Nil(t, err)
	assert.Equal(t, metadata, metadatas[0])
}
