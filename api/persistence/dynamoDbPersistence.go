package persistence

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/codetaming/indy-ingest/api/model"
	"log"
	"os"
)

type DynamoPersistence struct{}

type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string { return e.msg }

var ddb *dynamodb.DynamoDB

const datasetTableEnv = "DATASET_TABLE"
const metadataTableEnv = "METADATA_TABLE"

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

	tableName := aws.String(os.Getenv(datasetTableEnv))

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

	tableName := aws.String(os.Getenv(metadataTableEnv))

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
		tableName = aws.String(os.Getenv(datasetTableEnv))
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

func (DynamoPersistence) ListDatasets() ([]model.Dataset, error) {
	var (
		tableName = aws.String(os.Getenv(datasetTableEnv))
	)
	var queryInput = &dynamodb.QueryInput{
		TableName: tableName,
		KeyConditions: map[string]*dynamodb.Condition{
			"owner": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(model.DefaultOwner),
					},
				},
			},
		},
	}
	result, err := ddb.Query(queryInput)
	if err != nil {
		return nil, err
	} else {
		var datasets []model.Dataset
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &datasets)
		return datasets, nil
	}
}

func (DynamoPersistence) GetDataset(datasetId string) (model.Dataset, error) {
	var tableName = aws.String(os.Getenv(datasetTableEnv))
	var dataset model.Dataset
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
		log.Println("Error retrieving:" + err.Error())
		return dataset, err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &dataset)
	if err != nil {
		log.Println("Error unmarshalling:" + err.Error())
		return dataset, err
	}
	if dataset.DatasetId == "" {
		return dataset, &NotFoundError{datasetId}
	}
	log.Println("Returning dataset successfully")
	return dataset, nil
}

func (DynamoPersistence) GetMetadata(datasetId string, metadataId string) (model.Metadata, error) {
	var (
		tableName = aws.String(os.Getenv(metadataTableEnv))
	)
	result, err := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"dataset_id": {
				S: aws.String(datasetId),
			},
			"metadata_id": {
				S: aws.String(metadataId),
			},
		},
	})
	log.Print("table: " + os.Getenv(metadataTableEnv))
	log.Print("dataset_id: " + datasetId + ", metadata_id: " + metadataId)
	if err != nil {
		return model.Metadata{}, err
	}
	metadata := model.Metadata{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &metadata)
	return metadata, nil

}
