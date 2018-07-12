package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/persistence"
	"github.com/codetaming/indy-ingest/api/storage"
	"github.com/codetaming/indy-ingest/api/validator"
	"github.com/google/uuid"
	"time"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := new(persistence.DynamoPersistence)
	s := new(storage.S3Storage)
	return Do(request, p, p, s)
}

func Do(request events.APIGatewayProxyRequest, dec persistence.DatasetExistenceChecker, mp persistence.MetadataPersister, ms storage.MetadataStorer) (events.APIGatewayProxyResponse, error) {
	datasetId := request.PathParameters["id"]
	exists, err := checkDatasetExists(datasetId, dec)

	headers := map[string]string{"Content-Type": "application/json"}
	if err != nil {
		errorMessage := model.ErrorMessage{Message: err.Error()}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 500,
		}, nil
	}

	if !exists {
		errorMessage := model.ErrorMessage{Message: datasetId + " not found"}
		jsonErrorMessage, _ := json.Marshal(errorMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonErrorMessage),
			StatusCode: 404,
		}, nil
	}

	schemaUrl := request.Headers["describedBy"]
	bodyJson := request.Body

	result := validator.Validate(schemaUrl, bodyJson)

	if result.Valid {
		metadataRecord, metadataId, err := createMetadataRecord(datasetId, schemaUrl, mp)
		if err != nil {
			errorMessage := model.ErrorMessage{Message: err.Error()}
			jsonErrorMessage, _ := json.Marshal(errorMessage)
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				Body:       string(jsonErrorMessage),
				StatusCode: 500,
			}, nil
		}
		fileLocation, err := createMetadataFile(datasetId, metadataId, bodyJson, ms)
		if err != nil {
			errorMessage := model.ErrorMessage{Message: err.Error()}
			jsonErrorMessage, _ := json.Marshal(errorMessage)
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				Body:       string(jsonErrorMessage),
				StatusCode: 500,
			}, nil
		}
		metadataSuccessMessage := model.MetadataSuccessMessage{
			Info: metadataRecord,
			File: fileLocation,
		}
		jsonMetadataSuccessMessage, _ := json.Marshal(metadataSuccessMessage)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonMetadataSuccessMessage),
			StatusCode: 201,
		}, nil
	} else {
		validationResultJson, _ := json.Marshal(result)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(validationResultJson),
			StatusCode: 400,
		}, nil
	}
}

func checkDatasetExists(datasetId string, p persistence.DatasetExistenceChecker) (bool, error) {
	return p.CheckDatasetIdExists(datasetId)
}

func createMetadataFile(datasetId string, metadataId string, bodyJson string, ms storage.MetadataStorer) (fileLocation string, err error) {
	key := datasetId + "/" + metadataId
	return ms.StoreMetadata(key, bodyJson)
}

func createMetadataRecord(datasetId string, schemaUrl string, mp persistence.MetadataPersister) (metadataRecord model.Metadata, metadataId string, err error) {
	uuid := uuid.Must(uuid.NewUUID()).String()
	m := model.Metadata{
		DatasetId:   datasetId,
		MetadataId:  uuid,
		DescribedBy: schemaUrl,
		Created:     time.Now(),
	}
	persistErr := mp.PersistMetadata(m)
	if persistErr != nil {
		return m, uuid, nil
	} else {
		return m, "", persistErr
	}
}

func main() {
	lambda.Start(Handler)
}
