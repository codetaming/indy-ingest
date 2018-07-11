package validator

import (
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/xeipuuv/gojsonschema"
)

func Validate(schemaUrl string, bodyJson string) model.ValidationResult {

	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(bodyJson)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		panic(err.Error())
	}

	var message string
	var errors []string

	if result.Valid() {
		message = "The document is valid"
	} else {
		message = "The document is not valid"
		for _, desc := range result.Errors() {
			errors = append(errors, desc.Description())
		}
	}

	vr := model.ValidationResult{
		Valid:   result.Valid(),
		Message: message,
		Errors:  errors,
	}

	return vr
}
