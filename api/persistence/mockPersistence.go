package persistence

import (
	"github.com/codetaming/indy-ingest/api/model"
)

type MockPersistence struct{}

func (MockPersistence) PersistDataset(dataset model.Dataset) (err error) {
	return nil
}

func (MockPersistence) PersistMetadata(metadata model.Metadata) (err error) {
	return nil
}

func (MockPersistence) CheckDatasetIdExists(datasetId string) (bool, error) {
	return true, nil
}

func (MockPersistence) ListDatasets() ([]model.Dataset, error) {
	return []model.Dataset{}, nil
}
