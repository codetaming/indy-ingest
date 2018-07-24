package utils

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/api/model"
)

func RespondToInternalError(err error) (events.APIGatewayProxyResponse, error) {
	return respond(500, err.Error())
}

func RespondToClientError(err error) (events.APIGatewayProxyResponse, error) {
	return respond(400, err.Error())
}

func RespondToNotFound(id string) (events.APIGatewayProxyResponse, error) {
	return respond(404, id+" not found")
}

func respond(code int, message string) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	errorMessage := model.ErrorMessage{Message: message}
	jsonErrorMessage, _ := json.Marshal(errorMessage)
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(jsonErrorMessage),
		StatusCode: code,
	}, nil
}
