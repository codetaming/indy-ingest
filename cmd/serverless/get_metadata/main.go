package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/codetaming/indy-ingest/internal/utils"
	"time"
)

//AWS Lambda entry point
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(request, new(persistence.DynamoDataStore), new(persistence.S3FileStore))
}

//Do executes the function allowing dependencies to be specified
func Do(request events.APIGatewayProxyRequest, p persistence.MetadataGetter, s persistence.MetadataRetriever) (events.APIGatewayProxyResponse, error) {
	datasetId := request.PathParameters["datasetId"]
	metadataId := request.PathParameters["metadataId"]
	metadataRecord, err := p.GetMetadata(datasetId, metadataId)
	if err != nil {
		return utils.RespondToError(err)
	}
	metadataContent, err := s.RetrieveMetadata(datasetId + "/" + metadataId)
	return respond(metadataRecord, metadataContent, err)
}

func respond(metadataRecord model.Metadata, metadataContent string, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return utils.RespondToError(err)
	}
	headers := map[string]string{
		"Content-Type": "application/json; schema=\"" + metadataRecord.DescribedBy + "\"",
		"Date":         metadataRecord.Created.Format(time.RFC1123),
	}
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(metadataContent),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
