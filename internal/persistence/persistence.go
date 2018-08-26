package persistence

import (
	"github.com/codetaming/indy-ingest/internal/model"
)

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
