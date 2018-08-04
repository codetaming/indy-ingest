package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/codetaming/indy-ingest/internal/publication"
	"github.com/codetaming/indy-ingest/internal/storage"
	"github.com/stretchr/testify/assert"
)

type testDefinition struct {
	request         events.APIGatewayProxyRequest
	expectedMessage string
	expectedCode    int
	err             error
}

var headers = map[string]string{
	"Link": `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>; rel="describedby"`,
}

var pathParameters = map[string]string{
	"datasetId": "12345",
}

func TestErroredHandler(t *testing.T) {
	tests := []testDefinition{
		{
			request: events.APIGatewayProxyRequest{
				PathParameters: pathParameters,
				Headers:        headers,
				Body:           ``},
			expectedMessage: "",
			expectedCode:    500,
			err:             nil,
		},
	}
	for _, test := range tests {
		response, err := ErroredHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedCode, response.StatusCode)
	}
}

func TestHandler(t *testing.T) {

	tests := []testDefinition{
		{
			request: events.APIGatewayProxyRequest{
				PathParameters: pathParameters,
				Headers:        headers,
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
			expectedMessage: "{\"info\":{\"dataset_id\":\"12345\",\"metadata_id\":\".+\",\"described_by\":\"https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism\",\"created\":\".+\"},\"file\":\"12345/\"}",
			expectedCode:    201,
			err:             nil,
		},
		{
			request: events.APIGatewayProxyRequest{
				Headers: headers,
				Body:    "{}"},
			expectedMessage: "{\"valid\":false,\"message\":\".+\",\"errors\":[\"describedBy is required\",\"schema_type is required\",\"biomaterial_core is required\",\"organ is required\"]}",
			expectedCode:    400,
			err:             nil,
		},
	}

	for _, test := range tests {
		response, err := MockHandler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedCode, response.StatusCode)
		if test.expectedCode == 201 {
			assert.Regexp(t, "http://test/dataset/12345/metadata/.+", response.Headers["Location"])
		}
	}
}

func ErroredHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := new(persistence.ErroredPersistence)
	s := new(storage.MockStorage)
	pub := new(publication.MockPublication)
	return Do(request, p, p, s, pub)
}

func MockHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := new(persistence.MockPersistence)
	s := new(storage.MockStorage)
	pub := new(publication.MockPublication)
	return Do(request, p, p, s, pub)
}
