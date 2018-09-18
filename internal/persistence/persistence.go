package persistence

import (
	"github.com/codetaming/indy-ingest/internal/model"
)

type DataStore interface {
	DatasetPersister
	MetadataPersister
	DatasetExistenceChecker
	DatasetLister
	DatasetGetter
	MetadataLister
	MetadataGetter
}

type FileStore interface {
	MetadataStorer
	MetadataRetriever
}

type MetadataStorer interface {
	StoreMetadata(key string, content string) (location string, error error)
}

type MetadataRetriever interface {
	RetrieveMetadata(key string) (string, error)
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string { return e.Message }

type DatasetPersister interface {
	PersistDataset(dataset model.Dataset) (err error)
}

type MetadataPersister interface {
	PersistMetadata(metadata model.Metadata) (err error)
}

type DatasetExistenceChecker interface {
	CheckDatasetIdExists(datasetId string) (bool, error)
}

type DatasetLister interface {
	ListDatasets() ([]model.Dataset, error)
}

type DatasetGetter interface {
	GetDataset(datasetId string) (model.Dataset, error)
}

type MetadataLister interface {
	ListMetadata(datasetId string) ([]model.Metadata, error)
}

type MetadataGetter interface {
	GetMetadata(datasetId string, metadataId string) (model.Metadata, error)
}
