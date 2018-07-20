package utils

import (
	"errors"
	"github.com/codetaming/indy-ingest/_vendor-20180711202013/github.com/stretchr/testify/assert"
	"testing"
)

type testDefinition struct {
	headers map[string]string
	url     string
	err     error
}

func TestReturnUrlForHeaders(t *testing.T) {
	schemaUrl := "https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism"
	tests := []testDefinition{
		{
			headers: map[string]string{
				"link": `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>; rel="describedby"`,
			},
			url: schemaUrl,
			err: nil,
		},
		{
			headers: map[string]string{
				"Link": `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>; rel="describedby"`,
			},
			url: schemaUrl,
			err: nil,
		},
		{
			headers: map[string]string{},
			url:     "",
			err:     errors.New("Link header must be provided"),
		},
		{
			headers: map[string]string{
				"Link": ``,
			},
			url: "",
			err: errors.New("Link header must be provided"),
		},
		{
			headers: map[string]string{
				"Link": `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>`,
			},
			url: schemaUrl,
			err: nil,
		},
	}
	for _, test := range tests {
		url, err := ExtractSchemaUrl(test.headers)
		if nil == test.err {
			assert.Equal(t, test.err, err)
		}
		if nil != test.err {
			assert.Equal(t, test.err, err)
		}
		assert.Equal(t, test.url, url)
	}
}
