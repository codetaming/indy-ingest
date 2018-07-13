package validator

import (
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

func Validate(schemaUrl string, bodyJson string) (model.ValidationResult, error) {
	log.Info("schemaUrl: " + schemaUrl)

	if schemaUrl == "" {
		return model.ValidationResult{}, errors.New("Schema URL cannot be empty")
	}

	if bodyJson == "" {
		return model.ValidationResult{}, errors.New("JSON cannot be empty")
	}
	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(bodyJson)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		return model.ValidationResult{}, err
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

	return vr, nil
}
