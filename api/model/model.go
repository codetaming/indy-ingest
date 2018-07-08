package model

import "time"

type Dataset struct {
	Owner     string    `dynamodbav:"owner",json:"owner"`
	DatasetId string    `dynamodbav:"dataset_id",json:"dataset_id"`
	Created   time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix",json:"created"`
}

type Metadata struct {
	DatasetId   string    `dynamodbav:"dataset_id",json:"dataset_id"`
	MetadataId  string    `dynamodbav:"metadata_id",json:"metadata_id"`
	DescribedBy string    `dynamodbav:"described_by",json:"described_by"`
	Created     time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix",json:"created"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

var DefaultOwner = "dan"

type MetadataSuccessMessage struct {
	Info Metadata `json:"info"`
	File string   `json:"file"`
}
