package main

import (
	"errors"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xeipuuv/gojsonschema"
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
	result := validate(schemaUrl, bodyJson)
	body, _ := json.Marshal(result)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil

}

type ValidationResult struct {
	Valid   bool
	Message string
	Errors  []string
}

func validate(schemaUrl string, bodyJson string) (ValidationResult) {

	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(bodyJson)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		panic(err.Error())
	}

	var message string
	var errors []string

	if result.Valid() {
		message = "The document is valid"
	} else {
		message = "The document is not valid"
		for _, desc := range result.Errors() {
			errors = append(errors, desc.Description())
		}
	}

	vr := ValidationResult{
		Valid:   result.Valid(),
		Message: message,
		Errors:  errors,
	}

	return vr
}

func main() {
	lambda.Start(Handler)
}
