package main_test

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/api/validate"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	headers := map[string]string{
		"describedBy": "https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism",
	}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body

			request: events.APIGatewayProxyRequest{
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
			expect: "{\"valid\":true,\"message\":\"The document is valid\",\"errors\":null}",
			err:    nil,
		},
		{
			request: events.APIGatewayProxyRequest{
				Headers: headers,
				Body:    "{}"},
			expect: "{\"valid\":false,\"message\":\"The document is not valid\",\"errors\":[\"describedBy is required\",\"schema_type is required\",\"biomaterial_core is required\",\"organ is required\"]}",
			err:    nil,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
