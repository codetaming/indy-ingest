package persistence

import (
	"github.com/codetaming/indy-ingest/api/model"
)

type DatasetPersister interface {
	PersistDataset(dataset model.Dataset) (err error)
}

type MetadataPersister interface {
	PersistMetadata(metadata model.Metadata)
}

type DatasetExistenceChecker interface {
	CheckDatasetIdExists(datasetId string) (bool, error)
}
