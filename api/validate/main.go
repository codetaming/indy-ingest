package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/validator"
	"github.com/tomnomnom/linkheader"
	"log"
)

var (
	ErrMetadataNotProvided = errors.New("no metadata was provided in the HTTP body")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print(request)

	headers := map[string]string{"Content-Type": "application/json"}

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrMetadataNotProvided
	}

	linkHeader := request.Headers["Link"]
	links := linkheader.Parse(linkHeader)
	var link linkheader.Link
	if len(links) == 1 {
		link = links[0]
	}

	if link.URL == "" || link.Rel != "describedby" {
		errorMessage := model.ErrorMessage{Message: "Link header is invalid, not unique or missing"}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 500,
		}, nil
	}

	bodyJson := request.Body
	result, err := validator.Validate(link.URL, bodyJson)

	if err != nil {
		errorMessage := model.ErrorMessage{Message: err.Error()}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 500,
		}, nil
	}

	body, err := json.Marshal(result)

	if err != nil {
		errorMessage := model.ErrorMessage{Message: err.Error()}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
