package persistence

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/codetaming/indy-ingest/internal/model"
	"log"
)

type DynamoDataStore struct {
	logger        *log.Logger
	ddb           *dynamodb.DynamoDB
	datasetTable  *string
	metadataTable *string
}

func NewDynamoDataStore(logger *log.Logger, region string, datasetTable string, metadataTable string) (DataStore, error) {
	if ses, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		logger.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
		return nil, err
	} else {
		return DynamoDataStore{
			logger:        logger,
			ddb:           dynamodb.New(ses),
			datasetTable:  &datasetTable,
			metadataTable: &metadataTable,
		}, nil
	}
}

func (d DynamoDataStore) PersistDataset(dataset model.Dataset) (err error) {
	av, err := dynamodbattribute.MarshalMap(dataset)
	if err != nil {
		d.logger.Panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: d.datasetTable,
	}
	if _, err := d.ddb.PutItem(input); err != nil {
		return err
	}
	return nil
}

func (d DynamoDataStore) PersistMetadata(metadata model.Metadata) (err error) {
	av, err := dynamodbattribute.MarshalMap(metadata)
	if err != nil {
		d.logger.Panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: d.metadataTable,
	}
	if _, err := d.ddb.PutItem(input); err != nil {
		return err
	}
	return nil
}

func (d DynamoDataStore) CheckDatasetIdExists(datasetId string) (bool, error) {
	result, err := d.ddb.GetItem(&dynamodb.GetItemInput{
		TableName: d.datasetTable,
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
	if len(result.Item) == 0 {
		return false, &NotFoundError{datasetId}
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &dataset)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d DynamoDataStore) ListDatasets() ([]model.Dataset, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: d.datasetTable,
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
	result, err := d.ddb.Query(queryInput)
	if err != nil {
		return nil, err
	} else {
		var datasets []model.Dataset
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &datasets)
		return datasets, nil
	}
}

func (d DynamoDataStore) GetDataset(datasetId string) (model.Dataset, error) {
	var dataset model.Dataset
	result, err := d.ddb.GetItem(&dynamodb.GetItemInput{
		TableName: d.datasetTable,
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
		return dataset, err
	}
	if len(result.Item) == 0 {
		return dataset, &NotFoundError{datasetId}
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &dataset)
	if err != nil {
		d.logger.Println("Error unmarshalling:" + err.Error())
		return dataset, err
	}
	return dataset, nil
}

func (d DynamoDataStore) ListMetadata(datasetId string) ([]model.Metadata, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: d.metadataTable,
		KeyConditions: map[string]*dynamodb.Condition{
			"dataset_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(datasetId),
					},
				},
			},
		},
	}
	result, err := d.ddb.Query(queryInput)
	if err != nil {
		return nil, err
	} else {
		var metadata []model.Metadata
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &metadata)
		return metadata, nil
	}
}

func (d DynamoDataStore) GetMetadata(datasetId string, metadataId string) (model.Metadata, error) {
	var metadata model.Metadata
	result, err := d.ddb.GetItem(&dynamodb.GetItemInput{
		TableName: d.metadataTable,
		Key: map[string]*dynamodb.AttributeValue{
			"dataset_id": {
				S: aws.String(datasetId),
			},
			"metadata_id": {
				S: aws.String(metadataId),
			},
		},
	})
	d.logger.Print("dataset_id: " + datasetId + ", metadata_id: " + metadataId)

	if err != nil {
		return metadata, err
	}
	if len(result.Item) == 0 {
		return metadata, &NotFoundError{metadataId}
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &metadata)
	if err != nil {
		d.logger.Println("Error unmarshalling:" + err.Error())
		return metadata, err
	}
	return metadata, nil

}
