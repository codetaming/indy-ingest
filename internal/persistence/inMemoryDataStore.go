package persistence

import (
	"github.com/codetaming/indy-ingest/internal/model"
	"log"
)

type InMemoryDataStore struct {
	logger    *log.Logger
	datasets  map[string]model.Dataset
	metadatas map[string]model.Metadata
}

func (ds *InMemoryDataStore) PersistDataset(dataset model.Dataset) (err error) {
	ds.datasets[dataset.DatasetId] = dataset
	return nil
}

func (ds *InMemoryDataStore) PersistMetadata(metadata model.Metadata) (err error) {
	ds.metadatas[metadata.MetadataId] = metadata
	return nil
}

func (ds *InMemoryDataStore) CheckDatasetIdExists(datasetId string) (bool, error) {
	_, exists := ds.datasets[datasetId]
	return exists, nil
}

func (ds *InMemoryDataStore) ListDatasets() ([]model.Dataset, error) {
	var arr []model.Dataset
	for _, v := range ds.datasets {
		arr = append(arr, v)
	}
	return arr, nil
}

func (ds *InMemoryDataStore) GetDataset(datasetId string) (model.Dataset, error) {
	return ds.datasets[datasetId], nil
}

func (ds *InMemoryDataStore) ListMetadata(datasetId string) ([]model.Metadata, error) {
	var arr []model.Metadata
	for _, v := range ds.metadatas {
		arr = append(arr, v)
	}
	return arr, nil
}

func (ds *InMemoryDataStore) GetMetadata(datasetId string, metadataId string) (model.Metadata, error) {
	return ds.metadatas[metadataId], nil
}

func NewInMemoryDataStore(logger *log.Logger) (DataStore, error) {
	return &InMemoryDataStore{
		logger:    logger,
		datasets:  make(map[string]model.Dataset),
		metadatas: make(map[string]model.Metadata),
	}, nil
}
