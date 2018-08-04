package persistence

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/codetaming/indy-ingest/internal/model"
)

type ErroredPersistence struct{}

func (ErroredPersistence) PersistDataset(dataset model.Dataset) (err error) {
	return awserr.New("Error", "an error has occurred", err)
}

func (ErroredPersistence) PersistMetadata(metadata model.Metadata) (err error) {
	return awserr.New("Error", "an error has occurred", err)
}

func (ErroredPersistence) CheckDatasetIdExists(datasetId string) (success bool, err error) {
	return false, awserr.New("Error", "an error has occurred", err)
}
