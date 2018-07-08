package main

import (
	"github.com/aws/aws-lambda-go/events"
	"log"
	"github.com/codetaming/indy-ingest/api/validator"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/satori/go.uuid"
	"time"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strings"
)

var ddb *dynamodb.DynamoDB
var s3Uploader *s3manager.Uploader

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	submissionId := request.PathParameters["id"]

	exists, err := checkSubmissionIdExists(submissionId)

	headers := map[string]string{"Content-Type": "application/json"}

	if err != nil {
		errorMessage := model.ErrorMessage{Message: err.Error()}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 500,
		}, nil
	}

	if !exists {
		errorMessage := model.ErrorMessage{Message: submissionId + " not found"}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 404,
		}, nil
	}

	schemaUrl := request.Headers["describedBy"]
	bodyJson := request.Body

	result := validator.Validate(schemaUrl, bodyJson)

	if result.Valid {
		metadataRecord, metadataId, err := createMetadataRecord(submissionId, schemaUrl)
		if err != nil {
			errorMessage := model.ErrorMessage{Message: err.Error()}
			jsonErrorMessage, _ := json.Marshal(errorMessage)
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				Body:       string(jsonErrorMessage),
				StatusCode: 500,
			}, nil
		}
		fileLocation, err := createMetadataFile(submissionId, metadataId, bodyJson)
		if err != nil {
			errorMessage := model.ErrorMessage{Message: err.Error()}
			jsonErrorMessage, _ := json.Marshal(errorMessage)
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				Body:       string(jsonErrorMessage),
				StatusCode: 500,
			}, nil
		}
		metadataSuccessMessage := model.MetadataSuccessMessage{
			Content: metadataRecord,
			File: fileLocation,
		}
		jsonMetadataSuccessMessage, _ := json.Marshal(metadataSuccessMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonMetadataSuccessMessage),
			StatusCode: 201,
		}, nil
	} else
	{
		validationResultJson, _ := json.Marshal(result)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(validationResultJson),
			StatusCode: 400,
		}, nil
	}
}

func createMetadataFile(submissionId string, metadataId, bodyJson string) (fileLocation string, err error) {
	key := submissionId + "/" + metadataId

	fmt.Println("Uploading file to S3...")
		upParams := &s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("METADATA_BUCKET")),
		Key:    &key,
		Body:   strings.NewReader(bodyJson),
	}
	result, err := s3Uploader.Upload(upParams)
	if err != nil {
		panic(fmt.Sprintf("failed to create S3 file, %v", err))
		return "", err
	}
	return result.Location, nil
}

func checkSubmissionIdExists(submissionId string) (bool, error) {

	var (
		tableName = aws.String(os.Getenv("SUBMISSIONS_TABLE"))
	)
	result, err := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"owner": {
				S: aws.String(model.DefaultOwner),
			},
			"submission_id": {
				S: aws.String(submissionId),
			},
		},
	})

	if err != nil {
		return false, err
	}

	submission := model.Submission{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &submission)

	if err != nil {
		return false, err
	}

	if submission.SubmissionId == "" {
		return false, nil
	}

	return true, nil
}

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session)
		s3Uploader = s3manager.NewUploader(session)
	}
}

func createMetadataRecord(submissionId string, schemaUrl string) (metadataRecord []byte, metadataId string, err error) {
	log.Println("Create Metadata")

	u := uuid.Must(uuid.NewV4()).String()
	t := time.Now()

	s := model.Metadata{
		SubmissionId: submissionId,
		MetadataId:   u,
		DescribedBy:  schemaUrl,
		Created:      t,
	}

	av, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		return nil, "", err
	}

	var (
		tableName = aws.String(os.Getenv("METADATA_TABLE"))
	)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		body, _ := json.Marshal(s)
		return body, u, nil
	} else {
		return nil, "", err
	}
}

func main() {
	lambda.Start(Handler)
}
