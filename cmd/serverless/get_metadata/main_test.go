package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/internal/persistence/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var headers = map[string]string{}

var pathParameters = map[string]string{
	"datasetId":  "12345",
	"metadataId": "67890",
}

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{
				PathParameters: pathParameters,
				Headers:        headers,
				Body:           ``},
			expect: "{}",
			err:    nil,
		}}
	for _, test := range tests {
		response, err := MockHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

func MockHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(request, new(mock.MockPersistence), new(mock.MockStorage))
}
