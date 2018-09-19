package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/internal/persistence/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerForCreateDataset(t *testing.T) {
	tests := []struct {
		request                events.APIGatewayProxyRequest
		expectedLocationHeader string
		expectedBody           string
		expectedStatusCode     int
		err                    error
	}{
		{
			request:                events.APIGatewayProxyRequest{Body: ""},
			expectedLocationHeader: "http://test/dataset/.+",
			expectedBody:           "{\"owner\":\".+\",\"dataset_id\":\".+\",\"created\":\".+\"}",
			expectedStatusCode:     201,
			err:                    nil,
		}}
	for _, test := range tests {
		response, err := MockHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Regexp(t, test.expectedLocationHeader, response.Headers["Location"])
		assert.Equal(t, test.expectedStatusCode, response.StatusCode)
		assert.Regexp(t, test.expectedBody, response.Body)
	}
}

func MockHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(new(mock.MockPersistence))
}
