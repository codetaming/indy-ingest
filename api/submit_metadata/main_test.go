package main_test

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/codetaming/indy-ingest/api/submit_metadata"
)

func TestHandler(t *testing.T) {

	headers := map[string]string{
		"describedBy": "https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism",
	}

	pathParameters := map[string]string{
		"id": "12345",
	}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expectedMessage string
		expectedCode int
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body

			request: events.APIGatewayProxyRequest{
				PathParameters: pathParameters,
				Headers: headers,
				Body: `{
    "organ": {
        "text": "brain",
        "ontology": "UBERON:0000955"
    },
    "schema_type": "biomaterial",
    "biomaterial_core": {
        "ncbi_taxon_id": [
            9606
        ],
        "biomaterial_id": "BT_S2_T",
        "has_input_biomaterial": "BT_S2",
        "biomaterial_description": "Tumor"
    },
    "organ_part": {
        "text": "temporal lobe"
    },
    "genus_species": [
        {
            "text": "Homo sapiens",
            "ontology": "NCBITaxon:9606"
        }
    ],
    "describedBy": "https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism"
}`},
			expectedMessage: "{\"SubmissionId\":\"12345\",\"MetadataId\":\"b74d1c3f-79e6-408a-9f40-a1c572bcb961\",\"DescribedBy\":\"https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism\",\"Created\":\"2018-07-07T12:40:41.206641545+01:00\"}",
			expectedCode: 201,
			err:    nil,
		},
		{
			request: events.APIGatewayProxyRequest{
				Headers: headers,
				Body:    "{}"},
			expectedMessage: "{\"Valid\":false,\"Message\":\"The document is not valid\",\"Errors\":[\"describedBy is required\",\"schema_type is required\",\"biomaterial_core is required\",\"organ is required\"]}",
			expectedCode: 400,
			err:    nil,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedMessage, response.Body)
		assert.Equal(t, test.expectedCode, response.StatusCode)
	}

}
