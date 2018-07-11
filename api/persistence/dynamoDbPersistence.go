package persistence

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"fmt"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoPersistence struct{}

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if ses, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		log.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(ses)
	}
}

func (DynamoPersistence) PersistDataset(dataset model.Dataset) (err error) {
	av, err := dynamodbattribute.MarshalMap(dataset)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	tableName := aws.String(os.Getenv("DATASET_TABLE"))

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return err
	}
	return nil
}

func (DynamoPersistence) PersistMetadata(metadata model.Metadata) (err error) {
	av, err := dynamodbattribute.MarshalMap(metadata)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		return err
	}

	tableName := aws.String(os.Getenv("METADATA_TABLE"))

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return err
	}
	return nil
}

func (DynamoPersistence) CheckDatasetIdExists(datasetId string) (bool, error) {
	var (
		tableName = aws.String(os.Getenv("DATASET_TABLE"))
	)
	result, err := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"owner": {
				S: aws.String(model.DefaultOwner),
			},
			"dataset_id": {
				S: aws.String(datasetId),
			},
		},
	})

	if err != nil {
		return false, err
	}

	dataset := model.Dataset{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &dataset)

	if err != nil {
		return false, err
	}

	if dataset.DatasetId == "" {
		return false, nil
	}

	return true, nil
}

