package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	jsonMessage, _ := json.Marshal("This is a dummy endpoint")
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(jsonMessage),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
