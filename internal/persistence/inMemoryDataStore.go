package persistence

import (
	"github.com/codetaming/indy-ingest/internal/model"
	"log"
)

type InMemoryDataStore struct {
	logger *log.Logger
}

func (InMemoryDataStore) PersistDataset(dataset model.Dataset) (err error) {
	return nil
}

func (InMemoryDataStore) PersistMetadata(metadata model.Metadata) (err error) {
	panic("implement me")
}

func (InMemoryDataStore) CheckDatasetIdExists(datasetId string) (bool, error) {
	panic("implement me")
}

func (InMemoryDataStore) ListDatasets() ([]model.Dataset, error) {
	panic("implement me")
}

func (InMemoryDataStore) GetDataset(datasetId string) (model.Dataset, error) {
	panic("implement me")
}

func (InMemoryDataStore) ListMetadata(datasetId string) ([]model.Metadata, error) {
	panic("implement me")
}

func (InMemoryDataStore) GetMetadata(datasetId string, metadataId string) (model.Metadata, error) {
	panic("implement me")
}

func NewInMemoryDataStore(logger *log.Logger) (DataStore, error) {
	return InMemoryDataStore{
		logger: logger,
	}, nil
}
