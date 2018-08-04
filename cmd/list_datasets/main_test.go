package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{},
			expect:  "[]",
			err:     nil,
		}}
	for _, test := range tests {
		response, err := MockHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

func MockHandler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(new(persistence.MockPersistence))
}
