package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("List Datasets")
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(Handler)
}
