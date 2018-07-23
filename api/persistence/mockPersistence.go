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

func (MockPersistence) GetDataset(datasetId string) (model.Dataset, error) {
	return model.Dataset{
		DatasetId: datasetId,
	}, nil
}

func (MockPersistence) GetMetadata(datasetId string, metadataId string) (model.Metadata, error) {
	return model.Metadata{
		DatasetId:  datasetId,
		MetadataId: metadataId,
	}, nil
}
