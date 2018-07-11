package main

import (
	"github.com/aws/aws-lambda-go/events"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"time"
	"github.com/codetaming/indy-ingest/api/persistence"
	"github.com/google/uuid"
)

func Handler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return respond(createDataSet(new(persistence.DynamoPersistence)))
}

func MockHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return respond(createDataSet(new(persistence.MockPersistence)))
}

func createDataSet(p persistence.DatasetPersister) (model.Dataset, error) {
	d := model.Dataset{
		Owner:     model.DefaultOwner,
		DatasetId: uuid.Must(uuid.NewUUID()).String(),
		Created:   time.Now(),
	}
	return d, p.PersistDataset(d)
}

func respond(d model.Dataset, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(d)
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
