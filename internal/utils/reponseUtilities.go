package utils

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"log"
)

func RespondToError(err error) (events.APIGatewayProxyResponse, error) {
	switch t := err.(type) {
	case *persistence.NotFoundError:
		log.Println("NotFoundError", t)
		return RespondToNotFound(err.Error())
	default:
		return RespondToInternalError(err)
	}
}

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
