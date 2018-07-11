package main_test

import (
	"testing"

	"github.com/codetaming/indy-ingest/api/create_dataset"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "{\"owner\":\".+\",\"dataset_id\":\".+\",\"created\":\".+\"}",
			err:     nil,
		}}
	for _, test := range tests {
		response, err := main.MockHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Regexp(t, test.expect, response.Body)
	}
}
