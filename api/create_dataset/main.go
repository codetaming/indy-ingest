package main

import (
	"github.com/aws/aws-lambda-go/events"
	"log"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/satori/go.uuid"
	"github.com/codetaming/indy-ingest/api/model"
	"time"
	"github.com/codetaming/indy-ingest/api/persistence"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Create Dataset")

	u := uuid.Must(uuid.NewV4()).String()
	t := time.Now()

	d := model.Dataset{
		Owner:     model.DefaultOwner,
		DatasetId: u,
		Created:   t,
	}

	err := persistence.PersistDataset(d)
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
