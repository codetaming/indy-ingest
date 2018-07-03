package main

import (
	"errors"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xeipuuv/gojsonschema"
	"encoding/json"
)

type Response struct {
	Message string `json:"message"`
}

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	schemaUrl := request.Headers["describedBy"]
	bodyJson := request.Body
	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(bodyJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	var resultStr string

	if result.Valid() {
		resultStr = "The document is valid\n"
	} else {
		resultStr = "The document is not valid. see errors :\n"
		for _, desc := range result.Errors() {
			resultStr = resultStr + "- %s\n" + desc.Description()
		}
	}

	body, _ := json.Marshal(resultStr)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
