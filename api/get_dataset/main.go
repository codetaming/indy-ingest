package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/persistence"
	"github.com/codetaming/indy-ingest/api/utils"
)

//AWS Lambda entry point
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(request, new(persistence.DynamoPersistence))
}

//Do executes the function allowing dependencies to be specified
func Do(request events.APIGatewayProxyRequest, p persistence.DatasetGetter) (events.APIGatewayProxyResponse, error) {
	return respond(p.GetDataset(request.PathParameters["datasetId"]))
}

func respond(dataset model.Dataset, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		utils.RespondToInternalError(err)
	}
	body, _ := json.Marshal(dataset)
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
