package validator

import (
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"log"
)

func Validate(schemaUrl string, bodyJson string) (model.ValidationResult, error) {
	log.Print("schemaUrl: " + schemaUrl)

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
	var validationErrors []string

	if result.Valid() {
		message = "The document is valid"
	} else {
		message = "The document is not valid"
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.Description())
		}
	}

	vr := model.ValidationResult{
		Valid:   result.Valid(),
		Message: message,
		Errors:  validationErrors,
	}

	return vr, nil
}
