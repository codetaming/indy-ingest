package model

import "time"

type Submission struct {
	Owner        string    `dynamodbav:"owner"`
	SubmissionId string    `dynamodbav:"submission_id"`
	Created      time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix"`
}

type Metadata struct {
	SubmissionId string    `dynamodbav:"submission_id"`
	MetadataId   string    `dynamodbav:"metadata_id"`
	DescribedBy  string    `dynamodbav:"described_by"`
	Created      time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

var DefaultOwner = "dan"

type MetadataSuccessMessage struct {
	Content []byte `json:"content"`
	File    string `json:"file"`
}
