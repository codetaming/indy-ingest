package main_test

import (
	"testing"

	"github.com/codetaming/indy-ingest/api/create_submission"

	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{
				Body: ""},
			expect: "\"\"",
			err:    nil,
		}}
	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
