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
)

var ddb *dynamodb.DynamoDB

type Metadata struct {
	SubmissionId string    `dynamodbav:"submission_id"`
	MetadataId   string    `dynamodbav:"metadata_id"`
	DescribedBy  string    `dynamodbav:"described_by"`
	Created      time.Time `dynamodbav:"created" type:"timestamp" timestampFormat:"unix"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	submissionId := request.PathParameters["id"]
	schemaUrl := request.Headers["describedBy"]
	bodyJson := request.Body
	result := validator.Validate(schemaUrl, bodyJson)
	body, _ := json.Marshal(result)

	if result.Valid {
		metadata, _ := createMetadata(submissionId, schemaUrl)
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(metadata),
			StatusCode: 201,
		}, nil
	} else
	{
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
			StatusCode: 400,
		}, nil
	}
}

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session)
	}
}

func createMetadata(submissionId string, schemaUrl string) ([]byte, error) {
	log.Println("Create Metadata")

	u := uuid.Must(uuid.NewV4()).String()
	t := time.Now()

	s := Metadata{
		SubmissionId: submissionId,
		MetadataId:   u,
		DescribedBy:  schemaUrl,
		Created:      t,
	}

	av, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	for key, value := range av {
		log.Println("Key:", key, "Value:", value)
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
		return body, nil
	} else {
		return nil, err
	}
}

func main() {
	lambda.Start(Handler)
}
