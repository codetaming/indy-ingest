package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Submission struct {
	owner        string `json:"owner"`
	submissionId string `json:"submission_id"`
}

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
	fmt.Println("Create Submission")

	var (
		tableName = aws.String(os.Getenv("SUBMISSIONS_TABLE"))
	)

	// Initialize product
	submission := &Submission{
		owner:        "",
		submissionId: "123",
	}

	// Write to DynamoDB
	item, _ := dynamodbattribute.MarshalMap(submission)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(submission)
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
