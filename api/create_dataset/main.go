package main

import (
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/satori/go.uuid"
	"github.com/codetaming/indy-ingest/api/model"
	"time"
)

var ddb *dynamodb.DynamoDB

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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Create Dataset")

	u := uuid.Must(uuid.NewV4()).String()
	t := time.Now()

	s := model.Dataset{
		Owner:     model.DefaultOwner,
		DatasetId: u,
		Created:   t,
	}

	av, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	for key, value := range av {
		log.Println("Key:", key, "Value:", value)
	}

	var (
		tableName = aws.String(os.Getenv("DATASET_TABLE"))
	)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(s)
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
