package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/persistence"
)

//AWS Lambda entry point
func Handler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(new(persistence.DynamoPersistence))
}

//Do executes the function allowing dependencies to be specified
func Do(p persistence.DatasetLister) (events.APIGatewayProxyResponse, error) {
	return respond(p.ListDatasets())
}

func respond(datasets []model.Dataset, err error) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	body, _ := json.Marshal(datasets)
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
