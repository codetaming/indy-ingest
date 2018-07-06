package main

import (
	"errors"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/validator"
	"encoding/json"
)

var (
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	schemaUrl := request.Headers["describedBy"]
	bodyJson := request.Body
	result := validator.Validate(schemaUrl, bodyJson)
	body, _ := json.Marshal(result)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil

}


func main() {
	lambda.Start(Handler)
}
