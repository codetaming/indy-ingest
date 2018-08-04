package validator

import (
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		schemaUrl string
		json      string
		expect    model.ValidationResult
		err       error
	}{
		{
			schemaUrl: "https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism",
			json: `{
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
}`,
			expect: model.ValidationResult{
				Valid:   true,
				Message: "The document is valid",
				Errors:  nil,
			},
			err: nil,
		},
		{
			schemaUrl: "",
			json:      "",
			expect: model.ValidationResult{
				Valid:   false,
				Message: "",
				Errors:  nil,
			},
			err: errors.New("Schema URL cannot be empty"),
		},
		{
			schemaUrl: "https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism",
			json:      "",
			expect: model.ValidationResult{
				Valid:   false,
				Message: "",
				Errors:  nil,
			},
			err: errors.New("Schema URL cannot be empty"),
		},
	}
	for _, test := range tests {
		validationResult, err := Validate(test.schemaUrl, test.json)
		assert.Equal(t, test.expect, validationResult)
		if test.err != nil {
			assert.NotNil(t, err)
		}
	}
}
