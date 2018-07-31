package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/persistence"
	"github.com/codetaming/indy-ingest/api/publication"
	"github.com/codetaming/indy-ingest/api/storage"
	"github.com/codetaming/indy-ingest/api/utils"
	"github.com/codetaming/indy-ingest/api/validator"
	"github.com/google/uuid"
	"os"
	"time"
)

//AWS Lambda entry point
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := new(persistence.DynamoPersistence)
	s := new(storage.S3Storage)
	pub := new(publication.SnsPublication)
	return Do(request, p, p, s, pub)
}

//Do executes the function allowing dependencies to be specified
func Do(request events.APIGatewayProxyRequest, dec persistence.DatasetExistenceChecker, mp persistence.MetadataPersister, ms storage.MetadataStorer, mcp publication.MetadataCreatedPublisher) (events.APIGatewayProxyResponse, error) {
	datasetId := request.PathParameters["datasetId"]
	_, err := checkDatasetExists(datasetId, dec)

	headers := map[string]string{"Content-Type": "application/json"}

	if err != nil {
		return utils.RespondToError(err)
	}

	schemaUrl, err := utils.ExtractSchemaUrl(request.Headers)
	if err != nil {
		return utils.RespondToClientError(err)
	}

	bodyJson := request.Body

	result, err := validator.Validate(schemaUrl, bodyJson)
	if err != nil {
		return utils.RespondToInternalError(err)
	}

	if result.Valid {

		metadataRecord, metadataId, err := createMetadataRecord(datasetId, schemaUrl, mp)
		if err != nil {
			return utils.RespondToInternalError(err)
		}

		fileLocation, err := createMetadataFile(datasetId, metadataId, bodyJson, ms)
		if err != nil {
			return utils.RespondToInternalError(err)
		}
		metadataSuccessMessage := model.MetadataSuccessMessage{
			Info: metadataRecord,
			File: fileLocation,
		}
		baseUrl := os.Getenv("BASE_URL")
		metadataUrl := baseUrl + "/dataset/" + datasetId + "/metadata/" + metadataId
		headers["Location"] = metadataUrl
		jsonMetadataSuccessMessage, _ := json.Marshal(metadataSuccessMessage)
		mcp.PublishMetadataCreated(metadataUrl)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string(jsonMetadataSuccessMessage),
			StatusCode: 201,
		}, nil
	}
	validationResultJson, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(validationResultJson),
		StatusCode: 400,
	}, nil
}

func checkDatasetExists(datasetId string, p persistence.DatasetExistenceChecker) (bool, error) {
	return p.CheckDatasetIdExists(datasetId)
}

func createMetadataFile(datasetId string, metadataId string, bodyJson string, ms storage.MetadataStorer) (fileLocation string, err error) {
	key := datasetId + "/" + metadataId
	return ms.StoreMetadata(key, bodyJson)
}

func createMetadataRecord(datasetID string, schemaUrl string, mp persistence.MetadataPersister) (metadataRecord model.Metadata, metadataId string, err error) {
	metadataUuid := uuid.Must(uuid.NewUUID()).String()
	m := model.Metadata{
		DatasetId:   datasetID,
		MetadataId:  metadataUuid,
		DescribedBy: schemaUrl,
		Created:     time.Now(),
	}
	persistErr := mp.PersistMetadata(m)
	if persistErr != nil {
		return m, "", persistErr
	}
	return m, metadataUuid, nil
}

func main() {
	lambda.Start(Handler)
}
