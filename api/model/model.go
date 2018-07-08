package model

import "time"

type Dataset struct {
	Owner     string    `dynamodbav:"owner"`
	DatasetId string    `dynamodbav:"dataset_id"`
	Created   time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix"`
}

type Metadata struct {
	DatasetId   string    `dynamodbav:"dataset_id"`
	MetadataId  string    `dynamodbav:"metadata_id"`
	DescribedBy string    `dynamodbav:"described_by"`
	Created     time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

var DefaultOwner = "dan"

type MetadataSuccessMessage struct {
	Content Metadata `json:"content"`
	File    string `json:"file"`
}
