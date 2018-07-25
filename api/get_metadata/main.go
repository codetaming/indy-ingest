package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/persistence"
	"github.com/codetaming/indy-ingest/api/storage"
	"github.com/codetaming/indy-ingest/api/utils"
)

//AWS Lambda entry point
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(request, new(persistence.DynamoPersistence), new(storage.MockStorage))
}

//Do executes the function allowing dependencies to be specified
func Do(request events.APIGatewayProxyRequest, p persistence.MetadataGetter, s storage.MetadataRetriever) (events.APIGatewayProxyResponse, error) {
	datasetId := request.PathParameters["datasetId"]
	metadataId := request.PathParameters["metadataId"]
	metadataRecord, err := p.GetMetadata(datasetId, metadataId)
	metadataContent, err := s.RetrieveMetadata(datasetId, metadataId)
	return respond(metadataRecord, metadataContent, err)
}

func respond(metadataRecord model.Metadata, metadataContent string, err error) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	if err != nil {
		utils.RespondToInternalError(err)
	}
	body, _ := json.Marshal(metadataRecord)
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
